package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

// getSessionLockFile returns the path to the session lock file
func getSessionLockFile() string {
	// Use /tmp which gets cleared on reboot
	tmpDir := os.TempDir()
	// Use UID to make it user-specific
	uid := os.Getuid()
	return filepath.Join(tmpDir, fmt.Sprintf("todo-shown-%d", uid))
}

// hasShownThisSession checks if quick view was already shown this login session
func hasShownThisSession() bool {
	lockFile := getSessionLockFile()
	_, err := os.Stat(lockFile)
	return err == nil // File exists = already shown
}

// markShownThisSession creates a lock file to indicate quick view was shown
func markShownThisSession() error {
	lockFile := getSessionLockFile()
	file, err := os.Create(lockFile)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func showQuickView(force bool) {
	// Check if already shown this session (unless forced)
	if !force && hasShownThisSession() {
		// Silently exit - already shown this session
		return
	}
	todos, err := loadTodos()
	if err != nil {
		fmt.Println("Error loading todos:", err)
		os.Exit(1)
	}

	// Filter to active todos only
	activeTodos := []Todo{}
	for _, todo := range todos {
		if !todo.Completed {
			activeTodos = append(activeTodos, todo)
		}
	}

	// Styles
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

	// Print header
	fmt.Println(titleStyle.Render("📋 Active Todos"))

	if len(activeTodos) == 0 {
		fmt.Println(emptyStyle.Render("No active todos! 🎉"))
		fmt.Println()
		return
	}

	// Print each active todo
	for i, todo := range activeTodos {
		// Todo title with ID
		fmt.Printf("%s %s\n",
			todoTitleStyle.Render(fmt.Sprintf("%d.", todo.ID)),
			todoTitleStyle.Render(todo.Title))

		// Description (if present)
		if todo.Description != "" {
			fmt.Printf("   %s\n", descStyle.Render(todo.Description))
		}

		// Sub-todos with progress
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

			// Progress indicator
			progressText := fmt.Sprintf("(%d/%d done)", completed, len(todo.SubTodos))
			fmt.Printf("   %s\n", progressStyle.Render(progressText))
		}

		// Add spacing between todos (but not after the last one)
		if i < len(activeTodos)-1 {
			fmt.Println()
		}
	}

	fmt.Println()
	fmt.Println(descStyle.Render("💡 Run 'todo' to open the full app"))
	fmt.Println()

	// Mark as shown for this session
	markShownThisSession()
}
