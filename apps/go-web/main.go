package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sleticalboy/golearn/apps/go-web/proto"
	"github.com/sleticalboy/golearn/apps/go-web/samples"
	"github.com/sleticalboy/golearn/apps/go-web/sse"
	"github.com/sleticalboy/golearn/apps/go-web/ws"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func favicon(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "favico.ico", nil)
}

func main() {
	r := gin.Default()

	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000"}, // 你的前端地址
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))
	r.Use(cors.Default())

	// 将日志文件写入文件
	// gin.DisableConsoleColor()
	// file, _ := os.OpenFile("../build/gin.log", os.O_CREATE|os.O_APPEND, 0644)
	// gin.DefaultWriter = io.MultiWriter(file)

	// 静态资源目录
	// r.Static("assets", "resources")
	r.StaticFS("static", gin.Dir("assets", false))
	// r.StaticFile("/favicon.ico", "resources/img/")
	r.GET("/favicon.ico", favicon)
	// 模板路径
	{
		dir, _ := os.Getwd()
		matches, err := filepath.Glob("apps/go-web/templates/**")
		log.Printf("wd: %s, matches: %v, err: %v", dir, matches, err)
	}
	r.LoadHTMLGlob("apps/go-web/templates/**")
	// samples
	samples.Init(r)
	// protobuf handler
	proto.Init(r)
	// server side event
	sse.Init(r)
	// websocket
	ws.Init(r)

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
