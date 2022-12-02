package main

import (
	"embed"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

// Function:     getRequest
//
// Description:  This method will issue a get request with the data sent
//
//	as json in the body.
//
// Inputs:
//
//	url        The url to send the request
//	data       An io.Reader pointing to a json string
func getRequest(url string, data io.Reader) string {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, data)
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err2 := client.Do(req)
	if err2 != nil {
		//
		// handle error. Most likely not running.
		//
		log.Fatal(err2)
	}
	body, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		log.Fatal(err3)
	}
	resp.Body.Close()
	return string(body)
}

func main() {
	//
	// Process the Command Line.
	//
	app := &cli.App{
		Name:     "EmailIt",
		Usage:    "A program for sending emails and working with text and scripts.",
		Version:  "v2.0.1",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Richard Guay",
				Email: "raguay@customct.com",
			},
		},
		Copyright: "(c) 2022 Richard Guay",
		HelpName:  "EmailIt",
		UsageText: "",
		Action: func(cCtx *cli.Context) error {
			if cCtx.Args().Len() == 0 {
				//
				// Only run the main if there isn't any arguments. Miss typing an argument will run here.
				//
				runFrontEnd()
			} else {
				//
				// We have a mailto protocal to process.
				//
				mailto := cCtx.Args().First()
				mailto = mailto[7:]                           // remove the 'mailto:'
				mailto = strings.ReplaceAll(mailto, "?", "&") // Change al ? To & - There should only be one.
				urimsg := fmt.Sprintf("http://localhost:9978/api/emailit/mailto?to=%s", mailto)
				result := getRequest(urimsg, strings.NewReader(""))
				fmt.Printf("\nThe server returned: %s\n", result[1:len(result)-1])
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "mkemail",
				Aliases: []string{"me"},
				Usage:   "Create an email",
				Action: func(cCtx *cli.Context) error {
					return BuildEmail()
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runFrontEnd() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "EmailIt",
		Width:             1022,
		Height:            608,
		MinWidth:          1022,
		MinHeight:         608,
		MaxWidth:          1022,
		MaxHeight:         608,
		DisableResize:     true,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 33, G: 37, B: 43, A: 255},
		Assets:            assets,
		LogLevel:          logger.DEBUG,
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		OnShutdown:        app.shutdown,
		CSSDragProperty:   "--wails-draggable",
		CSSDragValue:      "drag",
		Bind: []interface{}{
			app,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "EmailIt",
				Message: "Version 2.0.0 © 2022 Richard Guay <raguay@customct.com>",
				Icon:    icon,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
