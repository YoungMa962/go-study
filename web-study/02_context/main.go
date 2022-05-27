package main

import (
	"gee"
	"net/http"
)

func main() {
	engine := gee.NewEngine()
	engine.GET("/", func(context *gee.Context) {
		context.HTMLResponse(http.StatusOK, "<h1>Hello Word</h1>")
	})
	engine.GET("/hello", func(context *gee.Context) {
		context.StringResponse(http.StatusOK, "hello %s, you're at %s\\n", context.Query("name"), context.Path())
	})
	engine.POST("/login", func(context *gee.Context) {
		context.JSONResponse(http.StatusOK, gee.H{
			"username": context.Form("username"),
			"password": context.Form("password")},
		)
	})
	engine.POST("/json/login", func(context *gee.Context) {
		context.JSONResponse(http.StatusOK, context.Body())
	})

	engine.Run(":9999")
}
