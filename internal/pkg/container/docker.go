package container

import (
	"context"
	"errors"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
)

// CreateDynamicServer create a dynamic server with a PartyUUID in argument.
// CreateDynamicServer bind the dynamic server instance to networks in order to make it works
func CreateDynamicServer(partyID string, env []string) error {
	ctx := context.Background()
	var networkID []string
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	imageName := "autorace_dynamic:latest"
	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return err
	}
	for _, net := range networks {
		if strings.Contains(strings.ToLower(net.Name), "rabbitmq") ||
			strings.Contains(strings.ToLower(net.Name), "logs") ||
			strings.Contains(strings.ToLower(net.Name), "autorace_cache") {
			networkID = append(networkID, net.ID)
		}
	}
	if len(networkID) == 0 {
		return errors.New("unable to find networks")
	}
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Hostname:   "dynamic_" + partyID,
		Domainname: "",
		User:       "",
		Cmd: strslice.StrSlice{
			partyID,
		},
		ArgsEscaped: false,
		Image:       imageName,
		Entrypoint:  nil,
		Env:         env,
	}, &container.HostConfig{
		Binds:           nil,
		ContainerIDFile: "",
		NetworkMode:     "",
		RestartPolicy:   container.RestartPolicy{},
		AutoRemove:      true,
	},
		nil,
		nil,
		"dynamic_"+partyID)
	if err != nil {
		return err
	}
	for _, netID := range networkID {
		err = cli.NetworkConnect(ctx, netID, resp.ID, nil)
		if err != nil {
			return err
		}
	}
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}
