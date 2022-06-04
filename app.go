package main

import (
	"context"
)

// App application struct and other structs
type App struct {
	ctx context.Context
	err string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// domReady is called after the front-end dom has been loaded
func (a *App) domReady(ctx context.Context) {
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
}
