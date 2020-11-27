package environment

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetCurrentEnvironment_local(t *testing.T) {
	err := os.Setenv("ENVIRONMENT", "local")
	initEnvironment()
	assert.Nil(t, err, "while setting up environment variable :", err)
	expectedResult := LOCAL
	result := GetCurrentEnvironment()
	assert.Equal(t, expectedResult, result, "environment should be set to", expectedResult.String(), ", got", result.String())
}

func TestGetCurrentEnvironment_local2(t *testing.T) {
	err := os.Setenv("ENVIRONMENT", "LoCaL")
	initEnvironment()
	assert.Nil(t, err, "while setting up environment variable :", err)
	expectedResult := LOCAL
	result := GetCurrentEnvironment()
	assert.Equal(t, expectedResult, result, "environment should be set to", expectedResult.String(), ", got", result.String())
}

func TestGetCurrentEnvironment_dev(t *testing.T) {
	err := os.Setenv("ENVIRONMENT", "dev")
	initEnvironment()
	assert.Nil(t, err, "while setting up environment variable :", err)
	expectedResult := DEV
	result := GetCurrentEnvironment()
	assert.Equal(t, expectedResult, result, "environment should be set to", expectedResult.String(), ", got", result.String())
}

func TestGetCurrentEnvironment_preproduction(t *testing.T) {
	err := os.Setenv("ENVIRONMENT", "preproduction")
	initEnvironment()
	assert.Nil(t, err, "while setting up environment variable :", err)
	expectedResult := PREPRODUCTION
	result := GetCurrentEnvironment()
	assert.Equal(t, expectedResult, result, "environment should be set to", expectedResult.String(), ", got", result.String())
}

func TestGetCurrentEnvironment_production(t *testing.T) {
	err := os.Setenv("ENVIRONMENT", "production")
	initEnvironment()
	assert.Nil(t, err, "while setting up environment variable :", err)
	expectedResult := PRODUCTION
	result := GetCurrentEnvironment()
	assert.Equal(t, expectedResult, result, "environment should be set to", expectedResult.String(), ", got", result.String())
}
