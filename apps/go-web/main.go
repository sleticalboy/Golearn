package main

import (
	"fmt"
	"github.com/donovanhide/eventsource"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/testdata/protoexample"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func Html(ctx *gin.Context) {
	// curl http://127.0.0.1:8099/
	ctx.HTML(http.StatusOK, "index.html", "Hello Gin")
}

func Json(ctx *gin.Context) {
	// curl http://127.0.0.1:8099/api/v1/hello
	ctx.JSON(http.StatusOK, gin.H{
		"hello":    "Gin",
		"language": "Go",
		"desc":     "Go web 编程",
	})
	// ctx.JSONP(http.StatusOK, nil)
}

type LoginForm struct {
	User string `form:"user" binding:"required"`
	Pwd  string `form:"pwd" binding:"required"`
}

func FormLogin(ctx *gin.Context) {
	// form 绑定结构体
	// curl -X POST "http://localhost:8099/api/v1/login" --form user=user --form pwd=root
	var form = LoginForm{}
	if ctx.ShouldBind(&form) == nil {
		if form.User == "binli" && form.Pwd == "root" {
			ctx.JSON(http.StatusOK, map[string]string{
				"status": strconv.Itoa(http.StatusOK),
				"msg":    "OK",
				"data":   "Login success",
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, map[string]string{
				"status": strconv.Itoa(http.StatusUnauthorized),
				"msg":    "Unauthorized",
				"data":   "",
			})
		}
	}
}

func FormLogin2(ctx *gin.Context) {
	// 主动解析
	// curl -X POST "http://localhost:8099/api/v1/login2" --form user=user --form pwd=root_tt
	user, hit := ctx.GetPostForm("user")
	fmt.Printf("for login user: %s, hit: %t\n", user, hit)
	pwd, hit := ctx.GetPostForm("pwd")
	fmt.Printf("for login pwd: %s, hit: %t\n", pwd, hit)
	ctx.JSON(http.StatusOK, map[string]string{
		"status": strconv.Itoa(http.StatusOK),
		"msg":    "echo server",
		"data":   fmt.Sprintf(`{"user":"%s", "pwd":"%s"}`, user, pwd),
	})
}

func FormLogin3(ctx *gin.Context) {
	// curl -X POST "http://localhost:8099/api/v1/login3?id=123&page=5" -H "Content-Type: application/x-www-form-urlencoded" -d "user=binli&pwd=root"
	id, hit := ctx.GetQuery("id")
	fmt.Printf("for login id: %s, hit: %t\n", id, hit)
	page, hit := ctx.GetQuery("page")
	fmt.Printf("for login page: %s, hit: %t\n", page, hit)
	user, hit := ctx.GetPostForm("user")
	fmt.Printf("for login user: %s, hit: %t\n", user, hit)
	pwd, hit := ctx.GetPostForm("pwd")
	fmt.Printf("for login pwd: %s, hit: %t\n", pwd, hit)
	ctx.JSON(http.StatusOK, map[string]string{
		"status": strconv.Itoa(http.StatusOK),
		"msg":    "echo server",
		"data":   fmt.Sprintf(`{"user":"%s", "pwd":"%s", "id":"%s", "page":"%s"}`, user, pwd, id, page),
	})
}

type Xml struct {
	Ns    string
	Key   string
	Value string
}

func xml(ctx *gin.Context) {
	// curl -X GET "http://localhost:8099/api/v1/xml"
	// 这里的 obj 需要实现 xml.Marshaler 接口
	ctx.XML(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    "xml echor",
		"data": Xml{
			Ns:    "xxxx",
			Key:   "user",
			Value: "binli",
		},
	})
}

func handleJson(ctx *gin.Context) {
	// curl -X GET "http://localhost:8099/api/v1/json"
	// 使用结构体作为 response
	msg := struct {
		Name   string `json:"user"`
		Age    int
		Gender int
	}{
		Name:   "anonymous struct",
		Age:    23,
		Gender: 1,
	}
	ctx.JSON(http.StatusOK, msg)
}

func yaml(ctx *gin.Context) {
	// curl -X GET "http://localhost:8099/api/v1/yaml"
	ctx.YAML(http.StatusOK, map[string]any{
		"status": http.StatusOK,
		"msg":    "yaml echor",
		"data":   "yaml data",
	})
}

func protobuf(ctx *gin.Context) {
	// curl -X GET "http://localhost:8099/api/v1/protobuf"
	label := "test\ntttt\n"
	reps := []int64{1, 2}
	data := &protoexample.Test{
		Label: &label,
		Reps:  reps,
	}
	ctx.ProtoBuf(http.StatusOK, data)
}

func download(ctx *gin.Context) {
	// curl "http://127.0.0.1:8099/api/v1/download?fileId={0/1}"
	fid, hit := ctx.GetQuery("fileId")
	fmt.Printf("download file id: %s, hit: %t\n", fid, hit)
	if hit {
		name := ""
		if fid == "0" {
			name = "2023-03-20_11-47-36.aac"
		} else if fid == "1" {
			name = "PCMalaw.wav"
		}
		var pwd, _ = os.Getwd()
		if name != "" {
			dest := fmt.Sprintf("%s/../build/%s", pwd, name)
			fmt.Printf("download with: %s\n", dest)
			ctx.FileAttachment(dest, name)
		} else {
			dest := fmt.Sprintf("%s/../build/test.txt", pwd)
			fmt.Printf("download with: %s\n", dest)
			info, err := os.Stat(dest)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			file, err := os.Open(dest)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer func() { _ = file.Close() }()
			// 从 reader 读取数据，也可以从另一个请求读取数据写入到这个请求里面
			ctx.DataFromReader(http.StatusOK, info.Size(), "application/octet-stream", file,
				map[string]string{
					"Content-Disposition": `attachment; finename="test.txt"`,
				},
			)
		}
	}
	// 在协程中使用 ctx 时只能使用其副本
	go func(cc *gin.Context) {
		fmt.Printf("request url: %s\n", cc.Request.URL)
	}(ctx.Copy())
}

