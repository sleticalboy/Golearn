package core

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var liveRoomId string
var LIVE_GIFT_LIST = map[string]bool{
	"月下瀑布": true,
	"小甜甜":  true,
}
var Gval = &gvalue{
	GiftList:   make([]string, 0),
	GiftIDList: make([]string, 0),
	RankUser:   make([]string, 0),
}

//	LiveInfo {
//	 "room_url": "https://live.douyin.com/878478511352",
//	 "room_title": "",
//	 "room_id": "7519127586287504191",
//	 "wss_url": "wss://webcast5-ws-web-lf.douyin.com/webcast/im/push/v2/?room_id=7519127586287504191&compress=gzip&version_code=180800&webcast_sdk_version=1.0.14-beta.0&live_id=1&did_rule=3&user_unique_id=7692890738690032830&identity=audience&signature=642vnD3eBZYzEGUq&aid=6383&device_platform=web&browser_language=zh-CN&browser_platform=Win32&browser_name=Mozilla&browser_version=5.0+%28Windows+NT+10.0%3B+Win64%3B+x64%29+AppleWebKit%2F537.36+%28KHTML%2C+like+Gecko%29+Chrome%2F126.0.0.0+Safari%2F537.36+Edg%2F126.0.0.0",
//	 "ttwid": "1%7C9SEGPfK9oK2Ku60vf6jyt7h6JWbBu4N_-kwQdU-SPd8%7C1697721607%7Cc406088cffa073546db29932058720720521571b92ba67ba902a70e5aaffd5d6",
//	 "m3u8_url": "https://pull-hls-l11.douyincdn.com/thirdgame/stream-7519127303609862966_or4.m3u8?arch_hrchy=w1&exp_hrchy=w1&expire=1751292702&major_anchor_level=common&sign=10bffb21dfa081b1dc3e7e30eda3d7b0&t_id=037-202506232211422FD9CA8B3E0AE7271611-4BkLXt",
//	 "flv_url": "https://pull-flv-l11.douyincdn.com/thirdgame/stream-7519127303609862966.flv?arch_hrchy=w1\\u0026exp_hrchy=w1\\u0026expire=1751292702\\u0026major_anchor_level=common\\u0026only_audio=1\\u0026sign=be7289c18f329dbe569788cb301aa85d\\u0026t_id=037-202506232211422FD9CA8B3E0AE7271611-4BkLXt\\u0026unique_id=stream-7519127303609862966_844_flv",
//	 "user_info": {
//	   "id_str": "77891600554",
//	   "sec_uid": "MS4wLjABAAAAfVEt12vrB1PEkuEpDms2zoER8ZLXbfSCDdshc0Wq_c4",
//	   "nickname": "小鹏爻开味",
//	   "avatar_thumb": {
//	     "url_list": [
//	       "https://p3.douyinpic.com/aweme/100x100/aweme-avatar/tos-cn-i-0813c001_o46IECynABAAZ3AHATi4AezACAGgAC2QILfvCi.jpeg?from=3067671334",
//	       "https://p26.douyinpic.com/aweme/100x100/aweme-avatar/tos-cn-i-0813c001_o46IECynABAAZ3AHATi4AezACAGgAC2QILfvCi.jpeg?from=3067671334",
//	       "https://p11.douyinpic.com/aweme/100x100/aweme-avatar/tos-cn-i-0813c001_o46IECynABAAZ3AHATi4AezACAGgAC2QILfvCi.jpeg?from=3067671334"
//	     ]
//	   },
//	   "follow_info": {
//	     "follow_status": 0,
//	     "follow_status_str": "0"
//	   },
//	   "subscribe": {
//	     "is_member": false,
//	     "level": 0,
//	     "identity_type": 0,
//	     "buy_type": 0,
//	     "open": 0
//	   },
//	   "foreign_user": 0,
//	   "open_id_str": ""
//	 },
//	 "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0"
//	}
type LiveInfo struct {
	RoomUrl   string `json:"room_url"`
	RoomTitle string `json:"room_title"`
	RoomId    string `json:"room_id"`
	WssUrl    string `json:"wss_url"`
	Ttwid     string `json:"ttwid"`
	M3U8Url   string `json:"m3u8_url"`
	FlvUrl    string `json:"flv_url"`
	UserInfo  struct {
		IdStr       string `json:"id_str"`
		SecUid      string `json:"sec_uid"`
		Nickname    string `json:"nickname"`
		AvatarThumb struct {
			UrlList []string `json:"url_list"`
		} `json:"avatar_thumb"`
		FollowInfo struct {
			FollowStatus    int    `json:"follow_status"`
			FollowStatusStr string `json:"follow_status_str"`
		} `json:"follow_info"`
		Subscribe struct {
			IsMember     bool `json:"is_member"`
			Level        int  `json:"level"`
			IdentityType int  `json:"identity_type"`
			BuyType      int  `json:"buy_type"`
			Open         int  `json:"open"`
		} `json:"subscribe"`
		ForeignUser int    `json:"foreign_user"`
		OpenIdStr   string `json:"open_id_str"`
	} `json:"user_info"`
	UserAgent string `json:"user_agent"`
}

