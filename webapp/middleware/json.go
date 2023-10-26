package middleware

import "github.com/kataras/iris/v12"

func Json(ctx iris.Context) {
	ctx.Header("Content-Type", "application/json")
	ctx.Next()
}
