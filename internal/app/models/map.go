package models

import (
	"fmt"

	"github.com/clnbs/autorace/internal/pkg/mathtool"
)

//PartyMap is the representation of a racetrack
type PartyMap struct {
	TurnPoints []TurnPoint `json:"turnpoints"`
}

//TurnPoint are generated points in order to trace a racetrack
type TurnPoint struct {
	Position mathtool.Vector2 `json:"position"`
}

func (tp TurnPoint) String() string {
	return fmt.Sprintf("X : %f \t Y : %f", tp.Position.X, tp.Position.Y)
}

//CircuitMapConfig contain configuration to generate racetrack
type CircuitMapConfig struct {
	Seed     int     `json:"seed"`
	MaxPoint int     `json:"max_point"`
	MinPoint int     `json:"min_point"`
	XSize    float64 `json:"x_size"`
	YSize    float64 `json:"y_size"`
}

// String stringify a Circuit configuration
func (circuitMC CircuitMapConfig) String() string {
	str := "Maximum point number :" + fmt.Sprintf("%d", circuitMC.MaxPoint) + "\n"
	str += "Minimum point number :" + fmt.Sprintf("%d", circuitMC.MinPoint) + "\n"
	str += "X size :" + fmt.Sprintf("%f", circuitMC.XSize) + "\n"
	str += "Y size :" + fmt.Sprintf("%f", circuitMC.YSize) + "\n"
	return str
}

//Len implemented in order to call sort.Sort()
func (partyMap PartyMap) Len() int {
	return len(partyMap.TurnPoints)
}

//Swap implemented in order to call sort.Sort()
func (partyMap PartyMap) Swap(i, j int) {
	partyMap.TurnPoints[i], partyMap.TurnPoints[j] = partyMap.TurnPoints[j], partyMap.TurnPoints[i]
}

//Less implemented in order to call sort.Sort()
func (partyMap PartyMap) Less(i, j int) bool {
	if partyMap.TurnPoints[i].Position.X == partyMap.TurnPoints[j].Position.X {
		return partyMap.TurnPoints[i].Position.Y < partyMap.TurnPoints[j].Position.Y
	}
	return partyMap.TurnPoints[i].Position.X < partyMap.TurnPoints[j].Position.X
}
