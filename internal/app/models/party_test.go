package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func TestNewParty(t *testing.T) {
	party := new(Party)
	party.PartyName = "testing"
	party.PartyUUID = uuid.New()
	party.Players = make(map[string]*Player)
	party.CircuitConfig = CircuitMapConfig{
		Seed:     0,
		MaxPoint: 20,
		MinPoint: 10,
		XSize:    2000,
		YSize:    2000,
	}
	newPlayer := new(Player)
	newPlayer.PlayerUUID = uuid.New()
	newPlayer.PlayerName = "testing"
	party.Players[newPlayer.PlayerUUID.String()] = newPlayer
	party.MapCircuit.MapGeneration(party.CircuitConfig)

	tokenCreation := PartyCreationToken{
		ClientID:      uuid.New().String(),
		Seed:          321,
		PartyName:     "toto party",
		CircuitConfig: CircuitMapConfig{
			Seed:     321,
			MaxPoint: 250,
			MinPoint: 50,
			XSize:    4000,
			YSize:    4000,
		},
	}

	jsonToken, err := json.Marshal(tokenCreation)
	fmt.Println("token :", string(jsonToken))
	fmt.Println("##########################################")

	jsonParty, err := json.Marshal(party)
	if err != nil {
		t.Fatal("error while marshaling party :", err)
	}
	fmt.Println("party :", string(jsonParty))
	fmt.Println("##########################################")
	partyList := []string{
		"partyID_1",
		"partyID_2",
		"partyID_3",
		"partyID_4",
	}
	jsonPartyList, err := json.Marshal(partyList)
	fmt.Println("party list :", string(jsonPartyList))
}
