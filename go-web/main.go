package main

import (
	"com.binlee/goweb/samples"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func main() {
	engine := gin.Default()

	// 将日志文件写入文件
	// gin.DisableConsoleColor()
	// file, _ := os.OpenFile("../build/gin.log", os.O_CREATE|os.O_APPEND, 0644)
	// gin.DefaultWriter = io.MultiWriter(file)

	samples.Html(engine)
	samples.Json(engine)
	samples.FormLogin(engine)
	samples.Others(engine)
	samples.Uploads(engine)
	samples.Downloads(engine)
	err := engine.Run("127.0.0.1:8099")
	if err != nil {
		fmt.Println(err)
	}

	// 可以运行多个服务
	g := errgroup.Group{}
	g.Go(func() error {
		return nil
	})
}
