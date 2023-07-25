package main

import (
	"bytes"
	"encoding/json"
	// "encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func BuildEmail() error {
	//
	// create the Bubbletea interface for building an email.
	//
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	return nil
}

// Struct:		model
//
// description: The structure for the bubbletea interface for building a dialog.
type model struct {
	inputs  []textinput.Model // This contains the input fields for the labels
	focused int               // This is the currently focused input
	err     error             // this will contain any errors from the validators
}

type (
	tickMsg struct{}
	errMsg  error
)

type HttpEmailMsg struct {
	Account   string `json:"account"`
	To        string `json:"to" binding:"required"`
	From      string `json:"from" binding:"required"`
	Subject   string `json:"subject" binding:"required"`
	Body      string `json:"body" binding:"required"`
	ReturnMsg string `json:"returnMsg"`
}

const (
	tofield = iota
	fromfield
	accountfield
	subjectfield
	bodyfield
)

const (
	purple   = lipgloss.Color("#9580FF")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(purple)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

// Validator functions to ensure valid input
func nameValidator(s string) error {
	return nil
}

func stringValidator(s string) error {
	return nil
}

func initialModel() model {
	var inputs []textinput.Model = make([]textinput.Model, 5)

	inputs[tofield] = textinput.New()
	inputs[tofield].Placeholder = ""
	inputs[tofield].CharLimit = 100
	inputs[tofield].Width = 102
	inputs[tofield].Prompt = ""
	inputs[tofield].Validate = nameValidator

	inputs[fromfield] = textinput.New()
	inputs[fromfield].Placeholder = ""
	inputs[fromfield].CharLimit = 100
	inputs[fromfield].Width = 102
	inputs[fromfield].Prompt = ""
	inputs[fromfield].Validate = nameValidator

	inputs[subjectfield] = textinput.New()
	inputs[subjectfield].Placeholder = ""
	inputs[subjectfield].CharLimit = 100
	inputs[subjectfield].Width = 102
	inputs[subjectfield].Prompt = ""
	inputs[subjectfield].Validate = nameValidator

	inputs[accountfield] = textinput.New()
	inputs[accountfield].Placeholder = ""
	inputs[accountfield].CharLimit = 100
	inputs[accountfield].Width = 102
	inputs[accountfield].Prompt = ""
	inputs[accountfield].Validate = stringValidator

	inputs[bodyfield] = textinput.New()
	inputs[bodyfield].Placeholder = ""
	inputs[bodyfield].CharLimit = 100
	inputs[bodyfield].Width = 102
	inputs[bodyfield].Prompt = ""
	inputs[bodyfield].Validate = nameValidator

	return model{
		// Our list of acctions
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// nextInput focuses the next input field
func (m *model) nextInput() {
	//
	// Increment the focused item and wrap around if
	// too large.
	//
	m.focused = (m.focused + 1) % 6
}

// prevInput focuses the previous input field
func (m *model) prevInput() {
	//
	// Decrement the focused item.
	//
	m.focused--

	//
	// If less than zero, wrap around to the highest number.
	//
	if m.focused < 0 {
		m.focused = 4
	}
}

func (m model) SendMessage() tea.Msg {
	//
	// Create and send the email. Then quit.
	//
	bodyJson := HttpEmailMsg{
		Account: m.inputs[accountfield].Value(),
		From:    m.inputs[fromfield].Value(),
		To:      m.inputs[tofield].Value(),
		Subject: m.inputs[subjectfield].Value(),
		Body:    m.inputs[bodyfield].Value(),
	}
	body, err := json.Marshal(bodyJson)
	bodyStr := string(body[:])
	if err != nil {
	} else {
		result := SendHTTPQuery("PUT", "http://localhost:9978/api/emailit/send", bodyStr)
	}
	return tea.QuitMsg{}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, 5)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == 5 {
				//
				// This is the last input, save the inputs
				//
				return m, m.SendMessage
			} else {
				m.nextInput()
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := 0; i < 5; i++ {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func viewEmail(m model) string {
	result := fmt.Sprintf(
		` 
 %s  %s
 %s  %s
 %s  %s
 %s  %s
 %s  %s
 `,
		inputStyle.Width(8).Render("     To"),
		m.inputs[tofield].View(),
		inputStyle.Width(8).Render("   From"),
		m.inputs[fromfield].View(),
		inputStyle.Width(8).Render("Account"),
		m.inputs[accountfield].View(),
		inputStyle.Width(8).Render("Subject"),
		m.inputs[subjectfield].View(),
		inputStyle.Width(8).Render("   Body"),
		m.inputs[bodyfield].View())
	if m.focused == 5 {
		result += continueStyle.Render("Send Email ->")
	}
	result += "\n"
	return result
}

// Function:    View
//
// Description: The view on a model controls how it is displayed. It returns strings
//
//	for displaying to the user.
func (m model) View() string {

	//
	// Make sure the inputs are not selected.
	//
	for i := 0; i < 5; i++ {
		m.inputs[i].Blur()
	}

	//
	// Focus the current input.
	//
	if m.focused != 5 {
		m.inputs[m.focused].Focus()
	}

	//
	// Generate the view and return the result.
	//
	return viewEmail(m)
}

func SendHTTPQuery(method string, uri string, body string) string {
	var result string
	switch strings.Trim(method, " ") {
	case "GET":
		{
			resp, err := http.Get(uri)
			if err != nil {
				print(err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				print(err)
			}
			result = string(body)
		}

	case "POST":
		{
			resp, err := http.Post(uri, "application/json", bytes.NewBuffer([]byte(body+"\n")))
			if err != nil {
				print(err)
			}
			defer resp.Body.Close()
			repBody, errdec := ioutil.ReadAll(resp.Body)
			if errdec != nil {
				print(errdec)
			}
			result = string(repBody)
		}
	case "PUT":
		{
			client := http.Client{}

			// regardless of GET or POST, we can safely add the body
			req, err := http.NewRequest("PUT", uri, bytes.NewBuffer([]byte(body+"\n")))
			if err != nil {
				return err.Error()
			}
			headers := map[string]string{}
			headers["Method"] = "PUT"
			headers["Content-Type"] = "application/json"

			// for each header passed, add the header value to the request
			for k, v := range headers {
				req.Header.Set(k, v)
			}

			// finally, do the request
			res, err := client.Do(req)
			if err != nil {
				return err.Error()
			}

			if res == nil {
				return fmt.Sprint("error: calling %s returned empty response", uri)
			}

			responseData, err := io.ReadAll(res.Body)
			if err != nil {
				return err.Error()
			}

			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				return fmt.Sprintf("error calling %s:\nstatus: %s\nresponseData: %s", uri, res.Status, responseData)
			}

			result = string(responseData)
		}
	}
	return result
}
