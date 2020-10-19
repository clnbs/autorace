package mathtool

import (
	"math"
	"testing"
)

func TestNewVector(t *testing.T) {
	expected := Vector2{
		X:     0,
		Y:     0,
		Angle: 0,
	}
	v := NewVector()
	if expected.X != v.X || expected.Y != v.Y || expected.Angle != v.Angle {
		t.Fatal("NewVector() should return a vector initialized at 0")
	}
}

func TestCross(t *testing.T) {
	a := Vector2{
		X: 5,
		Y: 9,
	}
	b := Vector2{
		X: 2,
		Y: 6,
	}
	o := Vector2{
		X: 1,
		Y: 0,
	}
	// Cross computation :
	// (a.X - o.X) * (b.Y - o.Y) - (a.Y - o.Y) * (b.X - o.X)
	// = (5 - 1) * (6 - 0) - (9 - 0) * (2 - 1)
	// = 4*6 - 9
	// = 24 - 9
	// = 15
	cross := Cross(o, a, b)
	if 15 != cross {
		t.Fatal("result should be 15, got :", cross)
	}
}

func TestDistance(t *testing.T) {
	origine := NewVector()
	pointA := Vector2{
		X:     3,
		Y:     4,
		Angle: 0,
	}
	distance := Distance(origine, pointA)
	if distance != 5 {
		t.Fatal("distance between", origine.String(), "and pointA", pointA.String(), "should be 5, got", distance)
	}
}

func TestVector2_Set(t *testing.T) {
	expected := Vector2{
		X:     3,
		Y:     4,
		Angle: 0,
	}
	v := NewVector()
	v = v.Set(3, 4)
	if expected.X != v.X || expected.Y != v.Y {
		t.Fatal("vector v", v.String(), "and expected value", expected.String(), "should be equals, they are not")
	}
}

func TestVector2_Equals(t *testing.T) {
	v := Vector2{
		X:     1,
		Y:     2,
		Angle: 0,
	}
	w := Vector2{
		X:     1,
		Y:     2,
		Angle: 0,
	}
	if !v.Equals(w) {
		t.Fatal("v", v.String(), "and", w.String(), "w are equals but function tells otherwise")
	}
}

func TestVector2_Add(t *testing.T) {
	v := Vector2{
		X:     1,
		Y:     1,
		Angle: 0,
	}
	w := Vector2{
		X:     1,
		Y:     1,
		Angle: 0,
	}
	expected := Vector2{
		X:     2,
		Y:     2,
		Angle: 0,
	}
	result := v.Add(w)
	if !expected.Equals(result) {
		t.Fatal("result should be equals to", expected.String(), "got", result.String())
	}
}

func TestVector2_Divide(t *testing.T) {
	v := Vector2{
		X:     4,
		Y:     4,
		Angle: 0,
	}
	expected := Vector2{
		X:     2,
		Y:     2,
		Angle: 0,
	}
	result := v.Divide(2)
	if !expected.Equals(result) {
		t.Fatal("result should be equals to", expected.String(), "got", result.String(), "when dividing by 2")
	}
	result = result.Divide(0)
	if !expected.Equals(result) {
		t.Fatal("result should be equals to", expected.String(), "got", result.String(), "when testing dividing by 0")
	}
}

func TestVector2_Rotate(t *testing.T) {
	v := Vector2{
		X:     0,
		Y:     0,
		Angle: 0,
	}
	expected := Vector2{
		X:     0,
		Y:     0,
		Angle: math.Pi,
	}
	result := v.Rotate(math.Pi)
	if result.Angle != expected.Angle {
		t.Fatal("result angle after rotating v should be", expected.Angle, "got", result.Angle)
	}
	result = result.Rotate(math.Pi)
	expected.Angle = 2 * math.Pi
	if result.Angle != expected.Angle {
		t.Fatal("result angle after rotating v should be", expected.Angle, "got", result.Angle)
	}
}

func TestVector2_ScaleWithAngle(t *testing.T) {
	v := Vector2{
		X:     0,
		Y:     0,
		Angle: 0,
	}
	expected := Vector2{
		X:     2,
		Y:     0,
		Angle: 0,
	}
	result := v.ScaleWithAngle(2)
	if !expected.Equals(result) {
		t.Fatal("after scaling, result should be equals to", expected.String(), "got", result.String())
	}
	expected = Vector2{
		X:     0,
		Y:     2,
		Angle: 0,
	}
	v = v.Rotate(math.Pi / 2)
	result = v.ScaleWithAngle(2)
	if !expected.Equals(result) {
		t.Fatal("after scaling, result should be equals to", expected.String(), "got", result.String())
	}
}

