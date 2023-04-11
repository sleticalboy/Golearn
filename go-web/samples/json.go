package samples

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Json() {
	engine := gin.Default()
	fmt.Printf("json run... %v\n", engine.BasePath())

	engine.GET("/api/v1/hello", func(context *gin.Context) {
		context.JSON(http.StatusOK, map[string]string{
			"hello":    "Gin",
			"language": "Go",
			"desc":     "Go web 编程",
		})
	})

	err := engine.Run("127.0.0.1:8099")
	if err != nil {
		fmt.Println(err)
	}
}
