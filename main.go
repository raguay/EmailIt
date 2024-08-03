package main

import (
	"embed"
	"encoding/json"
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
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

var wailsapp = NewApp()

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

var cliemail HttpEmailMsg

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
	cliemail.Account = "Default" // Set the default value for account.
	app := &cli.App{
		Name:     "EmailIt",
		Usage:    "A program for sending emails, working with text, and scripts.",
		Version:  "v2.1.0",
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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "a",
				Value:       "",
				Usage:       "Address to send an email",
				Destination: &cliemail.To,
			},
			&cli.StringFlag{
				Name:        "t",
				Value:       "",
				Usage:       "Attachment to the email",
				Destination: &cliemail.Attachment,
			},
			&cli.StringFlag{
				Name:        "s",
				Value:       "",
				Usage:       "Subject for the email",
				Destination: &cliemail.Subject,
			},
			&cli.StringFlag{
				Name:        "b",
				Value:       "",
				Usage:       "Body of the email",
				Destination: &cliemail.Body,
			},
		},
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
				Usage:   "Create an email using a TUI",
				Action: func(cCtx *cli.Context) error {
					BuildEmail()
					return (nil)
				},
			},
			{
				Name:    "notes",
				Aliases: []string{"n"},
				Usage:   "Open the notes.",
				Action: func(cCtx *cli.Context) error {
					result := getRequest("http://localhost:9978/api/notes/open", strings.NewReader(""))
					fmt.Printf("\nThe server returned: %s\n", result[1:len(result)-1])
					return (nil)
				},
			},
			{
				Name:    "emailit",
				Aliases: []string{"e"},
				Usage:   "Open the EmailIt email sending application.",
				Action: func(cCtx *cli.Context) error {
					result := getRequest("http://localhost:9978/api/emailit/open", strings.NewReader(""))
					fmt.Printf("\nThe server returned: %s\n", result[1:len(result)-1])
					return (nil)
				},
			},
			{
				Name:    "scriptline",
				Aliases: []string{"sl"},
				Usage:   "Open the ScriptLine application.",
				Action: func(cCtx *cli.Context) error {
					result := getRequest("http://localhost:9978/api/scriptline/open", strings.NewReader(""))
					fmt.Printf("\nThe server returned: %s\n", result[1:len(result)-1])
					return (nil)
				},
			},
			{
				Name:    "sendemail",
				Aliases: []string{"se"},
				Usage:   "Send the email directly. No GUI or TUI.",
				Action: func(cCtx *cli.Context) error {
					//
					// Create and send the email. Then quit.
					//
					bodyJson := HttpEmailMsg{
						Account:    cliemail.Account,
						From:       "Default",
						To:         cliemail.To,
						Subject:    cliemail.Subject,
						Body:       cliemail.Body,
						Attachment: cliemail.Attachment,
					}
					body, err := json.Marshal(bodyJson)
					bodyStr := string(body[:])
					if err == nil {
						//
						// Send the email then!
						//
						SendHTTPQuery("PUT", "http://localhost:9978/api/emailit/send", bodyStr)
					}

					return (nil)
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runFrontEnd() {
	//
	// Run the different server locations.
	//
	go dirserver(wailsapp, "Archive", "/Users/raguay/Dropbox/Richard/Notes/TheArchive")

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
		Frameless:         true,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		Assets:            assets,
		LogLevel:          logger.DEBUG,
		OnStartup:         wailsapp.startup,
		OnDomReady:        wailsapp.domReady,
		OnShutdown:        wailsapp.shutdown,
		CSSDragProperty:   "--wails-draggable",
		CSSDragValue:      "drag",
		Bind: []interface{}{
			wailsapp,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		// macOS platform specific items.
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "EmailIt",
				Message: "Version 2.0.0 © 2022 Richard Guay <raguay@customct.com>",
				Icon:    icon,
			},
		},
		// Linux platform specific items.
		Linux: &linux.Options{
			Icon:                icon,
			WindowIsTranslucent: true,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyAlways,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
