package container

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
)

// DockerClient hold client connection in order to pull image and start container
type DockerClient struct {
	Cli *client.Client
}

// PodFactory for Docker manager

// DockerContainerFactory implement PodFactory interface for Docker manager
type DockerContainerFactory struct {
	Client *DockerClient
}

// NewDockerContainerFactory create a DockerContainerFactory
func NewDockerContainerFactory() (*DockerContainerFactory, error) {
	var err error
	dockerContainerFactory := new(DockerContainerFactory)
	dockerContainerFactory.Client = new(DockerClient)
	dockerContainerFactory.Client.Cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return dockerContainerFactory, err
}

// PullPod pull image if it does not exist locally
func (dContainerFactory *DockerContainerFactory) PullPod(info *PodInfo) error {
	ctx := context.Background()
	imageName := createImageNameWithVersion(info)
	imageList, err := dContainerFactory.Client.Cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return err
	}
	for _, image := range imageList {
		if image.RepoTags[0] == imageName {
			return nil
		}
	}
	// Weird bug : when output is not catch, image is not pull properly
	reader, err := dContainerFactory.Client.Cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if reader == nil {
		return err
	}
	io.Copy(ioutil.Discard, reader)
	return err
}

// PodExecutor for Docker manager

// DockerContainerExecutor implement PodExecutor interface for Docker manager
type DockerContainerExecutor struct {
	Client *DockerClient
}

// NewDockerContainerExecutor create a DockerContainerExecutor
func NewDockerContainerExecutor() (*DockerContainerExecutor, error) {
	var err error
	dockerContainerExecutor := new(DockerContainerExecutor)
	dockerContainerExecutor.Client = new(DockerClient)
	dockerContainerExecutor.Client.Cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return dockerContainerExecutor, err
}

// Run start a Docker container
func (dContainerExecutor *DockerContainerExecutor) Run(info *PodInfo) error {
	ctx := context.Background()
	cmd := dContainerExecutor.prepareArgs(info)
	imageName := createImageNameWithVersion(info)

	containerID, err := dContainerExecutor.createContainer(ctx, &container.Config{
		Hostname: info.Hostname,
		Env:      info.Env,
		Cmd:      cmd,
		Image:    imageName,
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{},
		AutoRemove:    true,
	})
	if err != nil {
		return err
	}
	err = dContainerExecutor.attachNetworks(ctx, info, containerID)
	if err != nil {
		return err
	}
	return dContainerExecutor.Client.Cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}

func (dContainerExecutor *DockerContainerExecutor) attachNetworks(ctx context.Context, info *PodInfo, containerID string) error {
	var networkID []string
	nets, err := dContainerExecutor.Client.Cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return err
	}
	for _, net := range nets {
		toLowerCaseNet := strings.ToLower(net.Name)
		for _, netToAttach := range info.Networks {
			toLowerCaseNetToAttach := strings.ToLower(netToAttach)
			if strings.Contains(toLowerCaseNet, toLowerCaseNetToAttach) {
				networkID = append(networkID, net.ID)
			}
		}
	}
	if len(networkID) != len(info.Networks) {
		return errors.New("unable to find asked networks")
	}
	for _, netID := range networkID {
		err = dContainerExecutor.Client.Cli.NetworkConnect(ctx, netID, containerID, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dContainerExecutor *DockerContainerExecutor) prepareArgs(info *PodInfo) strslice.StrSlice {
	var cmd strslice.StrSlice
	for _, arg := range info.Args {
		cmd = append(cmd, arg)
	}
	return cmd
}

func (dContainerExecutor *DockerContainerExecutor) createContainer(ctx context.Context, config *container.Config, hostConfig *container.HostConfig) (string, error) {
	resp, err := dContainerExecutor.Client.Cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func createImageNameWithVersion(info *PodInfo) string {
	if info.Version == "" {
		return info.ImageName + ":latest"
	}
	return info.ImageName + ":" + info.Version
}
