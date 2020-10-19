package models

import (
	"math"
	"sort"

	"github.com/clnbs/autorace/internal/pkg/mathtool"
)

//MapGeneration cal every function to generate a racetrack in 8 steps :
// - Generate a cloud of dots
// - Get outsider dot using convex hull algorithm
// - Push point appart so they are not too close x3
// - Remove too sharp turns in order to have a clean racetrack
// - add some difficulties in the racetrack, convex hull give a too round racetrack
// - push apart dot in order to have a clean track
// - remove loop in track (because it's a 2D map)
func (partyMap *PartyMap) MapGeneration(config CircuitMapConfig) {
	pointCount := mathtool.RandomIntBetween(config.MinPoint, config.MaxPoint)
	partyMap.TurnPoints = make([]TurnPoint, pointCount)
	for i := 0; i < pointCount; i++ {
		partyMap.TurnPoints[i].Position.X = mathtool.RandomFloatBetween(0, config.XSize) - config.XSize/2
		partyMap.TurnPoints[i].Position.Y = mathtool.RandomFloatBetween(0, config.YSize) - config.YSize/2
	}
	partyMap.TurnPoints = partyMap.ConvexHull()

	for i := 0; i < 3; i++ {
		partyMap.PushApart()
	}

	partyMap.TurnPoints = partyMap.RemoveTooSharpTurn()
	partyMap.TurnPoints = partyMap.RacetrackDifficulty(0.02)

	for i := 0; i < 3; i++ {
		partyMap.PushApart()
	}
	//
	partyMap.TurnPoints = partyMap.RemoveLoop()
	partyMap.TurnPoints = partyMap.RemoveTooSharpTurn()
	partyMap.TurnPoints = SplineChain(partyMap.TurnPoints, 100, 1)
}

//ConvexHull select dot in a cloud of dots in order to get a circle like shape
func (partyMap *PartyMap) ConvexHull() []TurnPoint {
	sort.Sort(partyMap)
	var upper, lower = make([]TurnPoint, 0, len(partyMap.TurnPoints)), make([]TurnPoint, 0, len(partyMap.TurnPoints))

	// Build lower hull
	for i := 0; i < len(partyMap.TurnPoints); i++ {
		lowerLen := len(lower)
		for lowerLen >= 2 && mathtool.Cross(lower[lowerLen-2].Position, lower[lowerLen-1].Position, partyMap.TurnPoints[i].Position) <= 0 {
			lower = lower[:len(lower)-1]
			lowerLen--
		}
		lower = append(lower, partyMap.TurnPoints[i])
	}

	// Build upper hull
	for i := len(partyMap.TurnPoints) - 2; i >= 0; i-- {
		upperLen := len(upper)
		for upperLen >= 2 && mathtool.Cross(upper[upperLen-2].Position, upper[upperLen-1].Position, partyMap.TurnPoints[i].Position) <= 0 {
			upper = upper[:len(upper)-1]
			upperLen--
		}
		upper = append(upper, partyMap.TurnPoints[i])
	}
	convexhull := append(lower, upper...)

	//Sometime convexhull may select the same dot twice so we ensure that every points are unique except the firt and the last
	convexhullUnique := []TurnPoint{}
	for _, cp := range convexhull {
		seen := false
		for _, uniqueCp := range convexhullUnique {
			if uniqueCp.Position.Equals(cp.Position) {
				seen = true
				break
			}
		}
		if !seen {
			convexhullUnique = append(convexhullUnique, cp)
		}
	}
	convexhullUnique = append(convexhullUnique, convexhull[len(convexhull)-1])
	return convexhullUnique
}

