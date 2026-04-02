package main

import (
	"context"
	"log"
	"net/netip"
	"path/filepath"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

func checkOrCreateContainer(ctx context.Context, apiClient *client.Client, containerName string, imageName string) string {

	f := client.Filters{}
	f.Add("name", containerName)

	l, err := apiClient.ContainerList(ctx, client.ContainerListOptions{
		All:     true,
		Filters: f,
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(l.Items) != 0 {
		return l.Items[0].ID
	}

	hostPath, _ := filepath.Abs("./containers/" + containerName)
	hostIP := netip.MustParseAddr("0.0.0.0")
	port := network.MustParsePort("25565/tcp")

	portBindings := network.PortMap{
		port: []network.PortBinding{
			{
				HostIP:   hostIP,
				HostPort: "25565",
			},
		},
	}

	resp, err := apiClient.ContainerCreate(ctx, client.ContainerCreateOptions{
		Name:  containerName,
		Image: imageName,
		Config: &container.Config{
			Cmd:        []string{"java", "-jar", "server.jar"},
			WorkingDir: "/container",
			ExposedPorts: network.PortSet{
				port: struct{}{},
			},
		},
		HostConfig: &container.HostConfig{
			Binds: []string{
				hostPath + ":/container",
			},
			PortBindings: portBindings,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	return resp.ID
}
