package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"strconv"
	"time"

	clip "github.com/atotto/clipboard"
	cp "github.com/otiai10/copy"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	mail "gopkg.in/mail.v2"
)

// App application struct and other structs
type App struct {
	ctx  context.Context
	err  string
	proc *os.Process
}

type FileParts struct {
	Dir       string
	Name      string
	Extension string
}

type FileInfo struct {
	Dir       string
	Name      string
	Extension string
	IsDir     bool
	Size      int64
	Modtime   string
}

type Environment struct {
	Name   string   `json:"name" binding:"required"`
	EnvVar []string `json:"envVar" binding:"required"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (b *App) startup(ctx context.Context) {
	//
	// Save the context for other functions.
	//
	b.ctx = ctx

	//
	// We need to start the backend and setup the signaling.
	//
	go backend(b, ctx)
}

// domReady is called after the front-end dom has been loaded
func (b *App) domReady(ctx context.Context) {
}

// shutdown is called at application termination
func (b *App) shutdown(ctx context.Context) {
}

func (b *App) SystemOpenFile(prog string) {
	if b.FileExists(prog) {
		var ArgsArray [2]string
		var sar []string
		ArgsArray[1] = prog
		b.RunCommandLine("/usr/bin/open", ArgsArray[:], sar, "")
	}
}

func (b *App) GetExecutable() string {
	ex, err := os.Executable()
	if err != nil {
		b.err = err.Error()
	}
	return ex
}

func (b *App) Getwd() string {
	wd, err := os.Getwd()
	if err != nil {
		b.err = err.Error()
	}
	return wd
}

func (b *App) ReadFile(path string) string {
	b.err = ""
	contents, err := os.ReadFile(path)
	if err != nil {
		b.err = err.Error()
	}
	return string(contents[:])
}

func (b *App) GetHomeDir() string {
	b.err = ""
	hdir, err := os.UserHomeDir()
	if err != nil {
		b.err = err.Error()
	}
	return hdir
}

func (b *App) WriteFile(path string, data string) {
	err := os.WriteFile(path, []byte(data), 0666)
	if err != nil {
		b.err = err.Error()
	}
}

func (b *App) FileExists(path string) bool {
	b.err = ""
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func (b *App) DirExists(path string) bool {
	b.err = ""
	dstat, err := os.Stat(path)
	if err != nil {
		b.err = err.Error()
		return false
	}
	return dstat.IsDir()
}

func (b *App) SplitFile(path string) FileParts {
	b.err = ""
	var parts FileParts
	parts.Dir, parts.Name = filepath.Split(path)
	parts.Extension = filepath.Ext(path)
	return parts
}

func (b *App) ReadDir(path string) []FileInfo {
	b.err = ""
	var result []FileInfo
	result = make([]FileInfo, 0, 0)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		b.err = err.Error()
	} else {
		for _, file := range files {
			var fileInfo FileInfo
			fileInfo.Name = file.Name()
			fileInfo.Size = file.Size()
			fileInfo.IsDir = file.IsDir()
			fileInfo.Modtime = file.ModTime().Format(time.ANSIC)
			fileInfo.Dir = path
			fileInfo.Extension = filepath.Ext(file.Name())
			result = append(result, fileInfo)
		}
	}
	return result
}

func (b *App) MakeDir(path string) {
	b.err = ""
	err := os.MkdirAll(path, 0755)
	if err != nil {
		b.err = err.Error()
	}
}

func (b *App) MakeFile(path string) {
	b.err = ""
	b.WriteFile(path, "")
}

func (b *App) MoveEntries(from string, to string) {
	b.err = ""
	err := os.Rename(from, to)
	if err != nil {
		b.err = err.Error()
	}
}

func (b *App) RenameEntry(from string, to string) {
	b.err = ""
	err := os.Rename(from, to)
	if err != nil {
		b.err = err.Error()
	}
}

func (b *App) GetError() string {
	return b.err
}

func (b *App) CopyEntries(from string, to string) {
	b.err = ""
	info, err := os.Stat(from)
	if os.IsNotExist(err) {
		b.err = err.Error()
		return
	}
	if info.IsDir() {
		//
		// It's a directory! Do a deap copy.
		//
		err := cp.Copy(from, to)
		if err != nil {
			b.err = err.Error()
			return
		}
	} else {
		//
		// It's a file. Just copy it.
		//
		source, err := os.Open(from)
		if err != nil {
			b.err = err.Error()
			return
		}
		defer source.Close()

		destination, err := os.Create(to)
		if err != nil {
			b.err = err.Error()
			return
		}
		defer destination.Close()
		_, err = io.Copy(destination, source)

		if err != nil {
			b.err = err.Error()
		}
	}
}

func (b *App) DeleteEntries(path string) {
	b.err = ""
	err := os.RemoveAll(path)
	if err != nil {
		b.err = err.Error()
	}
}

func (b *App) RunCommandLine(cmd string, args []string, env []string, dir string) string {
	b.err = ""
	cmdline := exec.Command(cmd)
	cmdline.Args = args
	cmdline.Env = env
	cmdline.Dir = dir
	result, err := cmdline.CombinedOutput()
	if err != nil {
		b.err = err.Error()
	}

	return string(result[:])
}

func (b *App) GetClip() string {
	result, err := clip.ReadAll()
	if err != nil {
		b.err = err.Error()
	}
	return result
}

func (b *App) SetClip(msg string) {
	err := clip.WriteAll(msg)
	if err != nil {
		b.err = err.Error()
	}
}

func (b *App) GetEnvironment() []string {
	return os.Environ()
}

func (b *App) GetEnv(name string) string {
	return os.Getenv(name)
}

func (b *App) AppendPath(dir string, name string) string {
	return filepath.Join(dir, name)
}

func (b *App) Quit() {
	wailsruntime.Quit(b.ctx)
}

func (b *App) GetOSName() string {
	os := goruntime.GOOS
	result := ""
	switch os {
	case "windows":
		result = "windows"
		break
	case "darwin":
		result = "macos"
	case "linux":
		result = "linux"
	default:
		result = fmt.Sprintf("%s", os)
	}
	return result
}

func (b *App) GetCopyClip(name string) string {
	return "Not Implemented Yet"
}

func (b *App) GetFeedback(question string, defans string) string {
	return "Not Implemented Yet"
}

func (b *App) GetGitHubThemes() []string {
	var result []string
	return result
}

func (b *App) SendEmail(username string, from string, password string, host string, port string, toList string, msg string, msgText string, subject string) string {
	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", toList)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", msg)
	m.AddAlternative("text/plain", msgText)
	iport, _ := strconv.Atoi(port)
	d := mail.NewDialer(host, iport, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		b.err = err.Error()
		return b.err
	}
	return "Success"
}
