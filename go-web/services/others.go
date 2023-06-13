package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/testdata/protoexample"
	"net/http"
)

type Xml struct {
	Ns    string
	Key   string
	Value string
}

type XmlMap = gin.H

func Others(engine *gin.Engine) {
	fmt.Printf("other....%v\n", engine.BasePath())

	// curl -X GET "http://localhost:8099/api/v1/xml"
	engine.GET("/api/v1/xml", func(context *gin.Context) {
		// 这里的 obj 需要实现 xml.Marshaler 接口
		context.XML(http.StatusOK, XmlMap{
			"status": http.StatusOK,
			"msg":    "xml echor",
			"data": Xml{
				Ns:    "xxxx",
				Key:   "user",
				Value: "binli",
			},
		})
	})
	// curl -X GET "http://localhost:8099/api/v1/json"
	engine.GET("/api/v1/json", func(context *gin.Context) {
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
		context.JSON(http.StatusOK, msg)
	})
	// curl -X GET "http://localhost:8099/api/v1/yaml"
	engine.GET("/api/v1/yaml", func(context *gin.Context) {
		context.YAML(http.StatusOK, map[string]any{
			"status": http.StatusOK,
			"msg":    "yaml echor",
			"data":   "yaml data",
		})
	})
	// curl -X GET "http://localhost:8099/api/v1/protobuf"
	engine.GET("/api/v1/protobuf", func(context *gin.Context) {
		label := "test\ntttt\n"
		reps := []int64{1, 2}
		data := &protoexample.Test{
			Label: &label,
			Reps:  reps,
		}
		context.ProtoBuf(http.StatusOK, data)
	})
}
