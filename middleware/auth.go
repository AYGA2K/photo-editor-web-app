package middleware

import (
	"github.com/AYGA2K/photo-editor-web-app/webapp/controllers"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

func Authenticated(ctx iris.Context) {
	sessionid := ctx.GetCookie("sessionid")
	if sessionid == "" {
		ctx.Redirect("/login")
		return
	}
	sessionUuid, err := uuid.Parse(sessionid)
	if err != nil {
		ctx.Redirect("/login")
		return
	}
	user, err := controllers.GetUser(sessionUuid)
	if err != nil {
		// ctx.JSON(iris.Map{"code": 404, "message": "User not found"})
		// redirect to the login page if user is not authenticated

		ctx.Redirect("/login")
		return
	}
	ctx.Values().Set("user", user)

	ctx.Next()
}
