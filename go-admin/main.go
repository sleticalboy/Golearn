package main

import (
	"log"
	"os"
	"os/signal"

	_ "github.com/GoAdminGroup/go-admin/adapter/gin"               // web framework adapter
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite" // sql driver
	_ "github.com/GoAdminGroup/themes/sword"                       // ui theme

	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/gin-gonic/gin"

	"go-admin/models"
	"go-admin/pages"
	"go-admin/tables"
)

func main() {
	startServer()
}

func startServer() {
	// gin.SetMode(gin.ReleaseMode)
	// gin.DefaultWriter = io.Discard

	server := gin.Default()

	template.AddComp(chartjs.NewChart())

	eng := engine.Default()

	if err := eng.AddConfigFromYAML("./config.yml").
		AddGenerators(tables.Generators).Use(server); err != nil {
		panic(err)
	}

	server.Static("/uploads", "./uploads")

	eng.HTML("get", "/admin", pages.GetDashBoard)
	eng.HTMLFile("get", "/admin/hello", "./html/hello.tmpl", map[string]interface{}{
		"msg": "Hello world",
	})
	models.Init(eng.SqliteConnection())

	_ = server.Run(":8099")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Print("closing database connection")
	eng.SqliteConnection().Close()
}
