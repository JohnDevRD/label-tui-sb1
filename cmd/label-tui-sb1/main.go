package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/JohnDevRD/label-tui-sb1/internal/core"
	"github.com/JohnDevRD/label-tui-sb1/internal/tui"
)

func main() {
	templates, err := core.ListTemplates()
	if err != nil {
		log.Fatalf("loading templates: %v", err)
	}

	m := tui.New()
	m.SetTemplates(templates)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
