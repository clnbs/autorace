package mathtool

import (
	"math"
	"math/rand"
	"time"
)

//RandomIntBetween returns a random int between two values (int)
func RandomIntBetween(intOne, intTwo int) int {
	if intOne == intTwo {
		return intOne
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	if intOne <= 0 && intTwo <= 0 {
		intOne *= -1
		intTwo *= -1
		if intTwo < intOne {
			return (r1.Intn(intOne-intTwo) + intTwo) * -1
		}
		return (r1.Intn(intTwo-intOne) + intOne) * - 1
	}

	if intTwo < intOne {
		return r1.Intn(intOne-intTwo) + intTwo
	}
	return r1.Intn(intTwo-intOne) + intOne
}

//RandomFloatBetween returns a random float64 between two values (float64)
func RandomFloatBetween(floatOne, floatTwo float64) float64 {
	if floatOne == floatTwo {
		return floatOne
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	if floatTwo < floatOne {
		return floatTwo + r1.Float64()*(floatOne-floatTwo)
	}

	return floatOne + r1.Float64()*(floatTwo-floatOne)
}

//ClampFloat64 return x value between a and b value
func ClampFloat64(x, a, b float64) float64 {
	if a > b {
		a, b = b, a
	}
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

//ClampInt return x value between a and b value
func ClampInt(x, a, b int) int {
	if a > b {
		a, b = b, a
	}
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

//IsFloat64Between check if x value is between a and b value
func IsFloat64Between(x, a, b float64) bool {
	if (x < a && x > b) || (x > a && x < b) {
		return true
	}
	return false
}

//IsIntBetween check if x value is between a and b value
func IsIntBetween(x, a, b int) bool {
	if (x < a && x > b) || (x > a && x < b) {
		return true
	}
	return false
}

//Round give a float with a precision
func Round(x float64, prec int) float64 {
	output := math.Pow(10, float64(prec))
	return float64(math.Round(x*output)) / output
}
