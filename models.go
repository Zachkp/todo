package main

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
)

type SubTodo struct {
	ID        int
	Title     string
	Completed bool
}

type Todo struct {
	ID          int
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt time.Time
	SubTodos    []SubTodo
}

type viewMode int

const (
	tableView viewMode = iota
	detailView
	addView
	editView
)

type filterMode int

const (
	showAll filterMode = iota
	showActive
	showCompleted
)

type model struct {
	table          table.Model
	todos          []Todo
	mode           viewMode
	filter         filterMode
	titleInput     textinput.Model
	descInput      textarea.Model
	editingID      int
	selectedSubIdx int
	width          int
	height         int
}
