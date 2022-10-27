package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
	Script string `json:"script" binding:"required"`
	Text   string `json:"text"`
}

type TemplateMsg struct {
	Template string `json:"template" binding:"required"`
	Text     string `json:"text"`
}

type EmailMsg struct {
	Account string `json:"account"`
	To      string `json:"to" binding:"required"`
	From    string `json:"from" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

type EmptyMsg struct{}

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
			c.JSON(http.StatusOK, gin.H{
				"note": "no note",
			})
		} else {
			var notes []string
			json.Unmarshal(notesRaw, &notes)
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
		rt.EventsEmit(ctx, "envList", EmptyMsg{})
		running := true
		rt.EventsOnce(ctx, "envListReturn", func(optionalData ...interface{}) {
			c.JSON(http.StatusOK, optionalData[0])
			running = false
		})
		for running {
			time.Sleep(time.Millisecond)
		}
	})

	r.GET("/api/scripts/list", func(c *gin.Context) {
		rt.EventsEmit(ctx, "scriptList", EmptyMsg{})
		running := true
		rt.EventsOnce(ctx, "scriptListReturn", func(optionalData ...interface{}) {
			c.JSON(http.StatusOK, optionalData[0])
			running = false
		})
		for running {
			time.Sleep(time.Millisecond)
		}
	})

	r.PUT("/api/script/run", func(c *gin.Context) {
		var json ScriptMsg
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		rt.EventsEmit(ctx, "scriptRun", json)
		running := true
		rt.EventsOnce(ctx, "scriptRunReturn", func(optionalData ...interface{}) {
			c.JSON(http.StatusOK, optionalData[0])
			running = false
		})
		for running {
			time.Sleep(time.Millisecond)
		}
	})

	r.GET("/api/template/list", func(c *gin.Context) {
		rt.EventsEmit(ctx, "templateList", EmptyMsg{})
		running := true
		rt.EventsOnce(ctx, "templateListReturn", func(optionalData ...interface{}) {
			c.JSON(http.StatusOK, optionalData[0])
			running = false
		})
		for running {
			time.Sleep(time.Millisecond)
		}
	})

	r.PUT("/api/template/run", func(c *gin.Context) {
		var json TemplateMsg
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		rt.EventsEmit(ctx, "templateRun", json)
		running := true
		rt.EventsOnce(ctx, "templateRunReturn", func(optionalData ...interface{}) {
			c.JSON(http.StatusOK, optionalData[0])
			running = false
		})
		for running {
			time.Sleep(time.Millisecond)
		}
	})

	r.PUT("/api/emailit/send", func(c *gin.Context) {
		var json EmailMsg
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		rt.EventsEmit(ctx, "emailSend", json)
		running := true
		rt.EventsOnce(ctx, "emailSendReturn", func(optionalData ...interface{}) {
			c.JSON(http.StatusOK, optionalData[0])
			running = false
		})
		for running {
			time.Sleep(time.Millisecond)
		}
	})

	//
	// Run the server.
	//
	r.Run(":9978")
}
