package models

import (
	"fmt"
	"time"

	"github.com/clnbs/autorace/internal/pkg/mathtool"

	"github.com/google/uuid"
)

// Player represent a game's participant
type Player struct {
	PlayerName string          `json:"player_name"`
	PlayerUUID uuid.UUID       `json:"player_uuid"`
	Position   *PlayerPosition `json:"position"`
	Input      *PlayerInput    `json:"input,omitempty"`
}

// PlayerPosition represent a player's position in the race
type PlayerPosition struct {
	CurrentSpeed    float64          `json:"current_speed"`
	CurrentAngle    float64          `json:"current_angle"`
	CurrentPosition mathtool.Vector2 `json:"current_position"`
}

// PlayerCreationToken is issued to server when a client want to be registered server side
type PlayerCreationToken struct {
	SessionUUID uuid.UUID `json:"session_uuid"`
	PlayerName  string    `json:"player_name"`
}

// PlayerInput hold player input and is send to the server side in order to be processed
type PlayerInput struct {
	Acceleration  float64   `json:"acceleration"`
	Turning       float64   `json:"turning"`
	MessageNumber int       `json:"message_number,omitempty"`
	Timestamp     time.Time `json:"timestamp,omitempty"`
	PlayerUUID    uuid.UUID `json:"player_uuid"`
}

// PlayerToken is issued when a player ask for a particular action server side
type PlayerToken struct {
	ClientID string `json:"client_id"`
	PartyID  string `json:"party_id"`
}

// String stringify player position
func (pPostion PlayerPosition) String() string {
	str := fmt.Sprintf("Current angle : %.2f\n", pPostion.CurrentAngle)
	str += pPostion.CurrentPosition.String()
	return str
}

// String stringify player input
func (pInput *PlayerInput) String() string {
	var str string
	str += "Acceleration : " + fmt.Sprintf("%.2f", pInput.Acceleration) + "\n"
	str += "Turning : " + fmt.Sprintf("%.2f", pInput.Turning) + "\n"
	str += "Sending date : " + pInput.Timestamp.Format(time.RFC3339)
	return str
}

// String stringify a player representation
func (p Player) String() string {
	str := "Player name : " + p.PlayerName + "\n"
	str += p.Position.String()
	return str
}

// NewPlayer create a player instance from a given name
func NewPlayer(name string) *Player {
	player := new(Player)
	player.PlayerName = name
	player.PlayerUUID = uuid.New()
	player.Position = new(PlayerPosition)
	player.Input = new(PlayerInput)
	return player
}
