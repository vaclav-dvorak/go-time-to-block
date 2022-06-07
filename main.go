package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/sirupsen/logrus"
)

var (
	errStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9"))
	blockStyle = lipgloss.NewStyle().
			Width(13).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("63"))
)

type updateMsg struct{ resp resp }
type errMsg struct{ err error }

type resp struct {
	Top    int
	Bottom int
	Middle float64
}

type model struct {
	textInput textinput.Model
	spinner   spinner.Model
	date      time.Time
	err       error
	res       resp
	updating  bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "DD.MM.YYYY HH:II+GMT"
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 20

	s := spinner.NewModel()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("209"))

	return model{
		textInput: ti,
		spinner:   s,
		date:      time.Time{},
		err:       nil,
		res:       resp{},
		updating:  false,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			cmds = append(cmds, tea.Quit)
		case "enter":
			m.err = nil
			date, err := time.Parse("02.01.2006 15:04+MST", m.textInput.Value())
			if err != nil {
				m.err = err
				break
			}
			m.date = date
			m.updating = true
			cmds = append(cmds, m.spinner.Tick, updateTime(m))
		}

	case errMsg:
		m.updating = false
		m.err = msg.err
		m.res = resp{}

	case updateMsg:
		m.updating = false
		m.res = msg.resp
		m.err = nil

	case spinner.TickMsg:
		if m.updating {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() (ret string) {
	ret = ""
	ret += fmt.Sprintf("Input date you want to convert?\n\n%s", m.textInput.View())
	if m.err != nil {
		ret += fmt.Sprintf("\n\n\n%s\n", errStyle.Render(fmt.Sprintf("✘ - %s", m.err)))
	} else if m.updating {
		ret += fmt.Sprintf("\n\n\nConverting...%s\n", m.spinner.View())
	} else {
		ret += fmt.Sprintf("\n\n%s", renderResp(m.res))
	}
	ret += fmt.Sprintf("\n\n%s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("esc/q to quit"))
	return
}

func main() {
	p := tea.NewProgram(initialModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

func renderResp(r resp) string {
	if r.Top == 0 {
		return "\n\n"
	}
	return fmt.Sprintf("👆 After block:     %s\n👉 Exact blocktime: %s\n👇 Before block:    %s", blockStyle.Render(fmt.Sprintf("%d", r.Top)), blockStyle.Render(fmt.Sprintf("%.2f", r.Middle)), blockStyle.Render(fmt.Sprintf("%d", r.Bottom)))
}

func updateTime(m model) tea.Cmd {
	return func() tea.Msg {
		resp, err := getBlockData(m.date)
		if err != nil {
			return errMsg{err}
		}
		return updateMsg{resp}
	}
}
