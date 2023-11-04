package main

import (
	"fmt"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func markdownTranslate(app *App) gin.HandlerFunc {
	return func(c *gin.Context) {

		// before request

		// Perform the request.
		c.Next()

		// after request
		if strings.Contains(c.Request.URL.String(), ".md") {
			bodydata := app.ReadFile(app.AppendPath(app.GetHomeDir(), c.Request.URL.String()))
			c.Data(200, "text/html", []byte(bodydata))
		}
	}
}

func dirserver(app *App, dir string) {
	//
	// This will have the web server backend for EmailIt
	//
	r := gin.Default()
	r.Use(markdownTranslate(app))
	r.Use(static.Serve("/", static.LocalFile(dir, true)))

	//
	// Run the server.  TODO: I need to add a new unique port for each server.
	//
	err := r.Run(":9000")
	if err != nil {
		fmt.Print("\n Server error:\n", err.Error())
	}
}
