package samples

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Xml struct {
	Ns    string
	Key   string
	Value string
}

func Others() {
	engine := gin.Default()
	fmt.Printf("other....%v\n", engine.BasePath())

	// curl -X POST "http://localhost:8099/api/v1/xml" -H "Content-Type: application/xml; charset=utf-8" -d "user=binli&pwd=root"
	engine.GET("/api/v1/xml", func(context *gin.Context) {
		context.XML(http.StatusOK, map[string]any{
			"status": http.StatusOK,
			"msg":    "xml echor",
			"data": &Xml{
				Ns:    "xxxx",
				Key:   "user",
				Value: "binli",
			},
		})
	})

	err := engine.Run("127.0.0.1:8099")
	if err != nil {
		fmt.Println(err)
	}
}