type UploadError struct {
	Status int
	Msg    string
	Data   string
}

func (e UploadError) Error() string {
	return fmt.Sprintf(`status: %d, msg: %s, data: %s`, e.Status, e.Msg, e.Data)
}

func upload(ctx *gin.Context) {
	failed := func(err *UploadError) { ctx.JSON(err.Status, err.Error()) }

	succeeded := func(dests []string) {
		data := map[string]string{}
		for i, dest := range dests {
			data[fmt.Sprintf("file-%d", i)] = dest
		}
		ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "ok", "data": data})
	}

	saveFiles := func(ctx *gin.Context, files []*multipart.FileHeader) {
		if files == nil || len(files) == 0 {
			failed(&UploadError{Status: http.StatusBadRequest, Msg: "no files"})
			return
		}
		var pwd, _ = os.Getwd()
		dstFiles := make([]string, len(files))
		for i, file := range files {
			dest := fmt.Sprintf("%s/../build/%s", pwd, file.Filename)
			fmt.Printf("src: %s, dest: %s\n", file.Filename, dest)
			if err := ctx.SaveUploadedFile(file, dest); err != nil {
				failed(&UploadError{Status: http.StatusInternalServerError, Msg: "Save file error"})
				break
			}
			dstFiles[i] = dest
		}
		succeeded(dstFiles)
	}

	if form, _ := ctx.MultipartForm(); form != nil && len(form.File) != 0 {
		fmt.Printf("upload files: %v, value: %v\n", form.File, form.Value)
		// 多文件上传
		// curl -X POST http://127.0.0.1:8099/api/v1/upload -H "Content-Type: multipart/form-data" -F "file=@/path/to/your/file"
		saveFiles(ctx, form.File["upload[]"])
	}
	// 单文件上传
	// curl -X POST http://127.0.0.1:8099/api/v1/upload -H "Content-Type: multipart/form-data" -F "upload[]=@/path/to/your/file" -F "upload[]=@/path/to/your/file"
	if file, err := ctx.FormFile("file"); err == nil {
		fmt.Printf("upload file: %s\n", file.Filename)
		saveFiles(ctx, []*multipart.FileHeader{file})
	}
}

func login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0, "msg": "OK",
	})
}

func favicon(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "favico.ico", nil)
}

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
		time.Sleep(5 * time.Second)
		stream.Publish([]string{"events"}, &mysseEvent{})
		log.Println("Broadcasted time update")
	}
}

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
	go func() {
		for {
			var msg wsmessage
			if err := conn.ReadJSON(&msg); err != nil {
				log.Println("wsHello() Read error:", err)
				break
			}
			log.Printf("wsHello() Received: %+v", msg)
			// 回显消息
			if err = conn.WriteJSON(msg); err != nil {
				log.Println("wsHello() Write error:", err)
				break
			}
		}
	}()
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

func main() {
	r := gin.Default()

	// 将日志文件写入文件
	// gin.DisableConsoleColor()
	// file, _ := os.OpenFile("../build/gin.log", os.O_CREATE|os.O_APPEND, 0644)
	// gin.DefaultWriter = io.MultiWriter(file)

	// 静态资源目录
	// r.Static("assets", "resources")
	r.StaticFS("assets", gin.Dir("resources", false))
	// r.StaticFile("/favicon.ico", "resources/img/")
	r.GET("/favicon.ico", favicon)
	// 模板路径
	r.LoadHTMLGlob("templates/**")
	r.GET("/", Html)
	r.GET("/api/v1/hello", Json)

	r.POST("/api/v1/login", FormLogin)
	r.POST("/api/v1/login2", FormLogin2)
	r.POST("/api/v1/login3", FormLogin3)
	r.GET("/api/v1/xml", xml)
	r.GET("/api/v1/json", handleJson)
	r.GET("/api/v1/yaml", yaml)
	r.GET("/api/v1/protobuf", protobuf)
	r.POST("/api/v1/upload", upload)

	r.GET("/api/v1/download", download)

	// server side event
	go sseBroadcast()
	r.GET("/api/v1/sse", sse)
	r.GET("/api/v1/ssestream", ssestream)
	r.POST("/api/v1/ssesend", ssesend)

	// websocket
	go wsBroadcast()
	// // 静态文件（前端页面）
	// r.StaticFile("/", "./index.html")
	r.GET("/ws", ws)
	r.POST("/api/v1/broadcast", broadcastWs)

	// 可以运行多个服务
	g := errgroup.Group{}
	g.Go(func() error {
		fmt.Println("what ...")
		// r.RunTLS(":8443", "", "")
		return r.Run(":8099")
	})
	err := g.Wait()
	fmt.Printf("wait running %v ...\n", err)
}
