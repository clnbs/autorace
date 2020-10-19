package server

import (
	"github.com/clnbs/autorace/pkg/logger"
	"github.com/clnbs/autorace/pkg/systool"
	"math"
	"time"

	"github.com/clnbs/autorace/internal/app/models"
	"github.com/clnbs/autorace/internal/pkg/mathtool"
)

var (
	maxSpeed                = 500.0
	minSpeed                = -200.0
	decelerationFactor      = 0.2
	acceleration            = 500.0
	turningAcceleration     = 2.0
	maxGrassSpeed           = 200.0
	minGrassSpeed           = -50.0
	numberOfPositionToCheck = 500000
)

func (dServer *DynamicPartyServer) computeNewPosition(deltaTime float64) {
	logger.Trace(systool.TimeTrack(time.Now(), "compute players position"))
	for _, player := range dServer.party.Players {
		dServer.computeNewPlayerSpeed(player, deltaTime)
		dServer.computeNewPlayerAngle(player, deltaTime)
		player.Position.CurrentPosition.Y += (player.Position.CurrentSpeed * deltaTime) * math.Sin(player.Position.CurrentAngle)
		player.Position.CurrentPosition.X += (player.Position.CurrentSpeed * deltaTime) * math.Cos(player.Position.CurrentAngle)
	}
}

func (dServer *DynamicPartyServer) computeNewPlayerAngle(p *models.Player, deltaTime float64) {
	if p.Position.CurrentSpeed == 0 {
		return
	}
	p.Position.CurrentAngle += p.Input.Turning * deltaTime * turningAcceleration
}

func (dServer *DynamicPartyServer) computeNewPlayerSpeed(p *models.Player, deltaTime float64) {
	dServer.computeClosestRacetrackPointIndex(p)
	closestRacetrackPoint := dServer.closestRacetrackPointIndex[p.PlayerUUID.String()]
	if mathtool.Distance(p.Position.CurrentPosition, dServer.party.MapCircuit.TurnPoints[closestRacetrackPoint].Position) > 50 &&
		!mathtool.IsFloat64Between(p.Position.CurrentSpeed, minGrassSpeed, maxGrassSpeed) {
		dServer.computeDeceleration(p, 10.0)
	}
	if p.Input.Acceleration == 0 {
		dServer.computeDeceleration(p, 1.0)
		return
	}
	p.Position.CurrentSpeed += p.Input.Acceleration * deltaTime * acceleration
	p.Position.CurrentSpeed = mathtool.ClampFloat64(p.Position.CurrentSpeed, minSpeed, maxSpeed)
}

func (dServer *DynamicPartyServer) computeDeceleration(p *models.Player, multiplicatorFactor float64) {
	absSpeed := math.Abs(p.Position.CurrentSpeed)
	if absSpeed < 0.2 {
		p.Position.CurrentSpeed = 0
		return
	}
	decelerationDuration := math.Sqrt(absSpeed)
	decelerationSpeed := absSpeed / decelerationDuration
	if p.Position.CurrentSpeed > 0 {
		p.Position.CurrentSpeed -= decelerationSpeed * decelerationFactor * multiplicatorFactor
		return
	}
	p.Position.CurrentSpeed += decelerationSpeed * decelerationFactor * multiplicatorFactor
}

func (dServer *DynamicPartyServer) setCarAtStart() {
	startPosition := dServer.party.MapCircuit.TurnPoints[0].Position
	startAngle := mathtool.GetNormalizedDirection(dServer.party.MapCircuit.TurnPoints[0].Position, dServer.party.MapCircuit.TurnPoints[100].Position).GetVectorAngle() + math.Pi
	for _, player := range dServer.party.Players {
		player.Position.CurrentPosition = startPosition
		player.Position.CurrentAngle = startAngle
	}
}

func (dServer *DynamicPartyServer) computeClosestRacetrackPointIndex(player *models.Player) {
	var startIndex, endIndex, closestIndex int
	smallestDistance := math.MaxFloat64
	closestIndex = dServer.closestRacetrackPointIndex[player.PlayerUUID.String()]
	indexBoundaries := dServer.searchIndexBoundaries(closestIndex)
	playerPosition := player.Position.CurrentPosition
	for len(indexBoundaries) != 0 {
		startIndex, endIndex, indexBoundaries = indexBoundaries[0], indexBoundaries[1], indexBoundaries[2:]
		for index := startIndex; index <= endIndex; index++ {
			indexPosition := dServer.party.MapCircuit.TurnPoints[index].Position
			if mathtool.Distance(playerPosition, indexPosition) < smallestDistance {
				smallestDistance = mathtool.Distance(playerPosition, indexPosition)
				closestIndex = index
			}
		}
	}
	dServer.closestRacetrackPointIndex[player.PlayerUUID.String()] = closestIndex
}

func (dServer *DynamicPartyServer) searchIndexBoundaries(currentIndex int) []int {
	var indexBoundaries []int
	var realNumberOfPostionToCheck int
	lastIndexOfRacetrack := len(dServer.party.MapCircuit.TurnPoints) - 1
	if numberOfPositionToCheck > lastIndexOfRacetrack/2 {
		realNumberOfPostionToCheck = lastIndexOfRacetrack/2
	} else {
		realNumberOfPostionToCheck = numberOfPositionToCheck
	}
	if currentIndex-realNumberOfPostionToCheck < 0 {
		indexBoundaries = append(indexBoundaries, 0, currentIndex+realNumberOfPostionToCheck)
		indexBoundaries = append(indexBoundaries, lastIndexOfRacetrack+currentIndex-realNumberOfPostionToCheck, lastIndexOfRacetrack)
		return indexBoundaries
	}
	if currentIndex+realNumberOfPostionToCheck > lastIndexOfRacetrack {
		indexBoundaries = append(indexBoundaries, 0, currentIndex+realNumberOfPostionToCheck-lastIndexOfRacetrack)
		indexBoundaries = append(indexBoundaries, currentIndex-realNumberOfPostionToCheck, lastIndexOfRacetrack)
		return indexBoundaries
	}
	indexBoundaries = append(indexBoundaries, currentIndex-realNumberOfPostionToCheck, currentIndex+realNumberOfPostionToCheck)
	return indexBoundaries
}
