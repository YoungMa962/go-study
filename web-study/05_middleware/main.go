package main

import (
	"log"
	"net/http"
	"time"
)
import "gee"

func onlyForV2() gee.HandlerFunc {
	return func(context *gee.Context) {
		log.Printf("Time[%v],Group_V2:[Method:%s,Path:%s,Params%v]", time.Now(), context.Method(), context.Path(), context.Params())
		context.Next()
	}
}

func main() {
	r := gee.NewEngine()
	r.Use(gee.Logger())
	//r.GET("/index", func(c *gee.Context) {
	//	c.HTMLResponse(http.StatusOK, "<h1>Index Page</h1>")
	//})

	v2 := r.Group("/v2")
	v2.Use(func(context *gee.Context) {
		log.Printf("Time[%v],Group_V2:[Method:%s,Path:%s,Params%v]", time.Now().Format("2006/01/02 15:04:05"), context.Method(), context.Path(), context.Params())
		context.Next()
	})
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			c.StringResponse(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path())
		})
		v2.POST("/login", func(c *gee.Context) {
			c.JSONResponse(http.StatusOK, gee.H{
				"username": c.Form("username"),
				"password": c.Form("password"),
			})
		})
	}
	r.Run(":9999")

}
