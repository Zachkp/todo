package main

import (
	"fmt"
	"strconv"
	"strings"
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

	desc, subTodos := parseSubTodosFromDescription(description)

	newTodo := Todo{
		ID:          maxID + 1,
		Title:       title,
		Description: desc,
		Completed:   false,
		CreatedAt:   time.Now(),
		SubTodos:    subTodos,
	}

	m.todos = append(m.todos, newTodo)
	saveTodos(m.todos)
	m.updateTable()
}

func (m *model) updateTodo(id int, title, description string) {
	desc, subTodos := parseSubTodosFromDescription(description)

	for i, todo := range m.todos {
		if todo.ID == id {
			m.todos[i].Title = title
			m.todos[i].Description = desc
			m.todos[i].SubTodos = subTodos
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

		desc := todo.Description
		if len(todo.SubTodos) > 0 {
			completedCount := 0
			for _, sub := range todo.SubTodos {
				if sub.Completed {
					completedCount++
				}
			}
			desc += fmt.Sprintf(" (%d/%d done)", completedCount, len(todo.SubTodos))
		}

		rows = append(rows, table.Row{
			strconv.Itoa(todo.ID),
			todo.Title,
			desc,
			completed,
		})
	}
	m.table.SetRows(rows)
}

func (m *model) toggleSubTodo(idx int) {
	if idx < 0 {
		return
	}

	var todoIdx int
	var found bool

	if m.mode == detailView {
		currentTodo := m.getCurrentTodo()
		if currentTodo == nil {
			return
		}
		for i, todo := range m.todos {
			if todo.ID == currentTodo.ID {
				todoIdx = i
				found = true
				break
			}
		}
	} else if m.mode == editView && m.editingID != 0 {
		for i, todo := range m.todos {
			if todo.ID == m.editingID {
				todoIdx = i
				found = true
				break
			}
		}
	}

	if !found || idx >= len(m.todos[todoIdx].SubTodos) {
		return
	}

	m.todos[todoIdx].SubTodos[idx].Completed = !m.todos[todoIdx].SubTodos[idx].Completed

	if m.mode == detailView {
		saveTodos(m.todos)
		m.updateTable()
	}
}

func parseSubTodosFromDescription(description string) (string, []SubTodo) {
	lines := strings.Split(description, "\n")
	var descLines []string
	var subTodos []SubTodo
	subID := 1

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "- ") {
			subText := strings.TrimPrefix(trimmed, "- ")
			subText = strings.TrimSpace(subText)

			if subText != "" {
				subTodos = append(subTodos, SubTodo{
					ID:        subID,
					Title:     subText,
					Completed: false,
				})
				subID++
			}
		} else if trimmed != "" {
			descLines = append(descLines, line)
		}
	}

	cleanDesc := strings.Join(descLines, "\n")
	cleanDesc = strings.TrimSpace(cleanDesc)

	return cleanDesc, subTodos
}

