package main

import (
	"fmt"
	"hrms-penguin/internal/hrms"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

type formModal str

func initialModel() model {
	return model{
		choices:  []string{"Buy a", "Go b", "Do C"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
		}
	}
	return m, nil
}

func (m model) View() string {
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

	return s
}

func testClient() {
	// hrmsClient := hrms.Client{Host: ""}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	hrmsHost := os.Getenv("HRMS_HOST")
	hrmsUser := os.Getenv("HRMS_USER")
	hrmsPwd := os.Getenv("HRMS_USER")

	client := hrms.New(hrms.ClientOption{
		Host:     hrmsHost,
		UserName: hrmsUser,
		Pwd:      hrmsPwd,
	})

	client.Login()
	client.GetAction()

	// p := tea.NewProgram(initialModel())
	// if _, err := p.Run(); err != nil {
	// 	fmt.Printf("Error occurred: %v", err)
	// 	os.Exit(1)
	// }

}

func main() {

}
