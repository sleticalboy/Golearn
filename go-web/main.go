package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/testdata/protoexample"
	"golang.org/x/sync/errgroup"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

func Html(context *gin.Context) {
	// curl http://127.0.0.1:8099/
	context.HTML(http.StatusOK, "index.html", "Hello Gin")
}

func Json(context *gin.Context) {
	// curl http://127.0.0.1:8099/api/v1/hello
	context.JSON(http.StatusOK, gin.H{
		"hello":    "Gin",
		"language": "Go",
		"desc":     "Go web 编程",
	})
	// context.JSONP(http.StatusOK, nil)
}

type LoginForm struct {
	User string `form:"user" binding:"required"`
	Pwd  string `form:"pwd" binding:"required"`
}

func FormLogin(context *gin.Context) {
	// form 绑定结构体
	// curl -X POST "http://localhost:8099/api/v1/login" --form user=user --form pwd=root
	var form = LoginForm{}
	if context.ShouldBind(&form) == nil {
		if form.User == "binli" && form.Pwd == "root" {
			context.JSON(http.StatusOK, map[string]string{
				"status": strconv.Itoa(http.StatusOK),
				"msg":    "OK",
				"data":   "Login success",
			})
		} else {
			context.JSON(http.StatusUnauthorized, map[string]string{
				"status": strconv.Itoa(http.StatusUnauthorized),
				"msg":    "Unauthorized",
				"data":   "",
			})
		}
	}
}

func FormLogin2(context *gin.Context) {
	// 主动解析
	// curl -X POST "http://localhost:8099/api/v1/login2" --form user=user --form pwd=root_tt
	user, hit := context.GetPostForm("user")
	fmt.Printf("for login user: %s, hit: %t\n", user, hit)
	pwd, hit := context.GetPostForm("pwd")
	fmt.Printf("for login pwd: %s, hit: %t\n", pwd, hit)
	context.JSON(http.StatusOK, map[string]string{
		"status": strconv.Itoa(http.StatusOK),
		"msg":    "echo server",
		"data":   fmt.Sprintf(`{"user":"%s", "pwd":"%s"}`, user, pwd),
	})
}

func FormLogin3(context *gin.Context) {
	// curl -X POST "http://localhost:8099/api/v1/login3?id=123&page=5" -H "Content-Type: application/x-www-form-urlencoded" -d "user=binli&pwd=root"
	id, hit := context.GetQuery("id")
	fmt.Printf("for login id: %s, hit: %t\n", id, hit)
	page, hit := context.GetQuery("page")
	fmt.Printf("for login page: %s, hit: %t\n", page, hit)
	user, hit := context.GetPostForm("user")
	fmt.Printf("for login user: %s, hit: %t\n", user, hit)
	pwd, hit := context.GetPostForm("pwd")
	fmt.Printf("for login pwd: %s, hit: %t\n", pwd, hit)
	context.JSON(http.StatusOK, map[string]string{
		"status": strconv.Itoa(http.StatusOK),
		"msg":    "echo server",
		"data":   fmt.Sprintf(`{"user":"%s", "pwd":"%s", "id":"%s", "page":"%s"}`, user, pwd, id, page),
	})
}

type Xml struct {
	Ns    string
	Key   string
	Value string
}

func xml(context *gin.Context) {
	// curl -X GET "http://localhost:8099/api/v1/xml"
	// 这里的 obj 需要实现 xml.Marshaler 接口
	context.XML(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    "xml echor",
		"data": Xml{
			Ns:    "xxxx",
			Key:   "user",
			Value: "binli",
		},
	})
}

func json(context *gin.Context) {
	// curl -X GET "http://localhost:8099/api/v1/json"
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
}

func yaml(context *gin.Context) {
	// curl -X GET "http://localhost:8099/api/v1/yaml"
	context.YAML(http.StatusOK, map[string]any{
		"status": http.StatusOK,
		"msg":    "yaml echor",
		"data":   "yaml data",
	})
}

