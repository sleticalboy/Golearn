package sse

import (
	"fmt"
	"github.com/donovanhide/eventsource"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func sse(ctx *gin.Context) {
	// 设置响应头
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	for i := 1; i <= 10; i++ {
		ctx.SSEvent("message", gin.H{
			"id":    i,
			"value": fmt.Sprintf("this is message %d", i),
		})
		// 刷新缓冲区
		flusher.Flush()
		time.Sleep(time.Millisecond * 250)
	}
}

type Client struct {
	ID      string
	Channel chan eventsource.Event
}

var clientId = 0
var stream = eventsource.NewServer()
var clients = make(map[string]Client)

func ssestream(ctx *gin.Context) {
	// return gin.WrapH(stream.Handler("events"))
	// 为每个客户端创建唯一ID
	clientId++
	id := strconv.Itoa(clientId)

	// 创建客户端通道
	ch := make(chan eventsource.Event)
	clients[id] = Client{ID: id, Channel: ch}

	// 设置响应头
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	// 使用eventsource处理连接
	stream.Handler("events")(ctx.Writer, ctx.Request)

	// 客户端断开时清理
	go func() {
		<-ctx.Request.Context().Done()
		delete(clients, id)
		close(ch)
		log.Printf("Client %s disconnected", id)
	}()
}

func ssesend(ctx *gin.Context) {
	message := ctx.PostForm("message")
	if message == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "message is required"})
		return
	}

	stream.Publish([]string{"events"}, &mysseEvent{})
	ctx.JSON(http.StatusOK, gin.H{"status": "message sent"})
}

type mysseEvent struct{}

func (m *mysseEvent) Id() string { return "" }

func (m *mysseEvent) Event() string { return "mysseevent" }

func (m *mysseEvent) Data() string { return "hello" }

var _ eventsource.Event = (*mysseEvent)(nil)

func sseBroadcast() {
	for {
		select {
		case <-time.After(5 * time.Second):
			stream.Publish([]string{"events"}, &mysseEvent{})
		default:
			//
		}
	}
}

func Init(r *gin.Engine) {
	go sseBroadcast()
	r.GET("/api/v1/sse", sse)
	r.GET("/api/v1/ssestream", ssestream)
	r.POST("/api/v1/ssesend", ssesend)
}
