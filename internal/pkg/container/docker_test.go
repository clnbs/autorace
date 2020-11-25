package container

import (
	"os"
	"strings"
	"testing"
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
	environment := os.Getenv("ENVIRONMENT")
	environment = strings.ToLower(environment)
	switch environment {
	case "preproduction":
		testDockerPodManager(t)
	default:
		testMockDockerPodManagerWithoutAction(t)
	}
}

func testDockerPodManager(t *testing.T) {
	containerInfo := &PodInfo{
		ImageName: "hello-world",
		Version:   "",
		Hostname:  "testing_hello-world",
		Env: []string{
			"ENVIRONMENT=" + os.Getenv("ENVIRONMENT"),
		},
		Args:     nil,
		Networks: nil,
	}
	dockerFactory, err := NewDockerContainerFactory()
	if err != nil {
		t.Fatal("while creating docker factory :", err)
	}
	dockerExecutor, err := NewDockerContainerExecutor()
	if err != nil {
		t.Fatal("while creating docker executor :", err)
	}
	dockerManager := NewPodManager(containerInfo, dockerFactory, dockerExecutor)
	err = dockerManager.PullPod()
	if err != nil {
		t.Fatal("while pulling pod :", err)
	}
	err = dockerManager.Run()
	if err != nil {
		t.Fatal("while running pod :", err)
	}
}

func testMockDockerPodManagerWithoutAction(t *testing.T) {
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
	if err != nil {
		t.Fatal("error while pulling image (this should not happen), err :", err)
	}

	err = dockerContainer.Run()
	if err != nil {
		t.Fatal("error while starting container (this should not happen), err :", err)
	}
}
