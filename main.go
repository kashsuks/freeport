package main

import (
	"fmt"
	"os"
	"freeport/ui"
	"freeport/api"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	server := api.NewServer("6767")
	go func() {
		if err := server.Start(); err != nil {
			fmt.Printf("API Server Error: %v\n", err)
		}
	}()

	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}