type gvalue struct {
	// 点赞总数
	LikeNum int

	// 评论总数
	CommitNum int

	// 礼物数量和价值
	GiftNum   int
	GiftValue int

	// 特殊礼物：月下瀑布等
	GiftList []string

	// 礼物 ID 列表：用于去重
	GiftIDList []string

	// 直播间人数
	MemberNum int

	// 在线观众排名
	RankUser []string

	// 解析到的消息
	ParsedMsg proto.Message

	allEnterUsers int // 所有进入直播间的人数
	validUsers    int // 非匿名人数
}

func init() {
	log.SetPrefix("\n")
}

// GetLiveInfo fetches room metadata in pure Go (no main.exe).
// Signature still stubbed until JS engine is wired.
func GetLiveInfo(liveUrl string) *LiveInfo {
	info, err := ParseLiveRoomURL(liveUrl)
	if err != nil {
		log.Fatalf("GetLiveInfo: %v", err)
	}
	return info
}

func WssServerStart(liveInfo *LiveInfo) {
	var headers = map[string][]string{
		"cookie":     {"ttwid=" + liveInfo.Ttwid},
		"user-agent": {liveInfo.UserAgent},
	}
	// 模拟 WebSocket 连接
	conn, cresp, cerr := websocket.DefaultDialer.Dial(liveInfo.WssUrl, headers)
	var cbuf []byte
	if cresp != nil {
		cbuf, _ = io.ReadAll(cresp.Body)
	}
	log.Printf("WebSocket 连接: %v -> %s", cerr, cbuf)
	if cerr != nil {
		os.Exit(1)
	}
	defer func() { _ = conn.Close() }()
	log.Println("WebSocket 已连接")
	liveRoomId = liveInfo.RoomId
	// 每隔 10s 发送一次心跳包
	go startHeartbeat(conn)
	// 取消息
	for {
		msgType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取消息失败:", err)
			break
		}
		onMessage(conn, msgType, message)
	}
}

func startHeartbeat(conn *websocket.Conn) {
	pingMsg := &PushFrame{PayloadType: "hb"}
	data, err := proto.Marshal(pingMsg)
	if err != nil {
		log.Printf("Failed to marshal Ping message: %v", err)
		return
	}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			log.Printf("WebSocket 发送心跳失败: %v", err)
			return
		}
		log.Printf("[sendHeartbeat] [发送 ping 心跳] [房间Id：%s]", liveRoomId)
	}
}

