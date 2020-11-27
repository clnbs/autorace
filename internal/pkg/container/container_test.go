package container

import (
	"github.com/clnbs/autorace/internal/pkg/environment"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var podInformationToCompareTo *PodInfo

func setup(t *testing.T) {
	podInformationToCompareTo = &PodInfo{
		ImageName: "dummyImage",
		Version:   "dummy-1.0.0",
		Hostname:  "dummy",
		Env: []string{
			"ENVIRONMENT=" + environment.GetCurrentEnvironment().String(),
			"LOG_LEVEL=debug",
		},
		Args: []string{
			"ls",
			"-lah",
		},
		Networks: []string{
			"dummy_net",
			"logs",
		},
	}
}

func TestNewPodInfo(t *testing.T) {
	setup(t)
	var imageName, version, hostname string
	var env, args, nets []string

	imageName = "dummyImage"
	version = "dummy-1.0.0"
	hostname = "dummy"
	env = []string{
		"ENVIRONMENT=" + environment.GetCurrentEnvironment().String(),
		"LOG_LEVEL=debug",
	}
	args = []string{
		"ls",
		"-lah",
	}
	nets = []string{
		"dummy_net",
		"logs",
	}
	info := NewPodInfo(imageName, version, hostname, env, args, nets)
	assert.Equal(t, podInformationToCompareTo, info, "information content should be equals with function NewPodInfo")
}

func TestNewPodInfoFromEnvironment(t *testing.T) {
	setup(t)
	info, err := NewPodInfoFromEnvironment()
	assert.Nil(t, err, "error while getting new pod info from environment variables :", err)
	assert.Equal(t, podInformationToCompareTo, info, "information content should be equals with TestNewPodInfoFromEnvironment")
}

func TestNewPodInfoFromYAMLFile(t *testing.T) {
	setup(t)
	configPath := os.Getenv("HOME") + "/go/src/github.com/clnbs/autorace/test/podManager/testPodInfo_dev.yaml"
	info, err := NewPodInfoFromYAMLFile(configPath)
	assert.Nil(t, err, "error while getting new pod info from YAML file :", err)
	assert.Equal(t, podInformationToCompareTo, info, "information content should be equals with NewPodInfoFromYAMLFile")
}
