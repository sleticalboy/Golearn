package samples

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func Downloads(engine *gin.Engine) {
	fmt.Printf("download run...%s\n", engine.BasePath())

	// curl "http://127.0.0.1:8099/api/v1/download?fileId={0/1}"
	engine.GET("/api/v1/download", func(context *gin.Context) {
		fid, hit := context.GetQuery("fileId")
		fmt.Printf("download file id: %s, hit: %t\n", fid, hit)
		if hit {
			var name string
			if fid == "0" {
				name = "2023-03-20_11-47-36.aac"
			} else {
				name = "PCMalaw.wav"
			}
			var pwd, _ = os.Getwd()
			dest := fmt.Sprintf("%s/../build/%s", pwd, name)
			fmt.Printf("download with: %s\n", dest)
			context.FileAttachment(dest, name)
		}
	})
}
