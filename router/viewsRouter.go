package router

import (
	"context"

	"github.com/AYGA2K/photo-editor-web-app/webapp/middleware"
	"github.com/AYGA2K/photo-editor-web-app/webapp/views"
	"github.com/kataras/iris/v12"
)

func Views(app *iris.Application) {
	app.Get("/", middleware.Authenticated, func(ctx iris.Context) {
		index := views.Index()
		index.Render(context.Background(), ctx.ResponseWriter())
	})
	app.Get("/signup", func(ctx iris.Context) {
		signup := views.Signup()
		signup.Render(context.Background(), ctx.ResponseWriter())
	})
	app.Get("/login", func(ctx iris.Context) {
		login := views.Login()
		login.Render(context.Background(), ctx.ResponseWriter())
	})
	app.Get("/home", middleware.Authenticated, func(ctx iris.Context) {
		home := views.Home()
		home.Render(context.Background(), ctx.ResponseWriter())
	})
	app.Get("/profile", middleware.Authenticated, func(ctx iris.Context) {
		profile := views.Profile()
		profile.Render(context.Background(), ctx.ResponseWriter())
	})
}
