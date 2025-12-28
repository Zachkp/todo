package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--quick" || os.Args[1] == "-q") {
		force := len(os.Args) > 2 && (os.Args[2] == "--force" || os.Args[2] == "-f")
		showQuickView(force)
		return
	}

	todos, err := loadTodos()
	if err != nil {
		fmt.Println("Error loading todos:", err)
		os.Exit(1)
	}

	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Title", Width: 30},
		{Title: "Description", Width: 60},
		{Title: "Done", Width: 6},
	}

	rows := []table.Row{}
	for _, todo := range todos {
		completed := "[ ]"
		if todo.Completed {
			completed = "[✓]"
		}
		rows = append(rows, table.Row{
			strconv.Itoa(todo.ID),
			todo.Title,
			todo.Description,
			completed,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	ti := textinput.New()
	ti.Placeholder = "Enter todo title"
	ti.CharLimit = 100
	ti.Width = 50

	ta := textarea.New()
	ta.Placeholder = "Enter todo description"
	ta.SetWidth(50)
	ta.SetHeight(5)

	m := model{
		table:          t,
		todos:          todos,
		mode:           tableView,
		titleInput:     ti,
		descInput:      ta,
		selectedSubIdx: 0,
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
