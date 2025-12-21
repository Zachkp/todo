package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func getTodoFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "todos.csv"
	}
	return filepath.Join(home, "Documents", "todos.csv")
}

var csvFile = getTodoFilePath()

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1).
			Margin(1, 2)

	popupStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("62")).
			Bold(true)
)

type Todo struct {
	ID          int
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type viewMode int

const (
	tableView viewMode = iota
	detailView
	addView
	editView
)

type model struct {
	table      table.Model
	todos      []Todo
	mode       viewMode
	titleInput textinput.Model
	descInput  textarea.Model
	editingID  int
	width      int
	height     int
}

func loadTodos() ([]Todo, error) {
	file, err := os.Open(csvFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Todo{}, nil
		}
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var todos []Todo
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 4 {
			continue
		}

		id, _ := strconv.Atoi(record[0])
		completed := record[3] == "true"

		var createdAt, completedAt time.Time
		if len(record) > 4 && record[4] != "" {
			createdAt, _ = time.Parse(time.RFC3339, record[4])
		}
		if len(record) > 5 && record[5] != "" {
			completedAt, _ = time.Parse(time.RFC3339, record[5])
		}

		todos = append(todos, Todo{
			ID:          id,
			Title:       record[1],
			Description: record[2],
			Completed:   completed,
			CreatedAt:   createdAt,
			CompletedAt: completedAt,
		})
	}

	return todos, nil
}

