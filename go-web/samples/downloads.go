package samples

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func Downloads(engine *gin.Engine) {
	fmt.Printf("download run...%s\n", engine.BasePath())

	// curl "http://127.0.0.1:8099/api/v1/download?fileId={0/1}"
	engine.GET("/api/v1/download", func(context *gin.Context) {
		fid, hit := context.GetQuery("fileId")
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
				context.FileAttachment(dest, name)
			} else {
				dest := fmt.Sprintf("%s/../build/test.txt", pwd)
				fmt.Printf("download with: %s\n", dest)
				info, err := os.Stat(dest)
				if err == nil {
					return
				}
				file, _ := os.Open(dest)
				defer func() {
					_ = file.Close()
				}()
				// 从 reader 读取数据，也可以从另一个请求读取数据写入到这个请求里面
				context.DataFromReader(http.StatusOK, info.Size(), "application/octet-stream", file,
					map[string]string{
						"Content-Disposition": `attachment; finename="test.txt"`,
					},
				)
			}
		}
		// 在协程中使用 context 时只能使用其副本
		go func(cc *gin.Context) {
			fmt.Printf("request url: %s\n", cc.Request.URL)
		}(context.Copy())
	})
}
