package main

import "github.com/charmbracelet/lipgloss"

var (
	// baseStyle is the main container style
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1).
			Margin(1, 2)

	// popupStyle is used for detail, add, and edit views
	popupStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)

	// titleStyle is used for popup titles
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("62")).
			Bold(true)

	// completedStyle is used for completed checkmarks (green)
	completedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	// incompleteStyle is used for incomplete checkmarks (gray)
	incompleteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)
