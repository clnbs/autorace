package models

import (
	"github.com/clnbs/autorace/internal/pkg/mathtool"
	"math"
)

// SplineChain creates a spline curve through a series of control points by
// chaining Catmull-Rom splines.
//
// Each curve segment for two control points will have pointsPerSegment
// points, including the control points. Control points are not duplicated
// at the seams of the segments.
// Each control point is included in the result exactly as given.
//
// The alpha value ranges from 0 to 1. An alpha value of 0.5 results in a
// centripetal spline, alpha=0 results in a uniform spline, and alpha=1 results
// in a chordal spline.
func SplineChain(turnPoints []TurnPoint, pointsPerSegment int, alpha float64) []TurnPoint {
	newTurnPoints := make([]TurnPoint, len(turnPoints)+2)
	copy(newTurnPoints[1:], turnPoints)

	turnPoint0 := turnPoints[0]
	turnPoint1 := turnPoints[1]
	// Additional extrapolated control point at the beginning
	newTurnPoints[0] = TurnPoint{
		Position: turnPoint0.Position.Subtract(turnPoint1.Position.Subtract(turnPoint0.Position)),
	}

	turnPointY := turnPoints[len(turnPoints)-2]
	turnPointZ := turnPoints[len(turnPoints)-1]
	// Additional extrapolated control point at the end
	newTurnPoints[len(newTurnPoints)-1] = TurnPoint{Position: turnPointZ.Position.Add(turnPointZ.Position.Subtract(turnPointY.Position))}

	return chain(newTurnPoints, pointsPerSegment, alpha)
}

func chain(controlPoints []TurnPoint, pointsPerSegment int, alpha float64) []TurnPoint {
	P := controlPoints
	nSegments := len(P) - 3
	curve := make([]TurnPoint, 0, nSegments*pointsPerSegment-(nSegments-1))
	for i := 0; i < nSegments; i++ {
		segment := Spline(P[i], P[i+1], P[i+2], P[i+3], pointsPerSegment, alpha)
		if i == 0 {
			curve = append(curve, segment...)
		} else {
			// Do not duplicate points at seams.
			curve = append(curve, segment[1:]...)
		}
	}
	return curve
}

// Spline calculates the Catmull-Rom spline curve defined by the points p0, p1,
// p2, p3. The resulting curve starts with p1 and ends with p2. Both these
// points are included exactly as given. The total number of points for the
// resulting curve is defined by nPoints.
//
// The alpha value ranges from 0 to 1. An alpha value of 0.5 results in a
// centripetal spline, alpha=0 results in a uniform spline, and alpha=1 results
// in a chordal spline.
func Spline(p0, p1, p2, p3 TurnPoint, nPoints int, alpha float64) []TurnPoint {

	tj := func(ti float64, pi, pj TurnPoint) float64 {
		return math.Pow(mathtool.Distance(pi.Position, pj.Position), alpha) + ti
	}

	t0 := float64(0)
	t1 := tj(t0, p0, p1)
	t2 := tj(t1, p1, p2)
	t3 := tj(t2, p2, p3)

	step := (t2 - t1) / float64(nPoints-1)

	spline := make([]TurnPoint, nPoints)
	spline[0] = p1
	for i := 1; i < nPoints-1; i++ {

		t := t1 + (float64(i) * step)

		a1 := p0.Position.SetLength((t1 - t) / (t1 - t0)).Add(p1.Position.SetLength((t - t0) / (t1 - t0)))
		a2 := p1.Position.SetLength((t2 - t) / (t2 - t1)).Add(p2.Position.SetLength((t - t1) / (t2 - t1)))
		a3 := p2.Position.SetLength((t3 - t) / (t3 - t2)).Add(p3.Position.SetLength((t - t2) / (t3 - t2)))

		b1 := a1.SetLength((t2 - t) / (t2 - t0)).Add(a2.SetLength((t - t0) / (t2 - t0)))
		b2 := a2.SetLength((t3 - t) / (t3 - t1)).Add(a3.SetLength((t - t1) / (t3 - t1)))

		c := b1.SetLength((t2 - t) / (t2 - t1)).Add(b2.SetLength((t - t1) / (t2 - t1)))
		spline[i] = TurnPoint{Position: c}
	}
	spline[nPoints-1] = p2
	return spline
}
