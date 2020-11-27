package environment

import (
	"os"
	"strings"
)

type Environment int

// possible running environment
const (
	// DEV  typically dev's computer
	DEV Environment = iota
	// LOCAL dev's computer with production environment set up locally
	LOCAL
	// PREPRODUCTION
	PREPRODUCTION
	// PRODUCTION
	PRODUCTION
)

var stringifyEnvironment map[string]Environment

var environmentList = []string{"dev", "local", "preproduction", "production"}

var currentEnvironment Environment

func init() {
	initEnvironment()
}

func initEnvironment() {
	stringifyEnvironment = make(map[string]Environment, len(environmentList))
	for index, env := range environmentList {
		stringifyEnvironment[env] = intToEnvironments(index)
	}
	currentEnvironment = stringToEnvironment(os.Getenv("ENVIRONMENT"))
}

func GetCurrentEnvironment() Environment {
	return currentEnvironment
}

func (env Environment) String() string {
	if env > PRODUCTION || env < DEV {
		return environmentList[0]
	}
	return environmentList[env]
}

func stringToEnvironment(strEnv string) Environment {
	strEnv = strings.ToLower(strEnv)
	if _, ok := stringifyEnvironment[strEnv]; !ok {
		return LOCAL
	}
	return stringifyEnvironment[strEnv]
}

func intToEnvironments(number int) Environment {
	if number >= len(environmentList) || number < 0 {
		return LOCAL
	}
	return Environment(number)
}
