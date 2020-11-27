package container

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

// PodFactory interface is used to bake application pod in order to be usable
type PodFactory interface {
	PullPod(info *PodInfo) error
}

// PodExecutor interface is used to start application pod once initialized by PodFactory
type PodExecutor interface {
	Run(info *PodInfo) error
}

// PodInfo hold Pod configuration
type PodInfo struct {
	ImageName string   `yaml:"image_name"`
	Version   string   `yaml:"version"`
	Hostname  string   `yaml:"hostname"`
	Env       []string `yaml:"env"`
	Args      []string `yaml:"args"`
	Networks  []string `yaml:"networks"`
}

// PodInformationFromYaml is used to extract Pod Information from a YAML file under the filed "pod_information"
type PodInformationFromYaml struct {
	PodInformation PodInfo `yaml:"pod_information"`
}

// NewPodInfo create a PodInfo instance from arguments passed to this method
func NewPodInfo(imageName, version, hostname string, env, args, networks []string) *PodInfo {
	return &PodInfo{ImageName: imageName, Version: version, Hostname: hostname, Env: env, Args: args, Networks: networks}
}

// NewPodInfoFromEnvironment creation a PodInfo instance from environment variable
func NewPodInfoFromEnvironment() (*PodInfo, error) {
	info := new(PodInfo)
	info.ImageName = os.Getenv("POD_IMAGE_NAME")
	info.Version = os.Getenv("POD_VERSION")
	info.Hostname = os.Getenv("POD_HOSTNAME")
	info.Env = getPodEnvironmentList("POD_ENV", ";")
	info.Args = getPodEnvironmentList("POD_ARGS", ";")
	info.Networks = getPodEnvironmentList("POD_NETWORKS", ";")
	if !verifyPodInformations(info) {
		return nil, errors.New("pod information are incomplete")
	}
	return info, nil
}

// NewPodInfoFromYAMLFile create a PodInfo instance from a YAML file
func NewPodInfoFromYAMLFile(pathToFile string) (*PodInfo, error) {
	infoFromYaml := new(PodInformationFromYaml)
	yamlFile, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, infoFromYaml)
	if err != nil {
		return nil, err
	}
	if !verifyPodInformations(&infoFromYaml.PodInformation) {
		return nil, errors.New("pod information are incomplete")
	}
	for index, env := range infoFromYaml.PodInformation.Env {
		if strings.Contains(env, "ENVIRONMENT") {
			infoFromYaml.PodInformation.Env[index] = "ENVIRONMENT=" + os.Getenv("ENVIRONMENT")
			break
		}
	}
	return &infoFromYaml.PodInformation, nil
}

// PodManager is a wrapper for managing pod within a go program.
// It contains pod's information, a pod factory and a pod executor
type PodManager struct {
	Info     *PodInfo
	Factory  PodFactory
	Executor PodExecutor
}

// Create a PodManager from a PodInfo, a PodFactory and a PodExecutor
func NewPodManager(info *PodInfo, factory PodFactory, executor PodExecutor) *PodManager {
	return &PodManager{Info: info, Factory: factory, Executor: executor}
}

// Run wrap up PodExecutor -- Not really SOLID friendly but make PodManager easier to use
func (podManager *PodManager) Run() error {
	return podManager.Executor.Run(podManager.Info)
}

// PullPod wrap up PodFactory -- Not really SOLID friendly but make PodManager easier to use
func (podManager *PodManager) PullPod() error {
	return podManager.Factory.PullPod(podManager.Info)
}

func getPodEnvironmentList(environmentVariableName, separator string) []string {
	unformattedEnvironmentList := os.Getenv(environmentVariableName)
	environmentList := strings.Split(unformattedEnvironmentList, separator)
	return environmentList
}

func verifyPodInformations(info *PodInfo) bool {
	if info.Hostname == "" || info.ImageName == "" {
		return false
	}
	return true
}
