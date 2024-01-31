package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Function:      BuildEmail()
//
// Description:   the main program loop for the TUI.
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
// description: The structure for the bubbletea interface for the email sending TUI.
type model struct {
	account    textinput.Model // This is for the account.
	to         textinput.Model // this is for the to address
	subject    textinput.Model // this is the subject line
	body       textarea.Model  // The body of the message.
	attachment textinput.Model // The attachments to the message
	focused    int             // This is the currently focused input
	err        string          // this will contain any errors from the validators
	list       list.Model      // list for choosing accounts
	sending    bool            // flag if sending or not.
}

// Type:          errMsg
//
// Description:   This is the error message type.
type (
	errMsg error
)

// Type:          HttpEmailMsg
//
// Description:   This is the structure for sending an email to the EmailIt program.
type HttpEmailMsg struct {
	Account    string `json:"account"`
	To         string `json:"to" binding:"required"`
	From       string `json:"from" binding:"required"`
	Subject    string `json:"subject" binding:"required"`
	Body       string `json:"body" binding:"required"`
	ReturnMsg  string `json:"returnMsg"`
	Attachment string `json:"attachment"`
}

type EmailAccounts struct {
	Default    bool   `json:"default" binding:"required"`
	FooterHTML string `json:"footerHTML" binding:"required"`
	From       string `json:"from" binding:"required"`
	HeaderHTML string `json:"headerHTML" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Password   string `json:"passwd" binding:"required"`
	Port       string `json:"port" binding:"required"`
	Signiture  string `json:"signiture" binding:"required"`
	SmtpServer string `json:"smtpserver" binding:"required"`
	Username   string `json:"usersname" binding:"required"`
}

// Constants:     These are the different constants for the currently focused field.
const (
	accountfield = iota
	tofield
	subjectfield
	bodyfield
	attachmentfield
	sendfield
)

// Variables:     These are the different styling options for the program.
var (
	labelStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#9580FF"))
	inputStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#9F70A9"))
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#80FFEA"))

	continueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F8F8F2"))
	textareaStyle = textarea.Style{
		Base:             lipgloss.NewStyle().Foreground(lipgloss.Color("#9580FF")).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#9580FF")),
		CursorLine:       lipgloss.NewStyle().Foreground(lipgloss.Color("#80FFEA")),
		CursorLineNumber: lipgloss.NewStyle().Foreground(lipgloss.Color("#80FFEA")),
		EndOfBuffer:      lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF80")),
		LineNumber:       lipgloss.NewStyle().Foreground(lipgloss.Color("#544158")),
		Placeholder:      lipgloss.NewStyle().Foreground(lipgloss.Color("#9580FF")),
		Prompt:           lipgloss.NewStyle().Foreground(lipgloss.Color("#9580FF")),
		Text:             lipgloss.NewStyle().Foreground(lipgloss.Color("#9F70A9")),
	}
)

// These items are for the accounts list.

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("#9580FF"))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("#80FFEA"))
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#9F70A9"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

//
// End of list items
//

// Function:          nameValidator
//
// Description:       This is a validator functions to ensure valid input.
//
// Inputs:
//
//	s       string to validate.
func nameValidator(s string) error {
	return nil
}

// Function:          stringValidator
//
// Description:       This is a validator functions to ensure valid input.
//
// Inputs:
//
//	s       string to validate.
func stringValidator(s string) error {
	return nil
}

// Function:        initialModel()
//
// Description:     This will create and initialize the model for our TUI.
func initialModel() model {
	var m model
	var b App
	var accounts []EmailAccounts

	//
	// Set the different widths. Should be based on terminal size and will
	// be replaced when the terminal size message is sent.
	//
	inputWidth := 92
	textareaWidth := 100

	//
	// Get the list of accounts.
	//
	hmdir := b.GetHomeDir()
	accountsPath := b.AppendPath(hmdir, ".config/EmailIt/emailaccounts.json")
	items := []list.Item{}
	numacc := 0
	if b.FileExists(accountsPath) {
		accountsText := b.ReadFile(accountsPath)
		err := json.Unmarshal([]byte(accountsText), &accounts)
		if err != nil {
			fmt.Printf("Couldn't parse the JSON accounts file! %s", err.Error())
			os.Exit(0)
		} else {
			numacc = len(accounts) + 1
			items = make([]list.Item, numacc)
			items[0] = item("Default")
			for i := 1; i < numacc; i++ {
				items[i] = item(accounts[i-1].Name)
			}
		}
	} else {
		fmt.Print("\nEmailIt doesn't have accounts setup. Please setup accounts before using the TUI.\n")
		os.Exit(0)
	}

	l := list.New(items, itemDelegate{}, 50, numacc+8)
	l.Title = "Which account to you want to use?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	m.list = l

	//
	// Create the different inputs.
	//
	m.to = textinput.New()
	m.to.Placeholder = ""
	m.to.CharLimit = 0
	m.to.Prompt = ""
	m.to.Validate = nameValidator
	m.to.TextStyle = inputStyle
	m.to.Cursor.Style = cursorStyle
	m.to.Width = inputWidth
	if cliemail.To != "" {
		m.to.SetValue(cliemail.To)
	}

	m.subject = textinput.New()
	m.subject.Placeholder = ""
	m.subject.CharLimit = 0
	m.subject.Prompt = ""
	m.subject.Validate = nameValidator
	m.subject.TextStyle = inputStyle
	m.subject.Cursor.Style = cursorStyle
	m.subject.Width = inputWidth
	if cliemail.Subject != "" {
		m.subject.SetValue(cliemail.Subject)
	}

	m.account = textinput.New()
	m.account.Placeholder = "                                    "
	m.account.CharLimit = 0
	m.account.Prompt = ""
	m.account.Validate = stringValidator
	m.account.TextStyle = inputStyle
	m.account.Cursor.Style = cursorStyle
	m.account.Width = inputWidth

	m.body = textarea.New()
	m.body.FocusedStyle = textareaStyle
	m.body.BlurredStyle = textareaStyle
	m.body.Prompt = ""
	m.body.SetWidth(textareaWidth)
	if cliemail.Body != "" {
		m.body.SetValue(cliemail.Body)
	}

	m.attachment = textinput.New()
	m.attachment.Placeholder = "                                    "
	m.attachment.CharLimit = 0
	m.attachment.Prompt = ""
	m.attachment.Validate = stringValidator
	m.attachment.TextStyle = inputStyle
	m.attachment.Cursor.Style = cursorStyle
	m.attachment.Width = inputWidth

	//
	// Set up the rest of the default values.
	//
	m.account.SetValue("Default                ")
	m.focused = tofield
	m.sending = false
	m.err = ""
	return m
}

// Function:      Init()
//
// Description:   This initializes the TUI.
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Function:      nextInput()
//
// Description:   focuses the next input field
func (m *model) nextInput() {
	//
	// Increment the focused item and wrap around if
	// too large.
	//
	m.focused = (m.focused + 1) % (sendfield + 1)
}

// Function:      prevInput()
//
// Description:   focuses the previous input field with proper underflow.
func (m *model) prevInput() {
	//
	// Decrement the focused item.
	//
	m.focused--

	//
	// If less than zero, wrap around to the subject field.
	//
	if m.focused < 0 {
		m.focused = subjectfield
	}
}

// Function:      SendMessage()
//
// Description:   This is the tea message to actually send the email.
func (m model) SendMessage() tea.Msg {
	//
	// Create and send the email. Then quit.
	//
	fixacct := strings.Trim(m.account.Value(), " ")
	fixto := strings.Trim(m.to.Value(), " ")
	bodyJson := HttpEmailMsg{
		Account:    fixacct,
		From:       "Default",
		To:         fixto,
		Subject:    m.subject.Value(),
		Body:       m.body.Value(),
		Attachment: m.attachment.Value(),
	}
	body, err := json.Marshal(bodyJson)
	bodyStr := string(body[:])
	if err != nil {
		m.err = err.Error()
	} else {
		//
		// Send the email then!
		//
		result := SendHTTPQuery("PUT", "http://localhost:9978/api/emailit/send", bodyStr)
		if result != "" {
			fmt.Print(labelStyle.Width(25 + len(result)).Render(fmt.Sprintf("The server returned: %s", result)))
		}
	}
	return tea.QuitMsg{}
}

// Function:      Update()
//
// Description:   This is the bubbletea function to update the graphics based
//
//	on the different messages received.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, 5)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == sendfield {
				//
				// This is the last input, save the inputs
				//
				if !m.sending {
					m.sending = true
					return m, m.SendMessage
				} else {
					return m, nil
				}
			} else if m.focused == accountfield {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.account.SetValue(string(i))
					m.nextInput()
				}
			} else if m.focused < sendfield {
				m.nextInput()
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			if m.focused != accountfield {
				m.prevInput()
			}
		case tea.KeyTab, tea.KeyCtrlN:
			if m.focused != accountfield {
				m.nextInput()
			}
		}

		//
		// This case is for getting the terminal size and setting up the
		// component sizes.
		//
	case tea.WindowSizeMsg:
		m.to.Width = msg.Width - 18
		m.account.Width = msg.Width - 18
		m.attachment.Width = msg.Width - 18
		m.body.SetWidth(msg.Width - 10)
		m.body.SetHeight(msg.Height - 10)
		m.list.SetWidth(msg.Width - 10)

	// We handle errors just like any other message
	case errMsg:
		m.err = msg.Error()
		return m, nil
	}

	//
	// Focus the current input.
	//
	switch m.focused {
	case tofield:
		m.to.Focus()
		m.account.Blur()
		m.subject.Blur()
		m.body.Blur()
		m.attachment.Blur()

	case accountfield:
		m.account.Focus()
		m.to.Blur()
		m.subject.Blur()
		m.body.Blur()
		m.attachment.Blur()

	case subjectfield:
		m.subject.Focus()
		m.to.Blur()
		m.account.Blur()
		m.body.Blur()
		m.attachment.Blur()

	case bodyfield:
		m.body.Focus()
		m.to.Blur()
		m.account.Blur()
		m.subject.Blur()
		m.attachment.Blur()

	case attachmentfield:
		m.body.Blur()
		m.to.Blur()
		m.account.Blur()
		m.subject.Blur()
		m.attachment.Focus()
	}

	//
	// Send the update message.
	//
	m.to, cmds[tofield] = m.to.Update(msg)
	m.subject, cmds[subjectfield] = m.subject.Update(msg)
	m.body, cmds[bodyfield] = m.body.Update(msg)
	m.attachment, cmds[attachmentfield] = m.attachment.Update(msg)
	if m.focused == accountfield {
		m.list, cmds[accountfield] = m.list.Update(msg)
	} else {
		m.account, cmds[accountfield] = m.account.Update(msg)
	}
	return m, tea.Batch(cmds...)
}

// Function:    View
//
// Description: The view on a model controls how it is displayed. It returns strings
//
//	for displaying to the user.
func (m model) View() string {
	//
	// Create the actual view string.
	//
	result := ""
	if m.focused == accountfield {
		result = m.list.View()
	} else if m.sending {
		result = labelStyle.Width(25).Render("The Email is sending....")
	} else {
		result = fmt.Sprintf(`%s  %s
%s  %s
%s  %s
%s
%s  %s`,
			labelStyle.Width(7).Render("Account"),
			m.account.View(),
			labelStyle.Width(7).Render("     To"),
			m.to.View(),
			labelStyle.Width(7).Render("Subject"),
			m.subject.View(),
			m.body.View(),
			labelStyle.Width(12).Render("Attachments"),
			m.attachment.View())
		//
		// Add the Send Email field.
		//
		result += continueStyle.Render("\n")
		if m.focused == sendfield {
			result += continueStyle.Render("Send Email -> ")
		}
	}
	result += "\n"

	//
	// return the result.
	//
	return result
}

// Function:        SendHTTPQuery()
//
// Description:     Send a HTTP request
//
// Inputs:
//
//	method      The HTTP method to use
//	uri         The URI to send the request.
//	body        The body of the request. The body is assumed to be JSON structure.
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
				return fmt.Sprintf("error: calling %s returned empty response", uri)
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
