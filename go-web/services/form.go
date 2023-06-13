package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type LoginForm struct {
	User string `form:"user" binding:"required"`
	Pwd  string `form:"pwd" binding:"required"`
}

func FormLogin(engine *gin.Engine) {
	fmt.Printf("form login....%v\n", engine.BasePath())
	// form 绑定结构体
	// curl -X POST "http://localhost:8099/api/v1/login" --form user=user --form pwd=root
	engine.POST("/api/v1/login", func(context *gin.Context) {
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
	})
	// 主动解析
	// curl -X POST "http://localhost:8099/api/v1/login2" --form user=user --form pwd=root_tt
	engine.POST("/api/v1/login2", func(context *gin.Context) {
		user, hit := context.GetPostForm("user")
		fmt.Printf("for login user: %s, hit: %t\n", user, hit)
		pwd, hit := context.GetPostForm("pwd")
		fmt.Printf("for login pwd: %s, hit: %t\n", pwd, hit)
		context.JSON(http.StatusOK, map[string]string{
			"status": strconv.Itoa(http.StatusOK),
			"msg":    "echo server",
			"data":   fmt.Sprintf(`{"user":"%s", "pwd":"%s"}`, user, pwd),
		})
	})
	// curl -X POST "http://localhost:8099/api/v1/login3?id=123&page=5" -H "Content-Type: application/x-www-form-urlencoded" -d "user=binli&pwd=root"
	engine.POST("/api/v1/login3", func(context *gin.Context) {
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
	})
}
