package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

func Html(engine *gin.Engine) {
	fmt.Printf("html run... %v\n", engine.BasePath())

	t, err := template.New("index").Parse(`<html><body><h1>{{.}}</h1></body></html>`)
	engine.SetHTMLTemplate(template.Must(t, err))
	// engine.LoadHTMLFiles()
	// engine.LoadHTMLGlob("")

	// curl http://127.0.0.1:8099/
	engine.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index", "Hello Gin")
	})
}
