package samples

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

	engine.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index", "Hello Gin")
	})
}
