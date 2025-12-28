package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

func getSessionLockFile() string {
	tmpDir := os.TempDir()
	uid := os.Getuid()
	return filepath.Join(tmpDir, fmt.Sprintf("todo-shown-%d", uid))
}

func hasShownThisSession() bool {
	lockFile := getSessionLockFile()
	_, err := os.Stat(lockFile)
	return err == nil
}

func markShownThisSession() error {
	lockFile := getSessionLockFile()
	file, err := os.Create(lockFile)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func isInteractiveShell() bool {
	shlvl := os.Getenv("SHLVL")
	if shlvl == "" || shlvl == "0" || shlvl == "1" {
		return false
	}

	shell := os.Getenv("SHELL")
	if shell == "" {
		return false
	}

	term := os.Getenv("TERM")
	if term == "" || term == "dumb" {
		return false
	}

	return true
}

func showQuickView(force bool) {
	if !isInteractiveShell() {
		return
	}

	if !force && hasShownThisSession() {
		return
	}

	todos, err := loadTodos()
	if err != nil {
		fmt.Println("Error loading todos:", err)
		os.Exit(1)
	}

	activeTodos := []Todo{}
	for _, todo := range todos {
		if !todo.Completed {
			activeTodos = append(activeTodos, todo)
		}
	}
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)

	todoTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	subTodoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("248")).
		MarginLeft(2)

	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214"))

	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true)

	fmt.Println(titleStyle.Render("📋 Active Todos"))

	if len(activeTodos) == 0 {
		fmt.Println(emptyStyle.Render("No active todos! 🎉"))
		fmt.Println()
		return
	}

	for i, todo := range activeTodos {
		fmt.Printf("%s %s\n",
			todoTitleStyle.Render(fmt.Sprintf("%d.", todo.ID)),
			todoTitleStyle.Render(todo.Title))

		if todo.Description != "" {
			fmt.Printf("   %s\n", descStyle.Render(todo.Description))
		}

		if len(todo.SubTodos) > 0 {
			completed := 0
			for _, sub := range todo.SubTodos {
				checkbox := "[ ]"
				if sub.Completed {
					checkbox = "[✓]"
					completed++
				}
				fmt.Printf("   %s %s\n",
					subTodoStyle.Render(checkbox),
					subTodoStyle.Render(sub.Title))
			}

			progressText := fmt.Sprintf("(%d/%d done)", completed, len(todo.SubTodos))
			fmt.Printf("   %s\n", progressStyle.Render(progressText))
		}

		if i < len(activeTodos)-1 {
			fmt.Println()
		}
	}

	fmt.Println()
	fmt.Println(descStyle.Render("💡 Run 'todo' to open the full app"))
	fmt.Println()

	markShownThisSession()
}
