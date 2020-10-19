package engine

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// MainWindowConfig stores Autorace's game and graphics configuration.
// It is use to set-up game's window
type MainWindowConfig struct {
	Title                string
	TargetFramePerSecond float32
	MinX                 float64
	MaxX                 float64
	MinY                 float64
	MaxY                 float64
	Vsync                bool
	Smooth               bool
}

//NewWindowConfiguration create an instance of MainWindowConfig
// TODO get configuration from files or server
func NewWindowConfiguration() MainWindowConfig {
	var newMainWindowConfig MainWindowConfig
	newMainWindowConfig.Title = "Autorace"
	newMainWindowConfig.Smooth = true
	newMainWindowConfig.Vsync = false
	newMainWindowConfig.TargetFramePerSecond = 144.0
	newMainWindowConfig.MinX = 0
	newMainWindowConfig.MinY = 0
	newMainWindowConfig.MaxX = 1024
	newMainWindowConfig.MaxY = 768
	return newMainWindowConfig
}

//GenerateMainWindowFromConfig generate main game window from a MainWindowConfig instance.
// GenerateMainWindowFromConfig extends MainGameWindow methods
func (mainGameWindow *MainGameWindow) GenerateMainWindowFromConfig() error {
	var err error
	// create configuration that fit Pixel's requirement from a MainWindowConfig struct
	pixelWindowConfig := pixelgl.WindowConfig{
		Title:       mainGameWindow.WindowConfiguration.Title,
		Icon:        nil,
		Bounds:      pixel.R(0, 0, 1024, 768),
		Monitor:     nil,
		Resizable:   true,
		Undecorated: false,
		NoIconify:   false,
		AlwaysOnTop: false,
		VSync:       mainGameWindow.WindowConfiguration.Vsync,
	}
	mainGameWindow.mainWindow, err = pixelgl.NewWindow(pixelWindowConfig)
	if err != nil {
		return err
	}
	mainGameWindow.mainWindow.SetSmooth(mainGameWindow.WindowConfiguration.Smooth)
	mainGameWindow.mainWindow.Clear(colornames.Forestgreen)
	return nil
}