func onMessage(ws *websocket.Conn, msgType int, message []byte) {
	_ = msgType
	// log.Printf("收到消息: %s", message)
	// 解析 PushFrame 消息
	var wssPackage PushFrame
	err := proto.Unmarshal(message, &wssPackage)
	if err != nil {
		log.Fatalf("Failed to parse PushFrame: %v", err)
	}
	logID := wssPackage.LogId
	decompressed, err := gzip.NewReader(bytes.NewReader(wssPackage.Payload))
	if err != nil {
		log.Printf("Failed to decompress payload: %v", err)
	}
	payloadBytes, _ := io.ReadAll(decompressed)

	var payloadPackage Response
	err = proto.Unmarshal(payloadBytes, &payloadPackage)
	if err != nil {
		log.Fatalf("Failed to parse Response: %v", err)
	}

	// 发送 Ack 包
	if payloadPackage.NeedAck {
		sendAck(ws, logID, payloadPackage.InternalExt)
	}

	for _, msg := range payloadPackage.MessagesList {
		var parsedMsg proto.Message
		switch msg.Method {
		// case "WebcastMatchAgainstScoreMessage":
		// 	parsedMsg = unPackMatchAgainstScoreMessage(msg.Payload)
		// case "WebcastLikeMessage":
		// 	parsedMsg = unPackWebcastLikeMessage(msg.Payload)
		case "WebcastMemberMessage":
			parsedMsg = unPackWebcastMemberMessage(msg.Payload)
		// case "WebcastGiftMessage":
		// 	parsedMsg = unPackWebcastGiftMessage(msg.Payload)
		case "WebcastChatMessage":
		// 	parsedMsg = unPackWebcastChatMessage(msg.Payload)
		// case "WebcastSocialMessage":
		// 	parsedMsg = unPackWebcastSocialMessage(msg.Payload)
		// case "WebcastRoomUserSeqMessage":
		// 	parsedMsg = unPackWebcastRoomUserSeqMessage(msg.Payload)
		// case "WebcastUpdateFanTicketMessage":
		// 	parsedMsg = unPackWebcastUpdateFanTicketMessage(msg.Payload)
		// case "WebcastCommonTextMessage":
		// 	parsedMsg = unPackWebcastCommonTextMessage(msg.Payload)
		// case "WebcastProductChangeMessage":
		// 	parsedMsg = WebcastProductChangeMessage(msg.Payload)
		// case "WebcastRoomStatsMessage":
		// 	fmt.Println("WebcastRoomStatsMessage")
		// 	parsedMsg = nil
		default:
			parsedMsg = nil
			// log.Printf("[onMessage] [待解析方法 %s 等待解析] [房间Id：%s]", msg.Method, liveRoomId)
		}
		Gval.ParsedMsg = parsedMsg
	}
}

func sendAck(ws *websocket.Conn, logId uint64, internalExt string) {
	ack := &PushFrame{
		PayloadType: "ack",
		LogId:       logId,
		// 假设 PayloadType 字段也用于承载 internalExt，否则根据实际结构调整字段名
		// PayloadType: internalExt, // 注意：重复赋值可能需修改字段名，如 InternalExt
	}
	ack.PayloadType = internalExt

	data, err := proto.Marshal(ack)
	if err != nil {
		log.Printf("Failed to marshal Ack: %v", err)
		return
	}
	err = ws.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		log.Printf("Failed to send Ack: %v", err)
		return
	}
	// log.Printf("[sendAck] [🌟发送Ack] [房间Id：%s]", liveRoomId)
}

var pjson = protojson.MarshalOptions{
	Multiline: false,
}

func unPackMatchAgainstScoreMessage(data []byte) proto.Message {
	var msg MatchAgainstScoreMessage
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("Failed to parse MatchAgainstScoreMessage: %v", err)
		return nil
	}
	// 打印日志
	log.Printf("[unPackMatchAgainstScoreMessage] [🌟联谊会分数消息] [房间Id：%s] | %+v", liveRoomId, pjson.Format(&msg))

	return &msg
}

func unPackWebcastLikeMessage(data []byte) proto.Message {
	var msg LikeMessage
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("Failed to parse LikeMessage: %v", err)
		return nil
	}
	// 更新全局点赞数
	Gval.LikeNum = int(msg.Total)
	// 打印日志
	log.Printf("[unPackWebcastLikeMessage] [直播间点赞统计：%d] [房间Id：%s] | %+v",
		msg.Total, liveRoomId, pjson.Format(&msg))
	return &msg
}

func unPackWebcastMemberMessage(data []byte) proto.Message {
	var msg MemberMessage
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("Failed to parse MemberMessage: %v", err)
		return nil
	}
	// 更新全局成员数（假设 Gval 有 MemberNum 字段）
	Gval.MemberNum = int(msg.MemberCount)
	// 打印日志
	log.Printf("[unPackWebcastMemberMessage] [直播间成员加入: %d] [房间Id：%s] | %+v",
		msg.MemberCount, liveRoomId, pjson.Format(&msg))

	Gval.allEnterUsers += 1
	if msg.User != nil && (msg.User.SecUid == "" || msg.User.Id == 111111) {
		Gval.validUsers += 1
		log.Printf("[unPackWebcastMemberMessage] anon/enter %d -> %d, %s uid=%v",
			Gval.allEnterUsers, Gval.validUsers, msg.User.NickName, msg.User.Id)
	}

	return &msg
}

