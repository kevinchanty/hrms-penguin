package tui

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Model struct {
	Date         []time.Time
	AttnRecord   []string
	table        table.Model
	textInput    textinput.Model
	windowWidth  int
	windowHeight int
}

func InitialModel() Model {
	columns := []table.Column{
		{Title: "Date", Width: 12},
		{Title: "Attn", Width: 13},
	}

	rows := []table.Row{
		{"2025-07-01", "09:37 18:22"},
		{"2025-07-02", "09:37 18:22"},
		{"2025-07-03", "09:37 18:22"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true)
	t.SetStyles(s)

	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return Model{
		Date:       make([]time.Time, 0),
		AttnRecord: make([]string, 0),
		table:      t,
		textInput:  ti,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.SetWindowTitle("HRMS Penguin - Vacation")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		m.windowWidth = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[0]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {

	// Table
	s := baseStyle.Render(m.table.View())

	// Title
	s = lipgloss.JoinVertical(lipgloss.Center, "Attendance Record Management", s, m.textInput.View())

	// center
	return lipgloss.Place(m.windowWidth, m.windowHeight, lipgloss.Center, lipgloss.Center, s)
}

func StartProgramme() {
	m := InitialModel()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
