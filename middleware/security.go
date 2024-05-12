package middleware

import (
	"github.com/kataras/iris/v12"
)

func Security(ctx iris.Context) {
	ctx.Header("X-XSS-Protection", "1; mode=block")
	ctx.Header("X-Content-Type-Options", "nosniff")
	ctx.Header("X-Download-Options", "noopen")
	ctx.Header("Strict-Transport-Security", "max-age=5184000")
	ctx.Header("X-Frame-Options", "DENY")
	ctx.Header("X-DNS-Prefetch-Control", "off")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
	ctx.Header("Content-Security-Policy", "default-src https:")

	ctx.Next()
}
