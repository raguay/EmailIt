package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

type NoteMsg struct {
	Message string `json:"note" binding:"required"`
}

type NoteChangeMsg struct {
	Note   string `json:"note" binding:"required"`
	Noteid string `json:"noteid" binding:"required"`
	Aw     string `json:"aw" binding:"required"`
}

type ScriptMsg struct {
	Script      string   `json:"script" binding:"required"`
	Text        string   `json:"text"`
	Env         string   `json:"env"`
	EnvVar      []string `json:"envVar"`
	CommandLine string   `json:"commandLine"`
	ReturnMsg   string   `json:"returnMsg"`
}

type TemplateMsg struct {
	Template  string `json:"template" binding:"required"`
	Text      string `json:"text"`
	ReturnMsg string `json:"returnMsg"`
}

type EmailMsg struct {
	Account   string `json:"account"`
	To        string `json:"to" binding:"required"`
	From      string `json:"from" binding:"required"`
	Subject   string `json:"subject" binding:"required"`
	Body      string `json:"body" binding:"required"`
	ReturnMsg string `json:"returnMsg"`
}

type CommandMsg struct {
	Action    string `json:"action"`
	ReturnMsg string `json:"returnMsg"`
}

type EmptyMsg struct {
	ReturnMsg string `json:"returnMsg"`
}

