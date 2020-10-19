package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestNewPlayer(t *testing.T) {
	token := PlayerCreationToken{
		SessionUUID: uuid.New(),
		PlayerName:  "toto",
	}
	bitifyToken, err := json.Marshal(token)
	if err != nil {
		t.Fatal("could not marshal player token :", err)
	}
	fmt.Println(string(bitifyToken))

	player := NewPlayer(token.PlayerName)
	bitifyPlayer, err := json.Marshal(player)
	if err != nil {
		t.Fatal("could not marshal player :", err)
	}
	fmt.Println(string(bitifyPlayer))

	player.Input = &PlayerInput{
		Acceleration:  0,
		Turning:       0,
		MessageNumber: 0,
		Timestamp:     time.Now(),
		PlayerUUID:    player.PlayerUUID,
	}
	bitifyPlayerInput, err := json.Marshal(player.Input)
	if err != nil {
		t.Fatal("could not marshal player input :", err)
	}
	fmt.Println(string(bitifyPlayerInput))
}