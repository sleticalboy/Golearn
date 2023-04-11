package main

import (
	"com.binlee/goweb/samples"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
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
}
