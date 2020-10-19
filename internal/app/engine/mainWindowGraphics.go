package engine

import (
	"errors"
	"math"

	"github.com/clnbs/autorace/internal/app/models"
	"github.com/clnbs/autorace/internal/pkg/mathtool"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

// MainWindowCamera handles camera for MainGameWindow.
type MainWindowCamera struct {
	Zoom      float64
	ZoomSpeed float64
	Position  pixel.Vec
	Camera    pixel.Matrix
}

//NewMainWindowCamera generates a fed MainWindowCamera instance with defaults values
func NewMainWindowCamera() *MainWindowCamera {
	newWindowCamera := new(MainWindowCamera)
	newWindowCamera.Zoom = 1.0
	newWindowCamera.ZoomSpeed = 1.2
	newWindowCamera.Position = pixel.ZV
	return newWindowCamera
}

// PrintGraphicComponents extends MainGameWindow methods. It is used to print
// graphical component of MainGameWindow :
// - main player's car
// - competitors' car
// - checkpoints
func (mainGameWindow *MainGameWindow) PrintGraphicComponents() {
	for _, cp := range mainGameWindow.GameInfo.CheckPoints {
		newMatrix := pixel.IM
		newMatrix = newMatrix.Moved(cp.Position)
		cp.CpSprite.Draw(mainGameWindow.mainWindow, newMatrix)
	}
	for _, competitor := range mainGameWindow.GameInfo.Competitors {
		if competitor.Act.Car.CarSprite == nil {
			continue
		}
		newMatrix := pixel.IM
		newMatrix = newMatrix.Rotated(pixel.ZV, competitor.Act.Car.Angle)
		newMatrix = newMatrix.Moved(competitor.Act.Car.Position)
		competitor.Act.Car.CarSprite.Draw(mainGameWindow.mainWindow, newMatrix)
	}
	{
		newMatrix := pixel.IM
		newMatrix = newMatrix.Rotated(pixel.ZV, mainGameWindow.GameInfo.ActorPlayer.Act.Car.Angle)
		newMatrix = newMatrix.Moved(mainGameWindow.GameInfo.ActorPlayer.Act.Car.Position)
		mainGameWindow.GameInfo.ActorPlayer.Act.Car.CarSprite.Draw(mainGameWindow.mainWindow, newMatrix)
	}
	mainGameWindow.mainWindow.Update()
}

//SetCamera is used to set zoom and camera location at each game tick
func (mainGameWindow *MainGameWindow) SetCamera() {
	mainGameWindow.CameraWindow.Camera = pixel.IM.
		Scaled(mainGameWindow.CameraWindow.Position, mainGameWindow.CameraWindow.Zoom).
		Moved(mainGameWindow.mainWindow.Bounds().Center().Sub(mainGameWindow.CameraWindow.Position))
	mainGameWindow.mainWindow.SetMatrix(mainGameWindow.CameraWindow.Camera)
}

//CenterCameraOnPlayer is used to follow player's movement
func (mainGameWindow *MainGameWindow) CenterCameraOnPlayer() {
	mainGameWindow.CameraWindow.Position = mainGameWindow.GameInfo.ActorPlayer.Act.Car.Position
	mainGameWindow.CameraWindow.Zoom *= math.Pow(
		mainGameWindow.CameraWindow.ZoomSpeed,
		mainGameWindow.mainWindow.MouseScroll().Y)
}

// GenerateCarsSprites create car sprite for every participants in the game in order to be printed later
func (mainGameWindow *MainGameWindow) GenerateCarsSprites() {
	// we create a car sprite if it is not already done for the main actor (actual player)
	if mainGameWindow.GameInfo.ActorPlayer.Act.Car == nil {
		randomColorNumber := mathtool.RandomIntBetween(1, len(possibleColor)) - 1
		mainGameWindow.GameInfo.ActorPlayer.Act.Car = models.NewCar(possibleColor[randomColorNumber])
	}
	// for every competitor, we create a car sprite if it is not already created
	for _, competitor := range mainGameWindow.GameInfo.Competitors {
		if competitor.Act.Car.CarSprite == nil {
			tmpPosition := competitor.Act.Car.Position
			tmpAngle := competitor.Act.Car.Angle
			randomColorNumber := mathtool.RandomIntBetween(1, len(possibleColor)) - 1
			competitor.Act.Car = models.NewCar(possibleColor[randomColorNumber])
			competitor.Act.Car.Position = tmpPosition
			competitor.Act.Car.Angle = tmpAngle
		}
	}
}

// generate racetrack graphically
func (mainGameWindow *MainGameWindow) constructRacetrack() error {
	if mainGameWindow.GameInfo.Party == nil {
		return errors.New("party is not feed")
	}
	mainGameWindow.ImdDrawer = imdraw.New(nil)
	mainGameWindow.ImdDrawer.Color = colornames.Gray
	mainGameWindow.ImdDrawer.EndShape = imdraw.RoundEndShape
	for _, cp := range mainGameWindow.GameInfo.Party.MapCircuit.TurnPoints {
		mainGameWindow.ImdDrawer.Push(pixel.V(cp.Position.X, cp.Position.Y))
	}
	mainGameWindow.ImdDrawer.Line(100)
	return nil
}

// clear MainGameWindow in order to print graphical components
func (mainGameWindow *MainGameWindow) clear() {
	mainGameWindow.mainWindow.Clear(colornames.Forestgreen)
	mainGameWindow.ImdDrawer.Draw(mainGameWindow.mainWindow)
}