func saveTodos(todos []Todo) error {
	file, err := os.Create(csvFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"ID", "Title", "Description", "Completed", "CreatedAt", "CompletedAt"})

	for _, todo := range todos {
		createdAt := ""
		if !todo.CreatedAt.IsZero() {
			createdAt = todo.CreatedAt.Format(time.RFC3339)
		}
		completedAt := ""
		if !todo.CompletedAt.IsZero() {
			completedAt = todo.CompletedAt.Format(time.RFC3339)
		}

		record := []string{
			strconv.Itoa(todo.ID),
			todo.Title,
			todo.Description,
			strconv.FormatBool(todo.Completed),
			createdAt,
			completedAt,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func (m *model) addTodo(title, description string) {
	maxID := 0
	for _, todo := range m.todos {
		if todo.ID > maxID {
			maxID = todo.ID
		}
	}

	newTodo := Todo{
		ID:          maxID + 1,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}

	m.todos = append(m.todos, newTodo)
	saveTodos(m.todos)
	m.updateTable()
}

func (m *model) updateTodo(id int, title, description string) {
	for i, todo := range m.todos {
		if todo.ID == id {
			m.todos[i].Title = title
			m.todos[i].Description = description
			break
		}
	}
	saveTodos(m.todos)
	m.updateTable()
}

func (m *model) deleteTodo(id int) {
	for i, todo := range m.todos {
		if todo.ID == id {
			m.todos = append(m.todos[:i], m.todos[i+1:]...)
			break
		}
	}

	for i := range m.todos {
		m.todos[i].ID = i + 1
	}

	saveTodos(m.todos)
	m.updateTable()
}

func (m *model) toggleComplete(id int) {
	for i, todo := range m.todos {
		if todo.ID == id {
			m.todos[i].Completed = !m.todos[i].Completed
			if m.todos[i].Completed {
				m.todos[i].CompletedAt = time.Now()
			} else {
				m.todos[i].CompletedAt = time.Time{}
			}
			break
		}
	}
	saveTodos(m.todos)
	m.updateTable()
}

func (m *model) getCurrentTodo() *Todo {
	if len(m.table.SelectedRow()) == 0 {
		return nil
	}
	id, _ := strconv.Atoi(m.table.SelectedRow()[0])
	for _, todo := range m.todos {
		if todo.ID == id {
			return &todo
		}
	}
	return nil
}

func (m *model) updateTable() {
	rows := []table.Row{}
	for _, todo := range m.todos {
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
	m.table.SetRows(rows)
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Account for top margin (1), bottom margin (1), padding (2), help text (2), borders (2)
		m.table.SetHeight(msg.Height - 10)

		availableWidth := m.width - 16
		if availableWidth < 40 {
			availableWidth = 40
		}

		remainingWidth := availableWidth - 10
		titleWidth := remainingWidth / 3
		if titleWidth < 15 {
			titleWidth = 15
		}
		descWidth := remainingWidth - titleWidth
		if descWidth < 20 {
			descWidth = 20
		}

		m.table.SetColumns([]table.Column{
			{Title: "ID", Width: 4},
			{Title: "Title", Width: titleWidth},
			{Title: "Description", Width: descWidth},
			{Title: "Done", Width: 6},
		})

	case tea.KeyMsg:
		switch m.mode {
		case tableView:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "a":
				m.mode = addView
				m.titleInput.Reset()
				m.descInput.Reset()
				m.titleInput.Focus()
				return m, nil
			case "e":
				todo := m.getCurrentTodo()
				if todo != nil {
					m.mode = editView
					m.editingID = todo.ID
					m.titleInput.SetValue(todo.Title)
					m.descInput.SetValue(todo.Description)
					m.titleInput.Focus()
				}
				return m, nil
			case "d":
				todo := m.getCurrentTodo()
				if todo != nil {
					m.deleteTodo(todo.ID)
				}
				return m, nil
			case "enter":
				m.mode = detailView
				return m, nil
			case " ":
				todo := m.getCurrentTodo()
				if todo != nil {
					m.toggleComplete(todo.ID)
				}
				return m, nil
			}

		case detailView:
			switch msg.String() {
			case "esc", "q", "enter":
				m.mode = tableView
				return m, nil
			case "d":
				todo := m.getCurrentTodo()
				if todo != nil {
					m.deleteTodo(todo.ID)
					m.mode = tableView
				}
				return m, nil
			case " ":
				todo := m.getCurrentTodo()
				if todo != nil {
					m.toggleComplete(todo.ID)
				}
				return m, nil
			}

		case addView, editView:
			switch msg.String() {
			case "esc":
				m.mode = tableView
				m.titleInput.Blur()
				m.descInput.Blur()
				return m, nil
			case "ctrl+s":
				title := strings.TrimSpace(m.titleInput.Value())
				desc := strings.TrimSpace(m.descInput.Value())
				if title != "" {
					if m.mode == addView {
						m.addTodo(title, desc)
					} else {
						m.updateTodo(m.editingID, title, desc)
					}
					m.mode = tableView
					m.titleInput.Blur()
					m.descInput.Blur()
				}
				return m, nil
			case "tab":
				if m.titleInput.Focused() {
					m.titleInput.Blur()
					m.descInput.Focus()
				} else {
					m.descInput.Blur()
					m.titleInput.Focus()
				}
				return m, nil
			}

			if m.titleInput.Focused() {
				m.titleInput, cmd = m.titleInput.Update(msg)
			} else {
				m.descInput, cmd = m.descInput.Update(msg)
			}
			return m, cmd
		}
	}

	if m.mode == tableView {
		m.table, cmd = m.table.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	switch m.mode {
	case detailView:
		todo := m.getCurrentTodo()
		if todo == nil {
			return "No todo selected"
		}

		status := "Incomplete"
		if todo.Completed {
			status = "Completed"
		}

		popupWidth := m.width - 20
		if popupWidth > 80 {
			popupWidth = 80
		}
		if popupWidth < 40 {
			popupWidth = 40
		}

		content := fmt.Sprintf("%s\n\n", titleStyle.Render(todo.Title))
		content += fmt.Sprintf("Status: %s\n\n", status)
		content += fmt.Sprintf("Description:\n%s\n\n", todo.Description)

		if !todo.CreatedAt.IsZero() {
			content += fmt.Sprintf("Created: %s\n", todo.CreatedAt.Format("Jan 2, 2006 at 3:04 PM"))
		}
		if !todo.CompletedAt.IsZero() {
			content += fmt.Sprintf("Completed: %s\n", todo.CompletedAt.Format("Jan 2, 2006 at 3:04 PM"))
		}
		content += "\n"

		content += lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(
			"[enter/esc] back  [space] toggle  [d] delete",
		)

		popup := popupStyle.Width(popupWidth).Render(content)
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			popup,
		)

	case addView, editView:
		title := "Add Todo"
		if m.mode == editView {
			title = "Edit Todo"
		}

		popupWidth := m.width - 20
		if popupWidth > 80 {
			popupWidth = 80
		}
		if popupWidth < 40 {
			popupWidth = 40
		}

		inputWidth := popupWidth - 10
		if inputWidth < 30 {
			inputWidth = 30
		}
		m.titleInput.Width = inputWidth
		m.descInput.SetWidth(inputWidth)

		content := titleStyle.Render(title) + "\n\n"
		content += "Title:\n" + m.titleInput.View() + "\n\n"
		content += "Description:\n" + m.descInput.View() + "\n\n"
		content += lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(
			"[ctrl+s] save  [tab] switch field  [esc] cancel",
		)

		popup := popupStyle.Width(popupWidth).Render(content)
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			popup,
		)

	default:
		help := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Width(m.width).
			Align(lipgloss.Center).
			Render("[a] add  [e] edit  [d] delete  [space] toggle  [enter] details  [q] quit")
		return baseStyle.Render(m.table.View()) + "\n" + help + "\n"
	}
}

func main() {
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
		table:      t,
		todos:      todos,
		mode:       tableView,
		titleInput: ti,
		descInput:  ta,
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
