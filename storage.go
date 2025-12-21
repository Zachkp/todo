package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func getTodoFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "todos.csv"
	}
	return filepath.Join(home, "Documents", "todos.csv")
}

var csvFile = getTodoFilePath()

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
