package main

import (
	"fmt"

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
	choices        []string
	cursor         int
	selected       map[int]struct{}
	answerField    textinput.Model
	textInputValue string
	styles         *Styles
	width          int // for storing windows width at init
	height         int // for storing windows height at init
}

func New() *Model {
	answerField := textinput.New()
	answerField.Placeholder = "I am your placeholder"
	answerField.Focus()

	return &Model{
		choices:     []string{"Buy a", "Go b", "Do C"},
		selected:    make(map[int]struct{}),
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
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
			m.textInputValue = m.answerField.Value()
			m.answerField.SetValue("")
		}
	}

	m.answerField, cmd = m.answerField.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	s := "What should we do in weekend?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"

	s += "\n" + m.textInputValue

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

	// tea
	// p := tea.NewProgram(New())
	// if _, err := p.Run(); err != nil {
	// 	log.Fatal(err)
	// }

	// actual client
	client := NewHrmsClient(ClientOption{
		Host:     "https://hrms.hktv.com.hk",
		UserName: "tychan",
		Pwd:      "196HRMS=",
	})
	client.Login()
	client.GetAction()

	// client parsing
	// client.ParseMainAction(haha)
}

var haha string = "<p>Missing Attendance record 欠缺出入勤紀錄:<br /> 2023-12-21<br />2023-12-28<br />2024-01-08<br />2024-01-11<br />2024-01-12<br />2024-01-15</p><p>Early leave:<br /> 2023-12-18<br />2024-01-17</p><p>Lateness 遲到:<br /> 2023-12-18<br />2024-01-04</p>"
