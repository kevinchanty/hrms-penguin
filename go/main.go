package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	action Action

	answerField    textinput.Model
	textInputValue string
	styles         *Styles
}

func New() *Model {
	answerField := textinput.New()
	answerField.Placeholder = "I am your placeholder"
	answerField.Focus()

	return &Model{
		answerField: answerField,
		styles:      DefaultStyles(),
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
		// switch msg.String() {
		// case "ctrl+c", "q":
		// 	return m, tea.Quit

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
		// }
	}

	m.answerField, cmd = m.answerField.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	s := "Your Main Action?\n\n"

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			s,
			m.styles.InputField.Render(
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

	want := &Action{
		missAttendance: make([]string, 0, 31),
		earlyLeave:     make([]string, 0, 31),
		lateness:       make([]string, 0, 31),
	}
	want.missAttendance = append(want.missAttendance, "2023-12-21", "2023-12-28", "2024-01-08", "2024-01-11", "2024-01-12", "2024-01-15")
	want.earlyLeave = append(want.earlyLeave, "2023-12-18", "2024-01-17")
	want.lateness = append(want.lateness, "2023-12-18", "2024-01-04")

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
