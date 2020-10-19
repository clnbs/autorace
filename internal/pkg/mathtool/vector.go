package mathtool

import (
	"fmt"
	"math"
)

// Orientation is used to represent orientation of two segment in an Enum style
type Orientation int

// possible orentation
const (
	// COLINEAR is when two segment are colinear
	COLINEAR Orientation = iota
	// CLOCKWISE is when orientation of two segments are clockwise
	CLOCKWISE
	// COUNTER_CLOCKWISE is when orientation of two segments are counter clockwise
	COUNTER_CLOCKWISE
)

//Vector2 represent a vector in 2D
type Vector2 struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Angle float64
}

// NewVector create a null vector
func NewVector() Vector2 {
	return Vector2{
		X:     0,
		Y:     0,
		Angle: 0,
	}
}

//String return an stringify Vector2
func (v Vector2) String() string {
	return fmt.Sprintf("X : %.2f\tY : %.2f", v.X, v.Y)
}

//Cross 2D cross product of OA and OB vectors, i.e. z-component of their 3D cross product.
// Returns a positive value, if OAB makes a counter-clockwise turn,
// negative for clockwise turn, and zero if the points are collinear.
func Cross(o, a, b Vector2) float64 {
	return (a.X-o.X)*(b.Y-o.Y) - (a.Y-o.Y)*(b.X-o.X)
}

//Distance return the distance between two vector
func Distance(a, b Vector2) float64 {
	x := a.X - b.X
	y := a.Y - b.Y
	return math.Sqrt(x*x + y*y)
}

//Set is used to set X and Y of a 2d vector
func (v Vector2) Set(x, y float64) Vector2 {
	v.X = x
	v.Y = y
	return v
}

//Equals return true if 2 vector are equals and false otherwise
func (v Vector2) Equals(w Vector2) bool {
	return math.Round(v.X) == math.Round(w.X) && math.Round(v.Y) == math.Round(w.Y)
}

//Add sum up two vectors
func (v Vector2) Add(w Vector2) Vector2 {
	v.X += w.X
	v.Y += w.Y
	return v
}

// Subtract subtract vector w to vector v
func (v Vector2) Subtract(w Vector2) Vector2 {
	v.X -= w.X
	v.Y -= w.Y
	return v
}

//Divide divide X and Y of a vector by a given float
func (v Vector2) Divide(s float64) Vector2 {
	if s == 0 {
		return v
	}
	v.X = v.X / float64(s)
	v.Y = v.Y / float64(s)
	return v
}

//Rotate is used to rotate a fictional angle of a point (which does not mean anything
//mathematically speaking).
func (v Vector2) Rotate(angle float64) Vector2 {
	v.Angle += angle
	return v
}

//ScaleWithAngle is used to scale a vector length in the direction of the angle sets
func (v Vector2) ScaleWithAngle(distance float64) Vector2 {
	v.X += math.Cos(v.Angle) * distance
	v.Y += math.Sin(v.Angle) * distance
	return v
}

//GetCentroid return the centroid of a vector slice
func GetCentroid(points []Vector2) Vector2 {
	var centroid Vector2
	for _, point := range points {
		centroid.X += point.X
		centroid.Y += point.Y
	}
	centroid.X /= float64(len(points))
	centroid.Y /= float64(len(points))
	return centroid
}

//GetNormalizedDirection return a normalized vector between two points
func GetNormalizedDirection(a, b Vector2) Vector2 {
	return Vector2{
		X:     (b.X - a.X) / Distance(a, b),
		Y:     (b.Y - a.Y) / Distance(a, b),
		Angle: 0,
	}
}

//GetVectorAngle return the angle form by the X axis and the vector in radian
func (v Vector2) GetVectorAngle() float64 {
	return math.Atan(v.Y / v.X)
}

//Normalized is used to normalized a vector length
func (v Vector2) Normalized() Vector2 {
	length := v.Length()
	return Vector2{
		X:     v.X / length,
		Y:     v.Y / length,
		Angle: 0,
	}
}

//Length return the length of a vector
func (v Vector2) Length() float64 {
	return math.Round(math.Sqrt(v.X*v.X + v.Y*v.Y))
}

//SetLength is used to set the length of a vector
func (v Vector2) SetLength(l float64) Vector2 {
	return Vector2{
		X:     v.X * l,
		Y:     v.Y * l,
		Angle: 0,
	}
}

//Move is use to move a point by a given vector
func (v Vector2) Move(w Vector2) Vector2 {
	return Vector2{
		X:     v.X + w.X,
		Y:     v.Y + w.Y,
		Angle: 0,
	}
}

//Dot return the dot product of two vectors
func Dot(v, w Vector2) float64 {
	return (v.X * w.X) + (v.Y * w.Y)
}

//GetVectorsAngle return the angle form by two vector
func GetVectorsAngle(v, w Vector2) float64 {
	//arccos[(xa * xb + ya * yb) / (√(xa2 + ya2) * √(xb2 + yb2))]
	return math.Acos((v.X*w.X + v.Y*w.Y) / (math.Sqrt(v.X*v.X+v.Y*v.Y) * math.Sqrt(w.X*w.X+w.Y*w.Y)))
}

// GetNormalizedAngle form by two vector
func GetNormalizedAngle(v, w Vector2) float64 {
	return normalizedAngle(w) - normalizedAngle(v)
}

func normalizedAngle(v Vector2) float64 {
	angle := math.Atan2(v.Y, v.X)
	if angle < 0 {
		return angle + 2*math.Pi
	}
	return angle
}

//check is point x is on the segment ab
func onSegment(x, a, b Vector2) bool {
	if (x.X <= math.Max(a.X, b.X)) && x.X >= math.Min(a.X, b.X) && x.Y <= math.Max(a.Y, b.Y) && x.Y >= math.Min(a.Y, b.Y) {
		return true
	}
	return false
}

//find orientation of ordered triplet vector
func orientation(a, b, c Vector2) Orientation {
	value := (b.Y-a.Y) * (c.X - b.X) - (b.X - a.X) * (c.Y - b.Y)
	if value == 0 {
		return COLINEAR
	}
	if value > 0 {
		return CLOCKWISE
	}
	return COUNTER_CLOCKWISE
}

//DoIntersect return true if segment AB intersect CD
func DoIntersect(a, b, c, d Vector2) bool {
	orientation1 := orientation(a, b, c)
	orientation2 := orientation(a, b, d)
	orientation3 := orientation(c, d, a)
	orientation4 := orientation(c, d, b)
	//General case
	if orientation1 != orientation2 && orientation3 != orientation4 {
		return true
	}
	//a, b, c are colinear (c on segment ab)
	if orientation1 == COLINEAR && onSegment(a, c, b) {
		return true
	}
	// a, d, b are colinear (d on segment ab)
	if orientation2 == COLINEAR && onSegment(a, d, b) {
		return true
	}
	// c, d, a are colinear (a on sedment cd)
	if orientation3 == COLINEAR && onSegment(c, a, d) {
		return true
	}
	// c, d, b are colinear (b on segment cd)
	if orientation4 == COLINEAR && onSegment(c, b, d) {
		return true
	}
	return false
}