func protobuf(context *gin.Context) {
	// curl -X GET "http://localhost:8099/api/v1/protobuf"
	label := "test\ntttt\n"
	reps := []int64{1, 2}
	data := &protoexample.Test{
		Label: &label,
		Reps:  reps,
	}
	context.ProtoBuf(http.StatusOK, data)
}

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
	// 在协程中使用 context 时只能使用其副本
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

func upload(context *gin.Context) {
	failed := func(err *UploadError) { context.JSON(err.Status, err.Error()) }

	succeeded := func(dests []string) {
		data := map[string]string{}
		for i, dest := range dests {
			data[fmt.Sprintf("file-%d", i)] = dest
		}
		context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "ok", "data": data})
	}

	saveFiles := func(context *gin.Context, files []*multipart.FileHeader) {
		if files == nil || len(files) == 0 {
			failed(&UploadError{Status: http.StatusBadRequest, Msg: "no files"})
			return
		}
		var pwd, _ = os.Getwd()
		dstFiles := make([]string, len(files))
		for i, file := range files {
			dest := fmt.Sprintf("%s/../build/%s", pwd, file.Filename)
			fmt.Printf("src: %s, dest: %s\n", file.Filename, dest)
			if err := context.SaveUploadedFile(file, dest); err != nil {
				failed(&UploadError{Status: http.StatusInternalServerError, Msg: "Save file error"})
				break
			}
			dstFiles[i] = dest
		}
		succeeded(dstFiles)
	}

	if form, _ := context.MultipartForm(); form != nil && len(form.File) != 0 {
		fmt.Printf("upload files: %v, value: %v\n", form.File, form.Value)
		// 多文件上传
		// curl -X POST http://127.0.0.1:8099/api/v1/upload -H "Content-Type: multipart/form-data" -F "file=@/path/to/your/file"
		saveFiles(context, form.File["upload[]"])
	}
	// 单文件上传
	// curl -X POST http://127.0.0.1:8099/api/v1/upload -H "Content-Type: multipart/form-data" -F "upload[]=@/path/to/your/file" -F "upload[]=@/path/to/your/file"
	if file, err := context.FormFile("file"); err == nil {
		fmt.Printf("upload file: %s\n", file.Filename)
		saveFiles(context, []*multipart.FileHeader{file})
	}
}

func login(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"code": 0, "msg": "OK",
	})
}

func favicon(context *gin.Context) {
	context.HTML(http.StatusOK, "favico.ico", nil)
}

func main() {
	server := gin.Default()

	// 将日志文件写入文件
	// gin.DisableConsoleColor()
	// file, _ := os.OpenFile("../build/gin.log", os.O_CREATE|os.O_APPEND, 0644)
	// gin.DefaultWriter = io.MultiWriter(file)

	// 静态资源目录
	// server.Static("assets", "resources")
	server.StaticFS("assets", gin.Dir("resources", false))
	// server.StaticFile("/favicon.ico", "resources/img/")
	server.GET("/favicon.ico", favicon)
	// 模板路径
	server.LoadHTMLGlob("resources/templates/**")
	server.GET("/", Html)
	server.GET("/api/v1/hello", Json)

	server.POST("/api/v1/login", FormLogin)
	server.POST("/api/v1/login2", FormLogin2)
	server.POST("/api/v1/login3", FormLogin3)
	server.GET("/api/v1/xml", xml)
	server.GET("/api/v1/json", json)
	server.GET("/api/v1/yaml", yaml)
	server.GET("/api/v1/protobuf", protobuf)
	server.POST("/api/v1/upload", upload)

	server.GET("/api/v1/download", download)

	// 可以运行多个服务
	g := errgroup.Group{}
	g.Go(func() error {
		fmt.Println("what ...")
		// server.RunTLS(":8443", "", "")
		return server.Run(":8099")
	})
	err := g.Wait()
	fmt.Printf("wait running %v ...\n", err)
}
