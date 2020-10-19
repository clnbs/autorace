package database

import (
	"fmt"
	"testing"
	"time"

	"github.com/clnbs/autorace/internal/app/models"
	"github.com/clnbs/autorace/pkg/systool"
)

func TestRedisClient_GetPartyList(t *testing.T) {
	defer fmt.Println(systool.TimeTrack(time.Now(), "TestRedisClient_GetPartyList"))
	rdsClient := NewRedisClient()
	err := rdsClient.SetPartyConfiguration("one", models.PartyCreationToken{
		ClientID:      "one",
		Seed:          0,
		PartyName:     "one",
		CircuitConfig: models.CircuitMapConfig{
			Seed:     0,
			MaxPoint: 100,
			MinPoint: 10,
			XSize:    100,
			YSize:    100,
		},
	})
	if err != nil {
		t.Fatal("error while inserting data in redis :", err)
	}
	err = rdsClient.SetPartyConfiguration("two", models.PartyCreationToken{
		ClientID:      "two",
		Seed:          0,
		PartyName:     "two",
		CircuitConfig: models.CircuitMapConfig{
			Seed:     0,
			MaxPoint: 100,
			MinPoint: 10,
			XSize:    100,
			YSize:    100,
		},
	})
	if err != nil {
		t.Fatal("error while inserting data in redis :", err)
	}
	err = rdsClient.SetPartyConfiguration("three", models.PartyCreationToken{
		ClientID:      "three",
		Seed:          0,
		PartyName:     "three",
		CircuitConfig: models.CircuitMapConfig{
			Seed:     0,
			MaxPoint: 100,
			MinPoint: 10,
			XSize:    100,
			YSize:    100,
		},
	})
	if err != nil {
		t.Fatal("error while inserting data in redis :", err)
	}
	partyList, err := rdsClient.GetPartyList()
	if err != nil {
		t.Fatal("error while getting party list :", err)
	}
	for _, party := range partyList {
		fmt.Println(party)
	}

	partyTokenOne, err := rdsClient.GetPartyCreationToken("one")
	if err != nil {
		t.Fatal("could not get party token one :", err)
	}
	fmt.Println(partyTokenOne)
	fmt.Println(partyTokenOne.CircuitConfig)
}

func TestRedisClient_GetPlayerOnParty(t *testing.T) {
	defer fmt.Println(systool.TimeTrack(time.Now(), "TestRedisClient_GetPlayerOnParty"))
	rdsClient := NewRedisClient()
	err := rdsClient.SetPlayerOnParty("one", "toto")
	if err != nil {
		t.Fatal("error while setting player on party :", err)
	}
	err = rdsClient.SetPlayerOnParty("one", "tata")
	if err != nil {
		t.Fatal("error while setting player on party :", err)
	}
	err = rdsClient.SetPlayerOnParty("one", "titi")
	if err != nil {
		t.Fatal("error while setting player on party :", err)
	}
	err = rdsClient.SetPlayerOnParty("one", "tutu")
	if err != nil {
		t.Fatal("error while setting player on party :", err)
	}
	players, err := rdsClient.GetPlayersOnParty("one")
	if err != nil {
		t.Fatal("error while getting player on party :", err)
	}
	for _, p := range players {
		fmt.Println(p)
	}
}
