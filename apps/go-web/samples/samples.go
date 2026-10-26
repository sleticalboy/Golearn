package samples

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func html(ctx *gin.Context) {
	// curl http://127.0.0.1:8099/
	ctx.HTML(http.StatusOK, "index.html", "Hello Gin")
}

type LoginForm struct {
	User string `form:"user" binding:"required"`
	Pwd  string `form:"pwd" binding:"required"`
}

func formlogin(ctx *gin.Context) {
	// form 绑定结构体
	// curl -X POST "http://localhost:8099/api/v1/login" --form user=user --form pwd=root
	var form = LoginForm{}
	if ctx.ShouldBind(&form) == nil {
		if form.User == "binli" && form.Pwd == "root" {
			ctx.JSON(http.StatusOK, map[string]string{
				"status": strconv.Itoa(http.StatusOK),
				"msg":    "OK",
				"data":   "Login success",
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, map[string]string{
				"status": strconv.Itoa(http.StatusUnauthorized),
				"msg":    "Unauthorized",
				"data":   "",
			})
		}
	}
}

func formlogin2(ctx *gin.Context) {
	// 主动解析
	// curl -X POST "http://localhost:8099/api/v1/login2" --form user=user --form pwd=root_tt
	user, hit := ctx.GetPostForm("user")
	fmt.Printf("for login user: %s, hit: %t\n", user, hit)
	pwd, hit := ctx.GetPostForm("pwd")
	fmt.Printf("for login pwd: %s, hit: %t\n", pwd, hit)
	ctx.JSON(http.StatusOK, map[string]string{
		"status": strconv.Itoa(http.StatusOK),
		"msg":    "echo server",
		"data":   fmt.Sprintf(`{"user":"%s", "pwd":"%s"}`, user, pwd),
	})
}

func formlogin3(ctx *gin.Context) {
	// curl -X POST "http://localhost:8099/api/v1/login3?id=123&page=5" -H "Content-Type: application/x-www-form-urlencoded" -d "user=binli&pwd=root"
	id, hit := ctx.GetQuery("id")
	fmt.Printf("for login id: %s, hit: %t\n", id, hit)
	page, hit := ctx.GetQuery("page")
	fmt.Printf("for login page: %s, hit: %t\n", page, hit)
	user, hit := ctx.GetPostForm("user")
	fmt.Printf("for login user: %s, hit: %t\n", user, hit)
	pwd, hit := ctx.GetPostForm("pwd")
	fmt.Printf("for login pwd: %s, hit: %t\n", pwd, hit)
	ctx.JSON(http.StatusOK, map[string]string{
		"status": strconv.Itoa(http.StatusOK),
		"msg":    "echo server",
		"data":   fmt.Sprintf(`{"user":"%s", "pwd":"%s", "id":"%s", "page":"%s"}`, user, pwd, id, page),
	})
}

type XmlMsg struct {
	Ns    string
	Key   string
	Value string
}

func xmlHandler(ctx *gin.Context) {
	// curl -X GET "http://localhost:8099/api/v1/xml"
	// 这里的 obj 需要实现 xml.Marshaler 接口
	ctx.XML(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    "xml echor",
		"data": XmlMsg{
			Ns:    "xxxx",
			Key:   "user",
			Value: "binli",
		},
	})
}

func jsonHandler(ctx *gin.Context) {
	// curl http://127.0.0.1:8099/api/v1/hello
	// ctx.JSON(http.StatusOK, gin.H{
	// 	"hello":    "Gin",
	// 	"language": "Go",
	// 	"desc":     "Go web 编程",
	// })
	// curl -X GET "http://localhost:8099/api/v1/json"
	// 使用结构体作为 response
	var msg = struct {
		Name   string `json:"user"`
		Age    int
		Gender int
	}{
		Name:   "anonymous struct",
		Age:    23,
		Gender: 1,
	}
	ctx.JSON(http.StatusOK, msg)
}

func yamlHandler(ctx *gin.Context) {
	// curl -X GET "http://localhost:8099/api/v1/yaml"
	ctx.YAML(http.StatusOK, map[string]any{
		"status": http.StatusOK,
		"msg":    "yaml echor",
		"data":   "yaml data",
	})
}

func Init(r *gin.Engine) {
	r.GET("/", html)
	r.POST("/api/v1/login", formlogin)
	r.POST("/api/v1/login2", formlogin2)
	r.POST("/api/v1/login3", formlogin3)
	r.GET("/api/v1/xml", xmlHandler)
	r.GET("/api/v1/json", jsonHandler)
	r.GET("/api/v1/yaml", yamlHandler)
}
