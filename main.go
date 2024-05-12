package main

import (
	"log"

	"github.com/AYGA2K/photo-editor-web-app/webapp/database"
	"github.com/AYGA2K/photo-editor-web-app/webapp/router"
	"github.com/kataras/iris/v12"
)

func main() {
	database.ConnectDb()
	app := iris.New()

	app.HandleDir("/dist", "./dist")
	app.HandleDir("/assets", "./assets")
	app.HandleDir("/uploads", "./uploads")
	// app.Use(middleware.Security)
	router.Routes(app)
	router.Views(app)
	err := app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
