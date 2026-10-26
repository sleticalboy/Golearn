package files

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"os"
)

func download(ctx *gin.Context) {
	// curl "http://127.0.0.1:8099/api/v1/download?fileId={0/1}"
	fid, hit := ctx.GetQuery("fileId")
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
			ctx.FileAttachment(dest, name)
		} else {
			dest := fmt.Sprintf("%s/../build/test.txt", pwd)
			fmt.Printf("download with: %s\n", dest)
			info, err := os.Stat(dest)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			file, err := os.Open(dest)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer func() { _ = file.Close() }()
			// 从 reader 读取数据，也可以从另一个请求读取数据写入到这个请求里面
			ctx.DataFromReader(http.StatusOK, info.Size(), "application/octet-stream", file,
				map[string]string{
					"Content-Disposition": `attachment; finename="test.txt"`,
				},
			)
		}
	}
	// 在协程中使用 ctx 时只能使用其副本
	go func(cc *gin.Context) {
		fmt.Printf("request url: %s\n", cc.Request.URL)
	}(ctx.Copy())
}

type UploadError struct {
	Status int
	Msg    string
	Data   string
}

func (e UploadError) Error() string {
	return fmt.Sprintf(`status: %d, msg: %s, data: %s`, e.Status, e.Msg, e.Data)
}

func upload(ctx *gin.Context) {
	failed := func(err *UploadError) { ctx.JSON(err.Status, err.Error()) }

	succeeded := func(dests []string) {
		data := map[string]string{}
		for i, dest := range dests {
			data[fmt.Sprintf("file-%d", i)] = dest
		}
		ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "ok", "data": data})
	}

	saveFiles := func(ctx *gin.Context, files []*multipart.FileHeader) {
		if files == nil || len(files) == 0 {
			failed(&UploadError{Status: http.StatusBadRequest, Msg: "no files"})
			return
		}
		var pwd, _ = os.Getwd()
		dstFiles := make([]string, len(files))
		for i, file := range files {
			dest := fmt.Sprintf("%s/../build/%s", pwd, file.Filename)
			fmt.Printf("src: %s, dest: %s\n", file.Filename, dest)
			if err := ctx.SaveUploadedFile(file, dest); err != nil {
				failed(&UploadError{Status: http.StatusInternalServerError, Msg: "Save file error"})
				break
			}
			dstFiles[i] = dest
		}
		succeeded(dstFiles)
	}

	// 多文件上传
	// curl -X POST http://127.0.0.1:8099/api/v1/upload -H "Content-Type: multipart/form-data" -F "file=@/path/to/your/file"
	if form, _ := ctx.MultipartForm(); form != nil && len(form.File) != 0 {
		fmt.Printf("upload files: %v, value: %v\n", form.File, form.Value)
		saveFiles(ctx, form.File["upload[]"])
		return
	}
	// 单文件上传
	// curl -X POST http://127.0.0.1:8099/api/v1/upload -H "Content-Type: multipart/form-data" -F "upload[]=@/path/to/your/file" -F "upload[]=@/path/to/your/file"
	if file, err := ctx.FormFile("file"); err == nil {
		fmt.Printf("upload file: %s\n", file.Filename)
		saveFiles(ctx, []*multipart.FileHeader{file})
	}
}

func Init(r *gin.Engine) {
	r.POST("/api/v1/upload", upload)
	r.GET("/api/v1/download", download)
}
