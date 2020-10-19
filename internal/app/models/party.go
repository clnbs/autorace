package models

import (
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var (
	// ErrorPlayerAlreadyInParty used to trigger an error
	ErrorPlayerAlreadyInParty = errors.New("player is already registered in party")
	// ErrorPlayerNotFound used to trigger an error
	ErrorPlayerNotFound = errors.New("player not found in party")
)

// State is used to represent game's state in a Enum style
type State int

// possible game state
const (
	// LOBBY is when a party still accept player
	LOBBY State = iota
	// RUN the party is running
	RUN
	// PAUSE the party is paused
	PAUSE
	// END the party has ended
	END
)

var states = [...]string{
	"Lobby",
	"Run",
	"Pause",
	"End",
}

// String stringify possible state
func (state State) String() string {
	if state < LOBBY || state > END {
		return "unknown state"
	}
	return states[state]
}

// NewPartyState transform a string into a possible state
func NewPartyState(state string) State {
	formatedState := strings.ToLower(state)
	possibleState := map[string]State{
		"lobby": LOBBY,
		"run":   RUN,
		"pause": PAUSE,
		"end":   END,
	}
	if _, ok := possibleState[state]; !ok {
		return RUN
	}
	return possibleState[formatedState]
}

// PartyCreationToken is used to ask server to start a party room (in a dynamic server instance)
type PartyCreationToken struct {
	ClientID      string           `json:"client_id"`
	Seed          int              `json:"seed"`
	PartyName     string           `json:"party_name"`
	CircuitConfig CircuitMapConfig `json:"circuit_config"`
}

// String stringify PartyCreationToken
func (clientToken PartyCreationToken) String() string {
	str := "client ID : " + clientToken.ClientID + "\n"
	str += "seed : " + strconv.FormatInt(int64(clientToken.Seed), 10) + "\n"
	str += "party name : " + clientToken.PartyName
	return str
}

// ChangeStateToken is send to a party instance (aka dynamic server) to change game state
type ChangeStateToken struct {
	PlayerToken  PlayerToken `json:"player_token"`
	DesiredState State       `json:"desired_state"`
}

// ChangeStateAck is send back to player who asked to change game state
type ChangeStateAck struct {
	PartyID      string `json:"party_id"`
	DesiredState State  `json:"desired_state"`
	NewState     State  `json:"new_state"`
	Message      string `json:"message,omitempty"`
}

// Party is a representation of party
type Party struct {
	PartyUUID uuid.UUID `json:"party_uuid"`
	PartyName string    `json:"party_name"`
	//bind players by players' uuid in order to lower complexity when looking for a particular player
	Players       map[string]*Player `json:"-"`
	MapCircuit    PartyMap           `json:"map_circuit"`
	CircuitConfig CircuitMapConfig   `json:"circuit_config"`
	state         State
}

// String stringify a Party
func (party *Party) String() string {
	str := "Party name : " + party.PartyName + "\n"
	str += "Party state : " + party.state.String() + "\n"
	str += "Party UUID : " + party.PartyUUID.String()
	return str
}

// NewParty create a Party from a party creation token
func NewParty(creationToken PartyCreationToken, stringifyUUID string) (*Party, error) {
	party := new(Party)
	var err error
	party.PartyName = creationToken.PartyName
	party.PartyUUID, err = uuid.Parse(stringifyUUID)
	if err != nil {
		return nil, err
	}
	party.Players = make(map[string]*Player)
	party.CircuitConfig = creationToken.CircuitConfig
	return party, nil
}

// AddPlayer is use to add a player in a party
func (party *Party) AddPlayer(player *Player) error {
	if _, ok := party.Players[player.PlayerUUID.String()]; ok {
		return ErrorPlayerAlreadyInParty
	}
	party.Players[player.PlayerUUID.String()] = player
	return nil
}

// RemovePlayer is use to remove player from a party
func (party *Party) RemovePlayer(player *Player) error {
	if _, ok := party.Players[player.PlayerUUID.String()]; !ok {
		return ErrorPlayerNotFound
	}
	delete(party.Players, player.PlayerUUID.String())
	return nil
}

// RemoveAllPlayer delete all registered player in a party
func (party *Party) RemoveAllPlayer() {
	party.Players = make(map[string]*Player)
}

// SetState change party's state
func (party *Party) SetState(s State) {
	party.state = s
}

// GetState returns party's state
func (party *Party) GetState() State {
	return party.state
}

// AddPlayerToken is used by client to be added in a party
type AddPlayerToken struct {
	PlayerUUID string `json:"player_uuid"`
	PartyUUID  string `json:"party_uuid"`
}