func backend(app *App, ctx context.Context) {
	//
	// This will have the web server backend for EmailIt
	//
	r := gin.Default()
	r.Use(gin.Recovery())

	//
	// Define the routes needed.
	//
	r.GET("/api/note/:noteid/:type", func(c *gin.Context) {
		noteid := c.Param("noteid")
		nid, _ := strconv.Atoi(noteid)
		hmdir := app.GetHomeDir()
		notesfilepath := filepath.Join(hmdir, ".config", "EmailIt", "notes.json")
		notesRaw, err := os.ReadFile(notesfilepath)
		if err != nil || nid < 0 || nid > 9 {
			c.JSON(http.StatusBadRequest, gin.H{
				"note": "no note",
			})
		} else {
			var notes []string
			err := json.Unmarshal(notesRaw, &notes)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			c.JSON(http.StatusOK, gin.H{
				"note": notes[nid],
			})
		}
	})

	r.PUT("/api/note/:noteid/:type", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "okay",
		})
		var json NoteMsg
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		noteid := c.Param("noteid")
		awtype := c.Param("type")
		nid, err := strconv.Atoi(noteid)

		//
		// Send it to the frontend.
		//
		if err != nil || (nid > 0 && nid < 9) {
			rt.EventsEmit(ctx, "notechange", NoteChangeMsg{
				Note:   json.Message,
				Noteid: noteid,
				Aw:     awtype,
			})
		}
	})

	r.GET("/api/script/env/list", func(c *gin.Context) {
		json := EmptyMsg{}
		//
		// Create return message name for this query.
		//
		umicro := time.Now().UnixMicro()
		json.ReturnMsg = fmt.Sprintf("%s%d", "EnvList", umicro)
		ch := make(chan int)

		rt.EventsEmit(ctx, "envList", json)
		rt.EventsOnce(ctx, json.ReturnMsg, func(optionalData ...interface{}) {
			c.JSON(http.StatusOK, optionalData[0])
			nwdata := 1
			ch <- nwdata
		})
		<-ch
		close(ch)
	})

	r.GET("/api/scripts/list", func(c *gin.Context) {
		json := EmptyMsg{}
		//
		// Create return message name for this query.
		//
		umicro := time.Now().UnixMicro()
		json.ReturnMsg = fmt.Sprintf("%s%d", "ScriptList", umicro)
		ch := make(chan int)

		rt.EventsEmit(ctx, "scriptList", json)
		rt.EventsOnce(ctx, json.ReturnMsg, func(optionalData ...interface{}) {
			c.JSON(http.StatusOK, optionalData[0])
			nwdata := 1
			ch <- nwdata
		})
		<-ch
		close(ch)
	})

	r.PUT("/api/script/run", func(c *gin.Context) {
		var json ScriptMsg
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//
		// Create return message name for this query.
		//
		umicro := time.Now().UnixMicro()
		json.ReturnMsg = fmt.Sprintf("%s%d", json.Script, umicro)
		ch := make(chan int)
		rt.EventsEmit(ctx, "scriptRun", json)
		rt.EventsOnce(ctx, json.ReturnMsg, func(optionalData ...interface{}) {
			c.JSON(http.StatusOK, optionalData[0])
			nwdata := 1
			ch <- nwdata
		})
		<-ch
		close(ch)
	})

	r.GET("/api/template/list", func(c *gin.Context) {
		//
		// Create return message name for this query.
		//
		json := EmptyMsg{}
		umicro := time.Now().UnixMicro()
		json.ReturnMsg = fmt.Sprintf("%s%d", "templatelist", umicro)
		ch := make(chan int)

		rt.EventsEmit(ctx, "templateList", json)
		rt.EventsOnce(ctx, json.ReturnMsg, func(optionalData ...interface{}) {
			c.JSON(http.StatusOK, optionalData[0])
			nwdata := 1
			ch <- nwdata
		})
		<-ch
		close(ch)
	})

	r.PUT("/api/template/run", func(c *gin.Context) {
		var json TemplateMsg
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//
		// Create return message name for this query.
		//
		umicro := time.Now().UnixMicro()
		json.ReturnMsg = strings.ReplaceAll(fmt.Sprintf("%s%d", json.Template, umicro), " ", "-")
		ch := make(chan int)

		rt.EventsEmit(ctx, "templateRun", json)
		rt.EventsOnce(ctx, json.ReturnMsg, func(optionalData ...interface{}) {
			c.JSON(http.StatusOK, optionalData[0])
			nwdata := 1
			ch <- nwdata
		})
		<-ch
		close(ch)
	})

	r.GET("/api/emailit/mailto", func(c *gin.Context) {
		to := c.DefaultQuery("to", "")
		subject := c.DefaultQuery("subject", "")
		body := c.DefaultQuery("body", "")

		email := EmailMsg{
			Account:   "Default",
			To:        to,
			From:      "",
			Subject:   subject,
			Body:      body,
			ReturnMsg: "",
		}
		//
		// Create return message name for this query.
		//
		umicro := time.Now().UnixMicro()
		email.ReturnMsg = strings.ReplaceAll(fmt.Sprintf("%s%d", to, umicro), " ", "-")
		ch := make(chan int)

		rt.EventsEmit(ctx, "EditEmail", email)
		rt.EventsOnce(ctx, email.ReturnMsg, func(optionalData ...interface{}) {
			c.JSON(http.StatusOK, optionalData[0])
			nwdata := 1
			ch <- nwdata
		})
		<-ch
		close(ch)
	})

	r.PUT("/api/emailit/send", func(c *gin.Context) {
		var json EmailMsg
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//
		// Create return message name for this query.
		//
		umicro := time.Now().UnixMicro()
		json.ReturnMsg = fmt.Sprintf("%s%d", "emailitSend", umicro)
		ch := make(chan int)

		rt.EventsEmit(ctx, "emailSend", json)
		rt.EventsOnce(ctx, json.ReturnMsg, func(optionalData ...interface{}) {
			c.JSON(http.StatusOK, optionalData[0])
			nwdata := 1
			ch <- nwdata
		})
		<-ch
		close(ch)
	})

  r.GET("/api/emailit/open", func(c *gin.Context) {
		//
		// Tell the main program to show the EmailIt program.
		//
    json := CommandMsg{ Action: "EmailIt" }
		ch := make(chan int)

    umicro := time.Now().UnixMicro()
		json.ReturnMsg = fmt.Sprintf("%s%d", "commandSend", umicro)
	
		rt.EventsEmit(ctx, "EmailIt", json)
		rt.EventsOnce(ctx, json.ReturnMsg, func(optionalData ...interface{}) {
			c.JSON(http.StatusOK, optionalData[0])
			nwdata := 1
			ch <- nwdata
		})
		<-ch
		close(ch)
	})

  r.GET("/api/notes/open", func(c *gin.Context) {
		//
		// Tell the main program to show the EmailIt program.
		//
    json := CommandMsg{ Action: "Notes" }
		ch := make(chan int)

    umicro := time.Now().UnixMicro()
		json.ReturnMsg = fmt.Sprintf("%s%d", "commandSend", umicro)
	
		rt.EventsEmit(ctx, "Notes", json)
		rt.EventsOnce(ctx, json.ReturnMsg, func(optionalData ...interface{}) {
			c.JSON(http.StatusOK, optionalData[0])
			nwdata := 1
			ch <- nwdata
		})
		<-ch
		close(ch)
	})

	r.GET("/api/scriptline/open", func(c *gin.Context) {
		//
		// Tell the main program to show the EmailIt program.
		//
    json := CommandMsg{ Action: "ScriptLine" }
		ch := make(chan int)

    umicro := time.Now().UnixMicro()
		json.ReturnMsg = fmt.Sprintf("%s%d", "commandSend", umicro)
	
		rt.EventsEmit(ctx, "ScriptLine", json)
		rt.EventsOnce(ctx, json.ReturnMsg, func(optionalData ...interface{}) {
			c.JSON(http.StatusOK, optionalData[0])
			nwdata := 1
			ch <- nwdata
		})
		<-ch
		close(ch)
	})

  //
	// Run the server.
	//
	err := r.Run(":9978")
	if err != nil {
		fmt.Print("\n Server error:\n", err.Error())
	}
}
