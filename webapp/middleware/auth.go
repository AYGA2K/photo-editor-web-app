package middleware

import (
	"github.com/AYGA2K/photo-editor-web-app/webapp/controllers"
	"github.com/AYGA2K/photo-editor-web-app/webapp/models"
	"github.com/kataras/iris/v12"
)

func Authenticated(ctx iris.Context) {
	json := new(models.Session)
	if err := ctx.ReadJSON(json); err != nil {
		ctx.JSON(iris.Map{"code": 400, "message": "Invalid JSON"})
		return
	}
	user, err := controllers.GetUser(json.Sessionid)
	if err != nil {
		ctx.JSON(iris.Map{"code": 404, "message": "User not found"})
		return
	}
	ctx.Values().Set("user", user)
	ctx.Next()
}
