package container

type PodFactory interface {
	PullPod(info *PodInfo) error
}

type PodExecutor interface {
	Run(info *PodInfo) error
}

type PodInfo struct {
	ImageName string
	Version   string
	Hostname  string
	Env       []string
	Args      []string
	Networks  []string
}

func NewPodInfo(imageName, version string, env, args, networks []string) *PodInfo {
	return &PodInfo{ImageName: imageName, Version: version, Env: env, Args: args, Networks: networks}
}

type PodManager struct {
	Info     *PodInfo
	Factory  PodFactory
	Executor PodExecutor
}

func NewPodManager(info *PodInfo, factory PodFactory, executor PodExecutor) *PodManager {
	return &PodManager{Info: info, Factory: factory, Executor: executor}
}

func (podManager *PodManager) Run() error {
	return podManager.Executor.Run(podManager.Info)
}

func (podManager *PodManager) PullPod() error {
	return podManager.Factory.PullPod(podManager.Info)
}