func TestGetCentroid(t *testing.T) {
	points := make([]Vector2, 4)
	points[0] = Vector2{
		X:     0,
		Y:     0,
		Angle: 0,
	}
	points[1] = Vector2{
		X:     2,
		Y:     0,
		Angle: 0,
	}
	points[2] = Vector2{
		X:     0,
		Y:     2,
		Angle: 0,
	}
	points[3] = Vector2{
		X:     2,
		Y:     2,
		Angle: 0,
	}
	expected := Vector2{
		X:     1,
		Y:     1,
		Angle: 0,
	}
	result := GetCentroid(points)
	if !expected.Equals(result) {
		t.Fatal("centroid result should be equals to", expected.String(), "got", result.String())
	}
}

func TestVector2_GetVectorAngle(t *testing.T) {
	v := Vector2{
		X:     4,
		Y:     0,
		Angle: 0,
	}
	expectedV := 0.0
	resultV := v.GetVectorAngle()
	if resultV != expectedV {
		t.Fatal("angle from X axis should be equals to", expectedV, "got", resultV)
	}

	w := Vector2{
		X:     4,
		Y:     4,
		Angle: 0,
	}
	expectedW := math.Pi / 4
	resultW := w.GetVectorAngle()
	if resultW != expectedW {
		t.Fatal("angle from X axis should be equals to", expectedW, "got", resultW)

	}

	u := Vector2{
		X:     0,
		Y:     4,
		Angle: 0,
	}
	expectedU := math.Pi / 2
	resultU := u.GetVectorAngle()
	if resultU != expectedU {
		t.Fatal("angle from X axis should be equals to", expectedU, "got", resultU)
	}
}

func TestVector2_Length(t *testing.T) {
	v := Vector2{
		X:     1,
		Y:     0,
		Angle: 0,
	}
	expectedV := 1.0
	resultV := v.Length()
	if expectedV != resultV {
		t.Fatal("length of v", v.String(), "should be equals to", expectedV, "got", resultV)
	}

	w := Vector2{
		X:     4,
		Y:     3,
		Angle: 0,
	}
	expectedW := 5.0
	resultW := w.Length()
	if expectedV != resultV {
		t.Fatal("length of w", w.String(), "should be equals to", expectedW, "got", resultW)
	}
}

func TestVector2_Normalized(t *testing.T) {
	v := Vector2{
		X:     4,
		Y:     8,
		Angle: 0,
	}
	expectedNormalizedVLength := 1.0
	normalizedV := v.Normalized()
	resultLengthNormalizedV := normalizedV.Length()
	if resultLengthNormalizedV != expectedNormalizedVLength {
		t.Fatal("length of normalized v should be", expectedNormalizedVLength, "got", resultLengthNormalizedV)
	}
}

func TestVector2_SetLength(t *testing.T) {
	pointA := Vector2{
		X:     4,
		Y:     5,
		Angle: 0,
	}
	pointB := Vector2{
		X:     1,
		Y:     1,
		Angle: 0,
	}
	vectorAB := GetNormalizedDirection(pointA, pointB)
	resultVectorAB := vectorAB.SetLength(10.0)
	resultLengthVectorAB := resultVectorAB.Length()
	expectedLength := 10.0
	if resultLengthVectorAB != expectedLength {
		t.Fatal("length of scaled vector AB should be", expectedLength, "got", resultLengthVectorAB)
	}
}

func TestVector2_Move(t *testing.T) {
	pointA := Vector2{
		X:     2,
		Y:     2,
		Angle: 0,
	}
	vectorV := Vector2{
		X:     3,
		Y:     4,
		Angle: 0,
	}
	expectedPoint := Vector2{
		X:     5,
		Y:     6,
		Angle: 0,
	}
	resutlPoint := pointA.Move(vectorV)
	if !expectedPoint.Equals(resutlPoint) {
		t.Fatal("result of moving point A", pointA.String(), "with vector V", vectorV.String(), "should be", expectedPoint.String(), "got", resutlPoint.String())
	}
}

func TestDot(t *testing.T) {
	v := Vector2{
		X:     1,
		Y:     0,
		Angle: 0,
	}
	u := Vector2{
		X:     0,
		Y:     1,
		Angle: 0,
	}
	w := Vector2{
		X:     2,
		Y:     2,
		Angle: 0,
	}
	s := Vector2{
		X:     1,
		Y:     1,
		Angle: 0,
	}
	expectedDotVU := 0.0
	expectedDotVW := 2.0
	expectedDotUW := 2.0
	expectedDotVS := 1.0
	resultDotVU := Dot(v, u)
	resultDotVW := Dot(v, w)
	resultDotUW := Dot(u, w)
	resultDotVS := Dot(v, s)
	if expectedDotVU != resultDotVU {
		t.Fatal("scalar product of vector V", v.String(), "and vector U", u.String(), "should be equals to", expectedDotVU, "got", resultDotVU)
	}

	if expectedDotUW != resultDotUW {
		t.Fatal("scalar product of vector U", u.String(), "and vector W", w.String(), "should be equals to", expectedDotUW, "got", resultDotUW)
	}

	if expectedDotVW != resultDotVW {
		t.Fatal("scalar product of vector V", v.String(), "and vector W", w.String(), "should be equals to", expectedDotUW, "got", resultDotUW)
	}

	if expectedDotVS != resultDotVS {
		t.Fatal("scalar product of vector V", v.String(), "and vector s", s.String(), "should be equals to", expectedDotVS, "got", resultDotVS)
	}
}

