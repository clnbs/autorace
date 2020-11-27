package container

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockDockerFactory struct{}

func (mockDockerFactory *MockDockerFactory) PullPod(info *PodInfo) error {
	return nil
}

type MockDockerExecutor struct{}

func (mockDockerExecutor *MockDockerExecutor) Run(info *PodInfo) error {
	return nil
}

func TestNewDockerPodManager(t *testing.T) {
	containerInfo, err := NewPodInfoFromYAMLFile(os.Getenv("HOME") + "/go/src/github.com/clnbs/autorace/test/podManager/testDockerPodManager.yaml")
	assert.Nil(t, err, "while ready pod configuration from YAML file : ", err)

	dockerFactory, err := NewDockerContainerFactory()
	assert.Nil(t, err, "while creating docker factory : ", err)

	dockerExecutor, err := NewDockerContainerExecutor()
	assert.Nil(t, err, "while creating docker executor : ", err)

	dockerManager := NewPodManager(containerInfo, dockerFactory, dockerExecutor)
	err = dockerManager.PullPod()
	assert.Nil(t, err, "while pulling pod : ", err)

	err = dockerManager.Run()
	assert.Nil(t, err, "while running pod : ", err)
}

func TestMockDockerPodManagerWithoutAction(t *testing.T) {
	mockContainerInfo := &PodInfo{
		ImageName: "mock",
		Version:   "latest",
		Hostname:  "mock",
		Env: []string{
			"MOCK_VALUE=mock",
			"ENVIRONMENT=" + os.Getenv("ENVIRONMENT"),
		},
		Args:     nil,
		Networks: []string{"rabbitmq", "logs", "autorace_cache"},
	}
	mockDockerFactory := &MockDockerFactory{}
	mockDockerExecutor := &MockDockerExecutor{}
	dockerContainer := NewPodManager(mockContainerInfo, mockDockerFactory, mockDockerExecutor)

	err := dockerContainer.PullPod()
	assert.Nil(t, err, "error while pulling image (this should not happen) :", err)

	err = dockerContainer.Run()
	assert.Nil(t, err, "error while starting container (this should not happen) ", err)
}
