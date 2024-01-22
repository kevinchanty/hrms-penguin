package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.InputField = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(s.BorderColor).
		Padding(1).
		Width(90)
	return s
}

type Model struct {
	width  int // for storing windows width at init
	height int // for storing windows height at init

	actions       [][]string
	actionTable   table.Table
	isActionReady bool

	answerField    textinput.Model
	textInputValue string
	textInputStyle *Styles
}

func New() *Model {
	answerField := textinput.New()
	answerField.Placeholder = "I am your placeholder"
	answerField.Focus()

	sampleAction := [][]string{
		{"2023-12-21", "Missing Attendance record 欠缺出入勤紀錄"}, {"2023-12-28", "Missing Attendance record 欠缺出入勤紀錄"},
		{"2024-01-08", "Missing Attendance record 欠缺出入勤紀錄"}, {"2024-01-11", "Missing Attendance record 欠缺出入勤紀錄"},
		{"2024-01-12", "Missing Attendance record 欠缺出入勤紀錄"}, {"2024-01-15", "Missing Attendance record 欠缺出入勤紀錄"}, {"2023-12-18", "Early leave"},
		{"2024-01-17", "Early leave"},
		{"2023-12-18", "Lateness 遲到"},
		{"2024-01-04", "Lateness 遲到"}}

	return &Model{
		answerField:    answerField,
		isActionReady:  false,
		textInputStyle: DefaultStyles(),
		actions:        sampleAction,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		// case "ctrl+c", "q":
		case "ctrl+c":
			return m, tea.Quit

			// case "up", "k":
			// 	if m.cursor > 0 {
			// 		m.cursor--
			// 	}
			// case "down", "j":
			// 	if m.cursor < len(m.choices)-1 {
			// 		m.cursor++
			// 	}
			// case "enter", " ":
			// 	_, ok := m.selected[m.cursor]
			// 	if ok {
			// 		delete(m.selected, m.cursor)
			// 	} else {
			// 		m.selected[m.cursor] = struct{}{}
			// 	}
			// 	m.textInputValue = m.answerField.Value()
			// 	m.answerField.SetValue("")
		}
	}

	m.answerField, cmd = m.answerField.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	actionStr := "Your Main Action:\n\n"

	actionStr += ""
	for _, s := range m.actions.missAttendance {
		actionStr += s
		actionStr += "\n"
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			actionStr,
			m.textInputStyle.InputField.Render(
				m.answerField.View(),
			),
		),
	)

}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal("err: %w", err)
	}
	defer f.Close()

	// tea
	p := tea.NewProgram(New())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	// actual client
	// client := NewHrmsClient(ClientOption{
	// 	Host:     "https://hrms.hktv.com.hk",
	// 	UserName: "tychan",
	// 	Pwd:      "196HRMS=",
	// })
	// client.Login()
	// client.GetAction()

	// client parsing
	// client.ParseMainAction(haha)
}
