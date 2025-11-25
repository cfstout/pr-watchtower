package main

import (
	"fmt"
	"os"

	"github.com/cfstout/pr-watchtower/internal/config"
	"github.com/cfstout/pr-watchtower/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(tui.InitialModel(cfg))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
