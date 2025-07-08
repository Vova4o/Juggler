package app

import (
	"fmt"

	"juggler/internal/config"
	"juggler/internal/juggler"
	"juggler/internal/web"
)

// App represents the main application
type App struct {
	juggler   *juggler.Juggler
	webServer *web.Server
	config    *config.Config
}

// NewApp creates a new application
func NewApp(cfg *config.Config) *App {
	// Create juggler without starting it - it will be configured from frontend
	j := juggler.NewJuggler(0, 0) // Initialize with empty configuration
	webServer := web.NewServer(j, cfg.WebPort)

	return &App{
		juggler:   j,
		webServer: webServer,
		config:    cfg,
	}
}

// Run runs the application
func (a *App) Run() error {
	fmt.Printf("🤹 Жонглер готов к работе!\n")
	fmt.Printf("Веб-интерфейс доступен по адресу: http://localhost:%d\n", a.config.WebPort)
	fmt.Printf("Используйте веб-интерфейс для настройки и управления жонглированием.\n\n")

	// Start web server - this will block
	a.webServer.Start()

	return nil
}
