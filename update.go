package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowResize(msg), nil
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	if m.mode == tableView {
		m.table, cmd = m.table.Update(msg)
	}

	return m, cmd
}

func (m model) handleWindowResize(msg tea.WindowSizeMsg) model {
	m.width = msg.Width
	m.height = msg.Height
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

	return m
}

func (m model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case tableView:
		return m.handleTableViewKeys(msg)
	case detailView:
		return m.handleDetailViewKeys(msg)
	case addView, editView:
		return m.handleEditViewKeys(msg)
	}
	return m, nil
}

func (m model) handleTableViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	return m, nil
}

func (m model) handleDetailViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	return m, nil
}

func (m model) handleEditViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

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
