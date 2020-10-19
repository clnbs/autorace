package engine

import (
	"time"

	"github.com/clnbs/autorace/internal/app/models"
	"github.com/clnbs/autorace/pkg/logger"

	"github.com/faiface/pixel/pixelgl"
)

//Controller extends MainGameWindow methods and handles all keyboard event.
// Controller also send player's input to a dynamic instance.
func (mainGameWindow *MainGameWindow) Controller(input *models.PlayerInput, lastTimePauseCalled time.Time) time.Time {
	var err error
	if mainGameWindow.mainWindow.Pressed(pixelgl.KeyLeft) {
		input.Turning += 1.0
	}
	if mainGameWindow.mainWindow.Pressed(pixelgl.KeyRight) {
		input.Turning -= 1.0
	}
	if mainGameWindow.mainWindow.Pressed(pixelgl.KeyUp) {
		input.Acceleration += 1.0
	}
	if mainGameWindow.mainWindow.Pressed(pixelgl.KeyDown) {
		input.Acceleration -= 1.0
	}
	if mainGameWindow.mainWindow.Pressed(pixelgl.KeyP) && time.Since(lastTimePauseCalled) > waitDuration {
		lastTimePauseCalled = time.Now()
		var desiredState models.State
		if mainGameWindow.GameInfo.Party.GetState() == models.RUN {
			desiredState = models.PAUSE
		} else {
			desiredState = models.RUN
		}
		go func() {
			err = mainGameWindow.GameInfo.SendState(desiredState)
			if err != nil {
				logger.Error("error while sending new state :", err)
			}
		}()
	}

	if mainGameWindow.mainWindow.Pressed(pixelgl.KeyEnd) {
		go func() {
			err = mainGameWindow.GameInfo.SendState(models.END)
			if err != nil {
				logger.Error("error while sending new state :", err)
			}
		}()
	}

	if mainGameWindow.GameInfo.Party.GetState() == models.RUN {
		mainGameWindow.GameInfo.ActorPlayer.Player.Input = input
		mainGameWindow.GameInfo.SendPlayerInput()
	}
	return lastTimePauseCalled
}