func TestGetVectorsAngle(t *testing.T) {
	v := Vector2{
		X:     1,
		Y:     0,
		Angle: 0,
	}
	u := Vector2{
		X:     0,
		Y:     1,
		Angle: 0,
	}
	w := Vector2{
		X:     1,
		Y:     1,
		Angle: 0,
	}
	s := Vector2{
		X:     -1,
		Y:     0,
		Angle: 0,
	}
	v = v.Normalized()
	u = u.Normalized()
	w = w.Normalized()
	s = s.Normalized()
	expectedAngleVU := math.Pi/2
	expectedAngleVW := math.Pi/4
	expectedAngleVS := math.Pi

	expectedAngleWU := math.Pi/4
	expectedAngleUS := math.Pi/2

	expectedAngleSW := (math.Pi*3)/4

	resultAngleVU := GetVectorsAngle(v, u)
	resultAngleVW := GetVectorsAngle(v, w)
	resultAngleVS := GetVectorsAngle(v, s)

	resultAngleWU := GetVectorsAngle(w, u)
	resultAngleUS := GetVectorsAngle(u, s)

	resultAngleSW := GetVectorsAngle(s, w)

	if Round(expectedAngleVU, 2) != Round(resultAngleVU, 2) {
		t.Fatal("angle form by v", v.String(), "and u", u.String(), "should be", expectedAngleVU, "got", resultAngleVU)
	}

	if Round(expectedAngleVW, 2) != Round(resultAngleVW, 2) {
		t.Fatal("angle form by v", v.String(), "and w", w.String(), "should be", expectedAngleVW, "got", resultAngleVW)
	}

	if Round(expectedAngleVS, 2) != Round(resultAngleVS, 2) {
		t.Fatal("angle form by v", v.String(), "and s", s.String(), "should be", expectedAngleVS, "got", resultAngleVS)
	}

	if Round(expectedAngleWU, 2) != Round(resultAngleWU, 2) {
		t.Fatal("angle form by w", w.String(), "and u", u.String(), "should be", expectedAngleWU, "got", resultAngleWU)
	}

	if Round(expectedAngleUS, 2) != Round(resultAngleUS, 2) {
		t.Fatal("angle form by u", u.String(), "and s", s.String(), "should be", expectedAngleUS, "got", resultAngleUS)
	}

	if Round(expectedAngleSW, 2) != Round(resultAngleSW, 2) {
		t.Fatal("angle form by s", s.String(), "and w", w.String(), "should be", expectedAngleSW, "got", resultAngleSW)
	}
}

func TestGetNormalizedAngle(t *testing.T) {
	v := Vector2{
		X:     1,
		Y:     0,
		Angle: 0,
	}
	u := Vector2{
		X:     0,
		Y:     -1,
		Angle: 0,
	}
	expectedAngle := (3*math.Pi)/2
	angle := GetNormalizedAngle(v, u)
	if expectedAngle != angle {
		t.Fatal("angle formed by v", v.String(), "and u", u.String(), "should be", expectedAngle, "got", angle)
	}
}

func TestDoIntersect(t *testing.T) {
	pointA := Vector2{
		X:     1,
		Y:     1,
		Angle: 0,
	}
	pointB := Vector2{
		X:     4,
		Y:     4,
		Angle: 0,
	}
	pointC := Vector2{
		X:     1,
		Y:     3,
		Angle: 0,
	}
	pointD := Vector2{
		X:     5,
		Y:     3,
		Angle: 0,
	}
	pointE := Vector2{
		X:     3,
		Y:     2,
		Angle: 0,
	}
	if !DoIntersect(pointA, pointB, pointC, pointD) {
		t.Fatal("segemnt AB should intersect CD")
	}
	if DoIntersect(pointA, pointD, pointC, pointB) {
		t.Fatal("segment AD should not intersect CB")
	}
	if !DoIntersect(pointA, pointD, pointE, pointB) {
		t.Fatal("segement AD and EB should intersect")
	}
}