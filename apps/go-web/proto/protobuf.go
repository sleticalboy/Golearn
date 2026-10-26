package proto

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/testdata/protoexample"
	"net/http"
)

func protobufHandler(ctx *gin.Context) {
	// curl -X GET "http://localhost:8099/api/v1/protobuf"
	label := "test\ntttt\n"
	reps := []int64{1, 2}
	data := &protoexample.Test{
		Label: &label,
		Reps:  reps,
	}
	ctx.ProtoBuf(http.StatusOK, data)
}

func Init(r *gin.Engine) {
	r.GET("/api/v1/protobuf", protobufHandler)
}
