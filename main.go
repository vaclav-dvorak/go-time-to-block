package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/sirupsen/logrus"
)

const inputDate = "22.02.2022 22:22+GMT"

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

type resp struct {
	Top    int
	Bottom int
	Middle float64
}

type model struct {
	textInput textinput.Model
	err       error
	res       resp
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "DD.MM.YYYY HH:II+GMT"
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
		res:       resp{},
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.err = nil
			date, err := time.Parse("02.01.2006 15:04+MST", m.textInput.Value())
			if err != nil {
				m.err = err
				return m, nil
			}
			r, err := getBlockData(date)
			if err != nil {
				m.err = err
				return m, nil
			}
			m.res = r
			return m, nil
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() (ret string) {
	ret = ""
	ret += fmt.Sprintf("Input date you want to convert?\n\n%s", m.textInput.View())
	if m.err != nil {
		ret += fmt.Sprintf("\n\n\n%s\n", errStyle.Render(fmt.Sprintf("✘ - %s", m.err)))
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
