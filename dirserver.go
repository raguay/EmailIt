package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func markdownTranslate(app *App, name string, dir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if strings.Contains(c.Request.URL.String(), ".md") {
			//
			// Get the path to the file.
			//
			newURL, _ := url.QueryUnescape(c.Request.URL.String())
			localFile := strings.Replace(newURL, name, "", 1)
			filePath := app.AppendPath(dir, localFile)
			bodydata := app.ReadFile(filePath)

			//
			// Translate the markdown to HTML.
			//
			htmlText := mdToHTML([]byte(bodydata))
			fmt.Printf("Translated text: %s\n", htmlText)

			//
			// send back the translated data.
			//
			c.Data(200, "text/html", htmlText)
			c.Next()
		}
	}
}

func dirserver(app *App, name string, dir string) {
	//
	// This will have the web server backend for EmailIt
	//
	r := gin.Default()
	r.Use(markdownTranslate(app, name, dir))
	r.Use(static.Serve("/"+name+"/", static.LocalFile(dir, true)))

	//
	// Run the server.  TODO: I need to add a new unique port for each server.
	//
	err := r.Run(":9000")
	if err != nil {
		fmt.Print("\n Server error:\n", err.Error())
	}
}
