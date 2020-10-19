package models

import (
	"strconv"

	"github.com/google/uuid"
)

// MainActor is main player representation. It contains all player's information
// and in game representation
type MainActor struct {
	Act    *Actor  `json:"actor"`
	Player *Player `json:"player"`
}

// CompetitorActor is competitors' representation.
type CompetitorActor struct {
	Act       *Actor          `json:"actor"`
	ActorUUID uuid.UUID       `json:"actor_uuid"`
	Position  *PlayerPosition `json:"position"`
}

// Actor contains in-game representation
type Actor struct {
	Car  *Car
	Name string `json:"name"`
	Rank int    `json:"rank"`
}

//String stringify Actor
func (act Actor) String() string {
	str := "Name : " + act.Name + "\n"
	str += "Rank : " + strconv.FormatInt(int64(act.Rank), 10) + "\n"
	return str
}

//String stringify CompetitorActor
func (cAct CompetitorActor) String() string {
	str := cAct.String()
	str += cAct.Position.String()
	return str
}
