package server

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/clnbs/autorace/internal/app/models"
	"github.com/clnbs/autorace/internal/pkg/mathtool"
	"testing"
)

func TestSyncMessageContent(t *testing.T) {
	competitors := []*models.CompetitorActor{
		&models.CompetitorActor{
			Act:       nil,
			ActorUUID: uuid.New(),
			Position:  &models.PlayerPosition{
				CurrentSpeed:    0,
				CurrentAngle:    0,
				CurrentPosition: mathtool.Vector2{
					X:     0,
					Y:     0,
					Angle: 0,
				},
			},
		},
		&models.CompetitorActor{
			Act:       nil,
			ActorUUID: uuid.New(),
			Position:  &models.PlayerPosition{
				CurrentSpeed:    0,
				CurrentAngle:    0,
				CurrentPosition: mathtool.Vector2{
					X:     0,
					Y:     0,
					Angle: 0,
				},
			},
		},
		&models.CompetitorActor{
			Act:       nil,
			ActorUUID: uuid.New(),
			Position:  &models.PlayerPosition{
				CurrentSpeed:    0,
				CurrentAngle:    0,
				CurrentPosition: mathtool.Vector2{
					X:     0,
					Y:     0,
					Angle: 0,
				},
			},
		},
	}
	syncMessage := SyncMessageContent{
		PartyState:  models.RUN,
		Competitors: competitors,
		MainActor:   &models.MainActor{
			Act:    nil,
			Player: models.NewPlayer("toto"),
		},
	}

	jsonSyncMessage, err := json.Marshal(syncMessage)
	if err != nil {
		t.Fatal("could not marshal sync message :", err)
	}
	fmt.Println("sync message :", string(jsonSyncMessage))
	fmt.Println("##########################################")

	stateAck := models.ChangeStateAck{
		PartyID:      "partyID",
		DesiredState: 2,
		NewState:     2,
		Message:      "OK",
	}
	jsonStateAck, err := json.Marshal(stateAck)
	if err != nil {
		t.Fatal("could not marshal state ack :", err)
	}
	fmt.Println("state ack :", string(jsonStateAck))
}