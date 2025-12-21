package main

import (
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
)

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
		// Apply filter
		if m.filter == showActive && todo.Completed {
			continue
		}
		if m.filter == showCompleted && !todo.Completed {
			continue
		}

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