func unPackWebcastGiftMessage(data []byte) proto.Message {
	var msg GiftMessage
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("Failed to parse GiftMessage: %v", err)
		return nil
	}
	// 打印日志
	log.Printf("[unPackWebcastGiftMessage] [直播间礼物消息] [房间Id：%s] | %+v", liveRoomId, pjson.Format(&msg))

	// 提取礼物信息
	giftName := ""
	if msg.Gift != nil {
		giftName = msg.Gift.Name
	}
	nickName := ""
	if msg.User != nil {
		nickName = msg.User.NickName
	}
	traceID := msg.TraceId
	// 判断是否为特殊礼物
	if _, exists := LIVE_GIFT_LIST[giftName]; exists {
		if !contains(Gval.GiftIDList, traceID) {
			log.Printf("抓到特殊礼物了: %s，用户名：%s", giftName, nickName)
			Gval.GiftList = append(Gval.GiftList, nickName+"-"+giftName)
			Gval.GiftIDList = append(Gval.GiftIDList, traceID)
		}
	}
	// 累加礼物数量和价值
	totalCount := int(msg.TotalCount)
	diamondCount := 0
	if msg.Gift != nil {
		diamondCount = int(msg.Gift.DiamondCount)
	}
	Gval.GiftNum += totalCount
	Gval.GiftValue += totalCount * diamondCount
	// 模拟发送礼物信息到自己的 WebSocket 服务（需实现 wsSender）
	// wsSender(fmt.Sprintf("收到礼物: %s，礼物数量:%d，礼物价值: %d", giftName, Gval.GiftNum, Gval.GiftValue))
	return &msg
}

// 辅助函数：判断字符串是否存在于切片中
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func unPackWebcastChatMessage(data []byte) proto.Message {
	var msg ChatMessage
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("Failed to parse ChatMessage: %v", err)
		return nil
	}
	// 更新全局弹幕计数（假设 Gval 有 CommitNum 字段）
	Gval.CommitNum++
	// 打印日志
	log.Printf("[unPackWebcastChatMessage] [直播间弹幕消息：%d] [房间Id：%s] | %+v",
		Gval.CommitNum, liveRoomId, pjson.Format(&msg))
	return &msg
}

func unPackWebcastSocialMessage(data []byte) proto.Message {
	var msg SocialMessage
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("Failed to parse SocialMessage: %v", err)
		return nil
	}
	// 打印日志
	log.Printf("[unPackWebcastSocialMessage] [直播间关注消息] [房间Id：%s] | %+v", liveRoomId, pjson.Format(&msg))
	return &msg
}

func unPackWebcastRoomUserSeqMessage(data []byte) proto.Message {
	var msg RoomUserSeqMessage
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("Failed to parse RoomUserSeqMessage: %v", err)
		return nil
	}
	// 打印日志
	log.Printf("[unPackWebcastRoomUserSeqMessage] [直播间用户序列消息] [房间Id：%s] | %+v", liveRoomId, pjson.Format(&msg))
	return &msg
}

func unPackWebcastUpdateFanTicketMessage(data []byte) proto.Message {
	var msg UpdateFanTicketMessage
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("Failed to parse UpdateFanTicketMessage: %v", err)
		return nil
	}
	// 打印日志
	log.Printf("[unPackWebcastUpdateFanTicketMessage] [直播间粉丝票更新消息] [房间Id：%s] | %+v", liveRoomId, pjson.Format(&msg))
	return &msg
}

func WebcastProductChangeMessage(data []byte) proto.Message {
	var msg ProductChangeMessage
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("Failed to parse ProductChangeMessage: %v", err)
		return nil
	}
	// 打印日志
	log.Printf("[WebcastProductChangeMessage] [直播间商品变更消息] [房间Id：%s] | %+v", liveRoomId, pjson.Format(&msg))
	return &msg
}

func unPackWebcastCommonTextMessage(data []byte) proto.Message {
	var msg CommonTextMessage
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("Failed to parse CommonTextMessage: %v", err)
		return nil
	}
	// 打印日志
	log.Printf("[unPackWebcastCommonTextMessage] [直播间通用文本消息] [房间Id：%s] | %+v", liveRoomId, pjson.Format(&msg))
	return &msg
}
