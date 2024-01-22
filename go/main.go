package main

import (
	"slices"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var tableXOffset int = 20
var tableYOffset int = 6

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Model struct {
	width  int // for storing windows width at init
	height int // for storing windows height at init

	actionTable   table.Model
	actionRows    []table.Row
	selectedDate  []string
	isActionReady bool

	mainStyle lipgloss.Style
}

func New() *Model {
	mainStyle := lipgloss.NewStyle().Padding(1)

	actionRows := []table.Row{
		{"2023-12-21", "Missing Attendance record 欠缺出入勤紀錄"}, {"2023-12-28", "Missing Attendance record 欠缺出入勤紀錄"},
		// {"2024-01-08", "Missing Attendance record 欠缺出入勤紀錄"}, {"2024-01-11", "Missing Attendance record 欠缺出入勤紀錄"},
		// {"2024-01-12", "Missing Attendance record 欠缺出入勤紀錄"}, {"2024-01-15", "Missing Attendance record 欠缺出入勤紀錄"},
		// {"2023-12-18", "Early leave"},
		// {"2024-01-17", "Early leave"},
		// {"2023-12-18", "Lateness 遲到"},
		// {"2024-01-04", "Lateness 遲到"},
	}

	actionsColumn := []table.Column{
		{Title: "Date", Width: 10},
		{Title: "type", Width: 60},
	}

	actionTable := table.New(
		table.WithColumns(actionsColumn),
		table.WithRows(actionRows),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Background(lipgloss.Color("212"))
	// Background(lipgloss.Color("0"))
	actionTable.SetStyles(s)

	return &Model{
		mainStyle: mainStyle,

		// answerField:    answerField,
		isActionReady: false,
		actionRows:    actionRows,
		actionTable:   actionTable,
		selectedDate:  []string{},
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
		m.actionTable.SetColumns(
			[]table.Column{
				{Title: "Date", Width: 10},
				{Title: "Type", Width: msg.Width - tableXOffset},
			})
		m.actionTable.SetHeight(msg.Height - tableYOffset)
	case tea.KeyMsg:
		switch msg.String() {
		// case "ctrl+c", "q":
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			if i := slices.Index(m.selectedDate, m.actionTable.SelectedRow()[0]); i > -1 {
				m.selectedDate = slices.Delete(m.selectedDate, i, i+1)
			} else {
				m.selectedDate = append(m.selectedDate, m.actionTable.SelectedRow()[0])
			}
			return m, cmd

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

	m.actionTable, cmd = m.actionTable.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	title := " 🐧 HRMS Penguin - Amend Attendance Record"

	selectedDateStr := "Selected:"

	for i, date := range m.selectedDate {
		if i == 0 {
			selectedDateStr = selectedDateStr + " " + date
		} else {
			selectedDateStr = selectedDateStr + ", " + date
		}
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		baseStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				m.actionTable.View(),
				selectedDateStr,
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

	modal := New()
	// tea
	p := tea.NewProgram(modal)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	// actual client
	// client := NewHrmsClient(ClientOption{
	// 	// Host:     "https://hrms.hktv.com.hk",
	// 	Host:     "http://localhost:8080",
	// 	UserName: "tychan",
	// 	Pwd:      "196HRMS=",
	// })
	// client.Login()
	// actions := client.GetAction()

	modal.actionRows = append(modal.actionRows, []table.Row{{"2023-12-21", "Missing Attendance record 欠缺出入勤紀錄"}, {"2023-12-21", "Missing Attendance record 欠缺出入勤紀錄"}, {"2023-12-21", "Missing Attendance record 欠缺出入勤紀錄"}}...)
	p.Run()

	// client parsing
	// client.ParseMainAction(haha)
}
