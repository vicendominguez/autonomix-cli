package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tim/autonomix-cli/config"
	"github.com/tim/autonomix-cli/pkg/manager"
	"github.com/tim/autonomix-cli/tui"
)

const SelfRepoURL = "https://github.com/timappledotcom/autonomix-cli"

func main() {
	// CLI Argument Handling
	if len(os.Args) > 1 {
		arg := os.Args[1]
		// Determine if "add" command or direct URL
		// "autonomix-cli https://..." or "autonomix-cli add https://..."
		urlToAdd := ""
		if arg == "add" && len(os.Args) > 2 {
			urlToAdd = os.Args[2]
		} else if len(os.Args) == 2 && (arg != "-h" && arg != "--help") {
			// Assume it's a URL if it has slashes, simple check
			if len(arg) > 8 { // https://...
				urlToAdd = arg
			}
		}

		if urlToAdd != "" {
			cfg, err := config.Load()
			if err != nil {
				fmt.Printf("Error loading config: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("Adding repository: %s...\n", urlToAdd)
			res, err := manager.AddApp(cfg, urlToAdd)
			if err != nil {
				fmt.Printf("Error adding app: %v\n", err)
				os.Exit(1)
			}
			
			if res.Created {
				fmt.Printf("Successfully added %s (Latest: %s)\n", res.App.Name, res.App.Latest)
			} else {
				fmt.Printf("Repository %s is already tracked.\n", res.App.Name)
			}
			return
		}
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Ensure self is tracked
	tracked := false
	for _, app := range cfg.Apps {
		if app.RepoURL == SelfRepoURL {
			tracked = true
			break
		}
	}
	if !tracked {
		cfg.Apps = append(cfg.Apps, config.App{
			Name:    "Autonomix CLI",
			RepoURL: SelfRepoURL,
			Version: "dev", // Initial version
		})
		config.Save(cfg)
	}

	p := tea.NewProgram(tui.NewModel(cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