//PushApart make sure dots are not too close to each other, otherwise, it moves them to the outer circle of the racetrack.
//It moves dots by adding two normalized vectors : B to A and centroid to A, where A and B are too close
func (partyMap *PartyMap) PushApart() {
	centroid := partyMap.getCentroid()
	minimalDistance := 500.0
	for index := 0; index < len(partyMap.TurnPoints); index++ {
		for indexToCompare := index + 1; indexToCompare < len(partyMap.TurnPoints); indexToCompare++ {
			pointA, pointB := partyMap.TurnPoints[index].Position, partyMap.TurnPoints[indexToCompare].Position
			if mathtool.Distance(pointA, pointB) < minimalDistance {
				vectorCentroidA := mathtool.GetNormalizedDirection(centroid, pointA)
				vectorCentroidB := mathtool.GetNormalizedDirection(centroid, pointB)
				// if A, B and centroid are aligned -- should not happen -- we move one point to the outer circle
				// and the other to the inner circle
				if vectorCentroidA.Equals(vectorCentroidB) {
					angle := vectorCentroidA.GetVectorAngle()

					partyMap.TurnPoints[index].Position.X = partyMap.TurnPoints[index].Position.X +
						(minimalDistance/2)*math.Cos(angle)
					partyMap.TurnPoints[index].Position.Y = partyMap.TurnPoints[index].Position.Y +
						(minimalDistance/2)*math.Sin(angle)

					partyMap.TurnPoints[indexToCompare].Position.X = partyMap.TurnPoints[indexToCompare].Position.X -
						(minimalDistance/2)*math.Cos(angle)
					partyMap.TurnPoints[indexToCompare].Position.Y = partyMap.TurnPoints[indexToCompare].Position.Y -
						(minimalDistance/2)*math.Sin(angle)
					continue
				}
				// we randomized which point are going to be moved around
				if mathtool.RandomIntBetween(0, 1) == 0 {
					vectorBtoA := mathtool.GetNormalizedDirection(pointB, pointA)
					movingDirectionA := vectorBtoA.Add(vectorCentroidA)
					movingDirectionA = movingDirectionA.Normalized()
					movingDirectionA = movingDirectionA.SetLength(minimalDistance)
					partyMap.TurnPoints[index].Position.Move(movingDirectionA)
					continue
				}
				vectorAtoB := mathtool.GetNormalizedDirection(pointA, pointB)
				movingDirectionB := vectorAtoB.Add(vectorCentroidB)
				movingDirectionB = movingDirectionB.Normalized()
				movingDirectionB = movingDirectionB.SetLength(minimalDistance)
				partyMap.TurnPoints[indexToCompare].Position.Move(movingDirectionB)
			}
		}
	}
	// as long as the first and the last point are the same, we have to re-register them
	partyMap.TurnPoints[len(partyMap.TurnPoints)-1] = partyMap.TurnPoints[0]
}

//RemoveTooSharpTurn remove checkpoint when they are in an impossible position
// We check dot three per three and calculate the angle they form. If it less than a threshold, we remove the 3rd dot
func (partyMap *PartyMap) RemoveTooSharpTurn() []TurnPoint {
	//if thy are less than 3 points, ends it -- should not happen
	if len(partyMap.TurnPoints) < 3 {
		return partyMap.TurnPoints
	}
	// we register point to remove in a map in order to not modify the slice in the loop
	toRemove := make(map[TurnPoint]bool)
	for index := 0; index < len(partyMap.TurnPoints)-2; index++ {
		//register our 3 dots
		a := partyMap.TurnPoints[index].Position
		b := partyMap.TurnPoints[index+1].Position
		c := partyMap.TurnPoints[index+2].Position
		//calculate angles
		vectorAtoB := mathtool.GetNormalizedDirection(a, b)
		vectorBtoC := mathtool.GetNormalizedDirection(b, c)
		vectorCtoB := mathtool.GetNormalizedDirection(c, b)
		vectorBtoA := mathtool.GetNormalizedDirection(b, a)
		ABCAngle1 := mathtool.GetVectorsAngle(vectorAtoB, vectorBtoC)
		ABCAngle2 := mathtool.GetVectorsAngle(vectorAtoB, vectorCtoB)
		ABCAngle3 := mathtool.GetVectorsAngle(vectorBtoA, vectorBtoC)
		ABCAngle4 := mathtool.GetVectorsAngle(vectorBtoA, vectorCtoB)

		//remove the c dot if it less than an arbitrary angle
		if ABCAngle1 < math.Pi/7 || ABCAngle2 < math.Pi/7 || ABCAngle3 < math.Pi/7 || ABCAngle4 < math.Pi/7 {
			toRemove[partyMap.TurnPoints[index+1]] = true
		}
	}
	//we register dots in a new slice except if they are in the delete list
	var newTurnPoints []TurnPoint
	for _, tp := range partyMap.TurnPoints {
		if _, ok := toRemove[tp]; !ok {
			newTurnPoints = append(newTurnPoints, tp)
		}
	}
	if len(newTurnPoints) == 0 {
		return partyMap.TurnPoints
	}
	if newTurnPoints[0] != newTurnPoints[len(newTurnPoints)-1] {
		newTurnPoints = append(newTurnPoints, newTurnPoints[0])
	}
	return newTurnPoints
}

