// cmd/main.go
package main

import (
	"github.com/joho/godotenv"
	"github.com/rivo/tview"
	"log"
	"termichat/config"
	"termichat/internal/ui"
)

func main() {
	// Load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Load config
	err = config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := tview.NewApplication()
	app.EnableMouse(true)

	// Set up the main layout using the ui package
	layout := ui.CreateMainLayout(app)

	// Start the application with the layout
	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}
