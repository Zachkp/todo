package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	switch m.mode {
	case detailView:
		return m.renderDetailView()
	case addView, editView:
		return m.renderEditView()
	default:
		return m.renderTableView()
	}
}

func (m model) renderTableView() string {
	// Show empty state if no todos
	if len(m.todos) == 0 {
		emptyMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Width(m.width).
			Align(lipgloss.Center).
			Render("No todos yet! Press 'a' to add your first todo.")

		help := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Width(m.width).
			Align(lipgloss.Center).
			Render("[a] add  [q] quit")

		content := "\n\n" + emptyMsg + "\n\n"
		return baseStyle.Render(content) + "\n" + help + "\n"
	}

	// Show filter status
	filterStatus := "All"
	switch m.filter {
	case showActive:
		filterStatus = "Active"
	case showCompleted:
		filterStatus = "Completed"
	}

	filterText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("62")).
		Width(m.width).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("Filter: %s", filterStatus))

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Width(m.width).
		Align(lipgloss.Center).
		Render("[a] add  [e] edit  [d] delete  [space] toggle  [f] filter  [enter] details  [q] quit")

	return baseStyle.Render(m.table.View()) + "\n" + filterText + "\n" + help + "\n"
}

func (m model) renderDetailView() string {
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

	if len(todo.SubTodos) > 0 {
		content += "\nSub-Todos:\n"
		for i, sub := range todo.SubTodos {
			checkbox := "[ ]"
			if sub.Completed {
				checkbox = "[✓]"
			}

			subLine := fmt.Sprintf("  %s %s", checkbox, sub.Title)
			if i == m.selectedSubIdx {
				subLine = lipgloss.NewStyle().
					Background(lipgloss.Color("57")).
					Foreground(lipgloss.Color("229")).
					Render(subLine)
			}
			content += subLine + "\n"
		}
	}

	content += "\n"

	helpText := "[enter/esc] back  [d] delete"
	if len(todo.SubTodos) > 0 {
		helpText = "[enter/esc] back  [space] toggle sub  [↑↓] navigate  [d] delete"
	}
	content += lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(helpText)

	popup := popupStyle.Width(popupWidth).Render(content)
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		popup,
	)
}

func (m model) renderEditView() string {
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
	content += "Description:\n" + m.descInput.View() + "\n"
	content += lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(
		"(Use '- ' at start of line for sub-todos)",
	) + "\n\n"

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
}
