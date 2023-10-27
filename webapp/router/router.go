package router

import (
	"github.com/AYGA2K/photo-editor-web-app/webapp/controllers"
	"github.com/AYGA2K/photo-editor-web-app/webapp/middleware"
	"github.com/kataras/iris/v12"
)

func Routes(app *iris.Application) {
	users(app)
}

func users(app *iris.Application) {
	user := app.Party("/user")
	{
		user.Post("/", controllers.CreateUser)
		user.Delete("/", middleware.Authenticated, controllers.DeleteUser)
		user.Post("/profile", middleware.Authenticated, controllers.GetUserInfo)
		user.Post("/login", controllers.Login)
		user.Delete("/logout", controllers.Logout)
		user.Put("/", middleware.Authenticated, controllers.ChangePassword)

	}
}
