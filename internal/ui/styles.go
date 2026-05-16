package ui

import "github.com/charmbracelet/lipgloss"

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12"))

	tabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("4")).
			Padding(0, 1)

	indexActiveStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("4"))

	tabInactiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Padding(0, 1)

	tabDisabledStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("8")).
				Padding(0, 1)

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("11"))

	taskCursorMarkStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("15"))

	taskDoneStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	taskPendingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("7"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	errStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)

	progressDoneStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("6"))

	progressCompleteStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("2"))

	progressEmptyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("8"))

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))
)
