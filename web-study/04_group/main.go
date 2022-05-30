package main

import (
	"gee"
	"net/http"
)

func main() {

	r := gee.NewEngine()
	r.GET("/index", func(c *gee.Context) {
		c.HTMLResponse(http.StatusOK, "<h1>Index Page</h1>")
	})

	v1 := r.Group("/v1")
	v1.GET("/", func(c *gee.Context) {
		c.HTMLResponse(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	v1.GET("/hello", func(c *gee.Context) {
		c.StringResponse(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path())
	})

	v2 := r.Group("/v2")
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