//RemoveLoop remove loop in racetrack
func (partyMap *PartyMap) RemoveLoop() []TurnPoint {
	toRemove := make(map[TurnPoint]bool)
	for indexFirstSegment := 0 ; indexFirstSegment < len(partyMap.TurnPoints)-4; indexFirstSegment++ {
		pointA, pointB := partyMap.TurnPoints[indexFirstSegment].Position, partyMap.TurnPoints[indexFirstSegment+1].Position
		for indexSecondSegment := indexFirstSegment+2; indexSecondSegment < len(partyMap.TurnPoints)-2; indexSecondSegment++ {
			pointC, pointD := partyMap.TurnPoints[indexSecondSegment].Position, partyMap.TurnPoints[indexSecondSegment+1].Position
			if mathtool.DoIntersect(pointA, pointB, pointC, pointD) {
				for indexToDelete := indexFirstSegment+1; indexToDelete < indexSecondSegment+1; indexToDelete++ {
					toRemove[partyMap.TurnPoints[indexToDelete]] = true
				}
			}
		}
	}
	var newTurnPoints []TurnPoint
	for _, tp := range partyMap.TurnPoints {
		if _, ok := toRemove[tp]; !ok {
			newTurnPoints = append(newTurnPoints, tp)
		}
	}
	if len(newTurnPoints) == 0 {
		return partyMap.TurnPoints
	}
	if newTurnPoints[0] != newTurnPoints[len(newTurnPoints)-1] {
		newTurnPoints = append(newTurnPoints, newTurnPoints[0])
	}
	return newTurnPoints
}

//RacetrackDifficulty spice up racetrack since Convex Hull give a too easy rounded shape.
// For each racetrack segment (between two turning point) longer than a given length, we add a "random" middle dot
func (partyMap *PartyMap) RacetrackDifficulty(spice float64) []TurnPoint {
	var raceTrack []TurnPoint
	//max distance to push the new turning point from current racetrack line
	maxDistance := 600.0
	//minimal distance between two point to break
	distanceToBreak := 500.0
	//we add the first turn point to the new racetrack
	raceTrack = append(raceTrack, partyMap.TurnPoints[0])
	for index := 1; index < len(partyMap.TurnPoints); index++ {
		pointA, pointB := partyMap.TurnPoints[index-1].Position, partyMap.TurnPoints[index].Position
		if mathtool.Distance(pointA, pointB) > distanceToBreak {
			// get the new randomized distance to push the new point
			newDistance := math.Pow(mathtool.RandomFloatBetween(0, 1), spice) * maxDistance
			var newTurnPoint TurnPoint
			// currently, the position of the new turn point is exactly in the middle of the line, may change in the futur
			newTurnPoint.Position.X = (pointA.X + pointB.X) / 2
			newTurnPoint.Position.Y = (pointA.Y + pointB.Y) / 2
			// we give to the turning point a random "angle" to be push to
			newTurnPoint.Position = newTurnPoint.Position.Rotate(mathtool.RandomFloatBetween(0, math.Pi*2))
			// then we scale it to the given angle
			newTurnPoint.Position = newTurnPoint.Position.ScaleWithAngle(newDistance)
			raceTrack = append(raceTrack, newTurnPoint)
		}
		raceTrack = append(raceTrack, partyMap.TurnPoints[index])
	}
	return raceTrack
}

func (partyMap *PartyMap) getCentroid() mathtool.Vector2 {
	var positionList []mathtool.Vector2
	for _, cp := range partyMap.TurnPoints {
		positionList = append(positionList, cp.Position)
	}
	return mathtool.GetCentroid(positionList[:len(positionList)-1])
}
