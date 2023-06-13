package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Json(engine *gin.Engine) {
	fmt.Printf("json run... %v\n", engine.BasePath())

	// curl http://127.0.0.1:8099/api/v1/hello
	engine.GET("/api/v1/hello", func(context *gin.Context) {
		context.JSON(http.StatusOK, map[string]string{
			"hello":    "Gin",
			"language": "Go",
			"desc":     "Go web 编程",
		})
	})
}
