package tui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	colorPrimary   = "#10B981"
	colorSecondary = "#34D399"
	colorAccent    = "#FBBF24"
	colorSuccess   = "#10B981"
	colorError     = "#EF4444"
	colorSurface   = "#0F172A"
	colorBase      = "#020617"
	colorText      = "#E2E8F0"
	colorMuted     = "#64748B"
	colorBorder    = "#1E293B"
)

var (
	AppStyle = lipgloss.NewStyle().
		Padding(1, 2).
		Background(lipgloss.Color(colorBase))

	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorPrimary)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colorPrimary)).
		Padding(0, 2).
		MarginBottom(1)

	SectionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorSecondary)).
		Bold(true).
		MarginTop(1).
		MarginBottom(1)

	MenuItemStyle = lipgloss.NewStyle().
		Padding(0, 2).
		Margin(0, 1)

	MenuSelectedStyle = lipgloss.NewStyle().
		Padding(0, 2).
		Margin(0, 1).
		Background(lipgloss.Color(colorPrimary)).
		Foreground(lipgloss.Color(colorText)).
		Bold(true)

	InputStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colorBorder)).
		Padding(0, 1).
		Width(50)

	InputFocusedStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colorPrimary)).
		Padding(0, 1).
		Width(50)

	ArticleItemStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Margin(0, 1).
		Width(70)

	ArticleSelectedStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Margin(0, 1).
		Width(70).
		Background(lipgloss.Color(colorPrimary)).
		Foreground(lipgloss.Color(colorText))

	ArticleCheckedStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Margin(0, 1).
		Width(70).
		Foreground(lipgloss.Color(colorSuccess))

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorError)).
		MarginTop(1)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorSuccess)).
		MarginTop(1)

	InfoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMuted)).
		MarginTop(1)

	HelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMuted)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colorBorder)).
		Padding(0, 1).
		MarginTop(1)

	StatusBarStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(colorSurface)).
		Foreground(lipgloss.Color(colorMuted)).
		Padding(0, 2)

	PreviewBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colorBorder)).
		Padding(1, 2).
		MaxWidth(80)

	ButtonStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(colorPrimary)).
		Foreground(lipgloss.Color(colorText)).
		Bold(true).
		Padding(0, 3).
		MarginTop(1)

	DimmedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMuted))

	SpinnerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorAccent)).
		Bold(true)
)

func RenderHeader(title string) string {
	return TitleStyle.Render(title)
}

func RenderMenu(items []string, selected int) string {
	var s string
	for i, item := range items {
		if i == selected {
			s += MenuSelectedStyle.Render("▸ " + item) + "\n"
		} else {
			s += MenuItemStyle.Render("  "+item) + "\n"
		}
	}
	return s
}

func RenderHelp(keys ...string) string {
	var s string
	for i, k := range keys {
		if i > 0 {
			s += "  •  "
		}
		s += k
	}
	return HelpStyle.Render(s)
}

func RenderStatus(text string) string {
	return StatusBarStyle.Render(text)
}
