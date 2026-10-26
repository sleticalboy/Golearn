package ws

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	dypb "github.com/sleticalboy/golearn/apps/go-web/ws/dlive"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// websocket
type wsclient struct {
	conn *websocket.Conn
}

type wsmessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			log.Printf("wsHello() CheckOrigin: %s, uri: %v\n", r.RemoteAddr, r.URL)
			return true // 允许所有跨域请求（生产环境应限制）
		},
	}
	wsClients = sync.Map{}
	broadcast = make(chan []byte)
)

// 处理广播消息
func wsBroadcast() {
	for {
		select {
		case <-time.After(time.Second * 5):
			continue
		case msg := <-broadcast:
			wsClients.Range(func(key, value any) bool {
				client := key.(*wsclient)
				if err := client.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.Println("wsbroadcast() Send error:", err)
					wsClients.Delete(client)
				}
				return true
			})
		}
	}
}

var Methods = []string{
	"WebcastChatMessage", "WebcastRanklistHourEntranceMessage", "WebcastRoomRankMessage",
	"WebcastPreviewCjRpMessage", "WebcastRoomMessage", "WebcastSocialMessage", "WebcastInRoomBannerMessage",
	"WebcastRoomStatsMessage", "WebcastMemberMessage", "LotteryInfoSyncData", "WebcastGameCPAnchorPromoteInfoMessage",
	"WebcastRoomIntroMessage", "WebcastRoomDataSyncMessage", "WebcastGameCPBaseMessage", "WebcastGiftSortMessage",
	"WebcastRoomStreamAdaptationMessage", "WebcastRoomUserSeqMessage",
}

func handleWsClientConn(client *wsclient) {
	defer func() {
		_ = client.conn.Close()
		wsClients.Delete(client)
	}()
	for {
		msgType, buf, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("wsHello() Read error:", err)
			break
		}
		if msgType == websocket.TextMessage {
			var msg wsmessage
			if err := json.Unmarshal(buf, &msg); err != nil {
				log.Println("ws() Read error:", err)
				break
			}
			log.Printf("wsHello() Received: %+v", msg)
			// 回显消息
			if err = client.conn.WriteJSON(msg); err != nil {
				log.Println("ws() Write error:", err)
				break
			}
		} else if msgType == websocket.BinaryMessage {
			// pb message
			var pf dypb.PushFrame
			if err := proto.Unmarshal(buf, &pf); err != nil {
				log.Println("ws() proto.Unmarshal error:", err)
				break
			}
			// gzip 解压 pf.Payload
			gzipr, err := gzip.NewReader(bytes.NewReader(pf.Payload))
			if err != nil {
				log.Println("ws() gzip.NewReader error:", err)
				break
			}
			decompressed, err := io.ReadAll(gzipr)
			if err != nil {
				log.Println("ws() io.ReadAll error:", err)
				_ = gzipr.Close()
				break
			}
			_ = gzipr.Close()
			// pf.Payload = nil
			// log.Printf("ws() Received pf: %+v", &pf)
			pf.Payload = decompressed
			var resp dypb.Response
			if err := proto.Unmarshal(decompressed, &resp); err != nil {
				log.Println("ws() proto.Unmarshal error:", err)
				break
			}
			for _, message := range resp.MessagesList {
				var msgerr error
				var pbmsg proto.Message = nil
				switch message.Method {
				// case "WebcastMatchAgainstScoreMessage":
				// 	var msg dypb.MatchAgainstScoreMessage
				// 	msgerr = proto.Unmarshal(message.Payload, &msg)
				// 	pbmsg = &msg
				// case "WebcastLikeMessage":
				// 	var msg dypb.LikeMessage
				// 	msgerr = proto.Unmarshal(message.Payload, &msg)
				// 	pbmsg = &msg
				case "WebcastMemberMessage":
					var msg dypb.MemberMessage
					msgerr = proto.Unmarshal(message.Payload, &msg)
					pbmsg = &msg
				// case "WebcastGiftMessage":
				// 	var msg dypb.GiftMessage
				// 	msgerr = proto.Unmarshal(message.Payload, &msg)
				// 	pbmsg = &msg
				case "WebcastChatMessage":
					var msg dypb.ChatMessage
					msgerr = proto.Unmarshal(message.Payload, &msg)
					pbmsg = &msg
				// case "WebcastSocialMessage":
				// 	var msg dypb.SocialMessage
				// 	msgerr = proto.Unmarshal(message.Payload, &msg)
				// 	pbmsg = &msg
				// case "WebcastPreviewCjRpMessage":
				// case "WebcastRoomRankMessage":
				// case "WebcastRoomUserSeqMessage":
				// case "WebcastGiftSortMessage":
				// 	fallthrough
				default:
					log.Printf("ws() Received unknown message: %s", message.Method)
				}
				if msgerr == nil {
					log.Printf("ws() Received resp.Method: %v, msg: %v", message.Method, pbmsg)
					if pbmsg != nil {
						if err := client.conn.WriteJSON(pbmsg); err != nil {
							log.Println("ws() Write error:", err)
						}
					}
				} else {
					log.Printf("ws() message proto.Unmarshal err: %v, resp.Method: %s", msgerr, message.Method)
				}
			}
		}
	}
}

func ws(ctx *gin.Context) {
	// 升级 HTTP 连接为 WebSocket
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	client := &wsclient{conn}
	wsClients.Store(client, true)

	// 处理 WebSocket 消息
	go handleWsClientConn(client)
}

func broadcastWs(ctx *gin.Context) {
	data, err := ctx.GetRawData()
	if err != nil {
		log.Println("wsbroadcast() Read error:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(data) > 0 {
		broadcast <- data
		ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "ok"})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "data is empty"})
	}
}

func Init(r *gin.Engine) {
	go wsBroadcast()
	// // 静态文件（前端页面）
	// r.StaticFile("/", "./index.html")
	r.GET("/ws", ws)
	r.POST("/api/v1/broadcast", broadcastWs)
}
