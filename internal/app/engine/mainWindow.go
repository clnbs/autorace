package engine

import (
	"errors"
	"fmt"
	"image/color"
	"time"

	"github.com/clnbs/autorace/internal/app/models"
	"github.com/clnbs/autorace/pkg/logger"

	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	waitDuration           = 5 * time.Second
	targetedFramePerSecond = 60
	possibleColor          = []color.RGBA{
		colornames.Blue,
		colornames.Red,
		colornames.Violet,
		colornames.Orange,
		colornames.Yellow,
		colornames.Pink,
	}
)

// MainGameWindow structure handle main game window and essential components
type MainGameWindow struct {
	WindowConfiguration MainWindowConfig
	CameraWindow        *MainWindowCamera
	mainWindow          *pixelgl.Window
	ImdDrawer           *imdraw.IMDraw
	GameInfo            *GameCommunication
	events              chan models.Event
}

// NewMainGameWindow create a MainGameWindow structure and feed some of the main components
func NewMainGameWindow(config MainWindowConfig, playerName, rabbitAddr string) (*MainGameWindow, error) {
	var err error
	newMainWindow := new(MainGameWindow)
	newMainWindow.CameraWindow = NewMainWindowCamera()
	newMainWindow.WindowConfiguration = config
	newMainWindow.events = make(chan models.Event)
	newMainWindow.GameInfo, err = NewGameCommunication(playerName, rabbitAddr, newMainWindow.events)
	if err != nil {
		return nil, err
	}
	err = newMainWindow.GameInfo.GetNewPlayer()
	if err != nil {
		return nil, err
	}
	newMainWindow.GenerateCarsSprites()
	return newMainWindow, nil
}

// HandleEvent handle events sent on events communication chan
func (mainGameWindow *MainGameWindow) HandleEvent() {
	event := <-mainGameWindow.events
	switch event.(type) {
	case models.AddCar:
		mainGameWindow.GenerateCarsSprites()
	}
}

// StartCommunicationDaemon start listener needed for the Autorace game. It get can only
// be started if GameInfo.Party is already registered. The only use case who can trigger
// this daemon is when the player created a party and not joining it.
func (mainGameWindow *MainGameWindow) StartCommunicationDaemon() error {
	if mainGameWindow.GameInfo.Party == nil {
		return errors.New("party is not set, could not start communication daemon")
	}
	readyToReceive := make(chan bool)
	go func() {
		err := mainGameWindow.GameInfo.HandleSync(mainGameWindow.GameInfo.Party.PartyUUID.String(), readyToReceive)
		if err != nil {
			logger.Error("error while receiving sync messages :", err)
			return
		}
	}()
	if !<-readyToReceive {
		return errors.New("unable to start HandleSync")
	}
	go func() {
		err := mainGameWindow.GameInfo.ReceiveParty(mainGameWindow.GameInfo.Party.PartyUUID.String(), readyToReceive)
		if err != nil {
			readyToReceive <- false
			logger.Error("error while receiving party :", err)
			return
		}
	}()
	if !<-readyToReceive {
		return errors.New("unable to start ReceiveParty")
	}
	go func() {
		err := mainGameWindow.GameInfo.HandleGameState(mainGameWindow.GameInfo.Party.PartyUUID.String(), readyToReceive)
		if err != nil {
			logger.Error("while receiving new game state :", err)
		}
	}()
	go func() {
		for {
			mainGameWindow.HandleEvent()
		}
	}()
	if !<-readyToReceive {
		return errors.New("unable to start HandleGameState")
	}
	return nil
}

//StartCommunicationDaemonWithPartyID start listener needed for the Autorace game. It can only be trigger
// with an already registered party.
func (mainGameWindow *MainGameWindow) StartCommunicationDaemonWithPartyID(partyID string, readyToReceive chan bool) error {
	ready := make(chan bool)
	go func() {
		err := mainGameWindow.GameInfo.HandleSync(partyID, ready)
		if err != nil {
			readyToReceive <- false
			logger.Error("error while receiving sync messages :", err)
			return
		}
	}()
	if !<-ready {
		return errors.New("unable to start HandleSync")
	}
	go func() {
		err := mainGameWindow.GameInfo.ReceiveParty(partyID, ready)
		if err != nil {
			readyToReceive <- false
			logger.Error("error while receiving party :", err)
			return
		}
	}()
	if !<-ready {
		return errors.New("unable to start ReceiveParty")
	}
	go func() {
		err := mainGameWindow.GameInfo.HandleGameState(partyID, ready)
		if err != nil {
			readyToReceive <- false
			logger.Error("error while receiving new state :", err)
			return
		}
	}()
	if !<-ready {
		return errors.New("unable to start HandleGameState")
	}
	go func() {
		ready <- true
		for {
			mainGameWindow.HandleEvent()
		}
	}()
	<-ready
	readyToReceive <- true
	return nil
}

// Run is used to print graphical components and handle game's logic
func (mainGameWindow *MainGameWindow) Run() {
	err := mainGameWindow.GenerateMainWindowFromConfig()
	if err != nil {
		logger.Error("error while creating new window :", err)
		return
	}
	// create racetrack graphical component
	err = mainGameWindow.constructRacetrack()
	if err != nil {
		logger.Error("error while constructing racetrack :", err)
		return
	}


	frames := 0
	gameTickerDuration := time.Duration((1.0/float64(targetedFramePerSecond))*1000000000) * time.Nanosecond
	gameTicker := time.NewTicker(gameTickerDuration)

	lastTimePauseCalled := time.Now().Add(-5 * time.Second)
	input := new(models.PlayerInput)
	input.PlayerUUID = mainGameWindow.GameInfo.ActorPlayer.Player.PlayerUUID

	// FPS updater -- TODO print FPS in a canvas
	go func() {
		second := time.NewTicker(time.Second)
		for {
			<-second.C
			mainGameWindow.mainWindow.SetTitle(fmt.Sprintf("%s | FPS: %d", mainGameWindow.WindowConfiguration.Title, frames))
			frames = 0
		}
	}()

	for !mainGameWindow.mainWindow.Closed() {
		<-gameTicker.C
		mainGameWindow.SetCamera()

		input.Timestamp = time.Now()
		input.Turning = 0.0
		input.Acceleration = 0.0
		input.MessageNumber = 0

		lastTimePauseCalled = mainGameWindow.Controller(input, lastTimePauseCalled)
		mainGameWindow.CenterCameraOnPlayer()
		mainGameWindow.clear()
		mainGameWindow.PrintGraphicComponents()

		if mainGameWindow.GameInfo.Party.GetState() == models.END {
			logger.Debug("End of the game, bye bye !")
			return
		}
		frames++
	}
}
