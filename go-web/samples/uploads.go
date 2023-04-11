package samples

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"os"
)

var pwd, _ = os.Getwd()

type UploadError struct {
	Status int
	Msg    string
	Data   string
}

func (e UploadError) Error() string {
	return fmt.Sprintf(`status: %d, msg: %s, data: %s`, e.Status, e.Msg, e.Data)
}

func saveFile(context *gin.Context, name string, file *multipart.FileHeader) (string, *UploadError) {
	dest := fmt.Sprintf("%s/../build/%s", pwd, file.Filename)
	fmt.Printf("src: %s, dest: %s\n", file.Filename, dest)
	err := context.SaveUploadedFile(file, dest)
	if err != nil {
		return "", &UploadError{
			Status: http.StatusInternalServerError,
			Msg:    "Save file error",
			Data:   "",
		}
	}
	return dest, nil
}

func saveFiles(context *gin.Context, k string, files []*multipart.FileHeader) ([]string, *UploadError) {
	dests := make([]string, len(files))
	for i, file := range files {
		if dest, err := saveFile(context, k, file); err == nil {
			dests[i] = dest
		} else {
			return nil, err
		}
	}
	return dests, nil
}

func response(context *gin.Context, dests []string, err *UploadError) {
	if err != nil {
		context.JSON(err.Status, err.Error())
	} else {
		data := map[string]string{}
		for i, dest := range dests {
			data[fmt.Sprintf("file-%d", i)] = dest
		}
		context.JSON(http.StatusOK, map[string]any{
			"status": http.StatusOK,
			"msg":    "ok",
			"data":   data,
		})
	}
}

func Uploads(engine *gin.Engine) {
	fmt.Printf("upload run...%s\n", engine.BasePath())

	engine.POST("/api/v1/upload", func(context *gin.Context) {
		if form, _ := context.MultipartForm(); form != nil && len(form.File) != 0 {
			fmt.Printf("upload files: %v, value: %v\n", form.File, form.Value)
			// 多文件上传
			// curl -X POST http://127.0.0.1:8099/api/v1/upload -H "Content-Type: multipart/form-data" -F "file=@/path/to/your/file"
			files := form.File["upload[]"]
			if files != nil && len(files) > 0 {
				dests, err := saveFiles(context, "", files)
				response(context, dests, err)
				return
			}
		}
		// 单文件上传
		// curl -X POST http://127.0.0.1:8099/api/v1/upload -H "Content-Type: multipart/form-data" -F "upload[]=@/path/to/your/file" -F "upload[]=@/path/to/your/file"
		if file, err := context.FormFile("file"); err == nil {
			fmt.Printf("upload file: %s\n", file.Filename)
			dest, err := saveFile(context, file.Filename, file)
			response(context, []string{0: dest}, err)
		}
	})
}
