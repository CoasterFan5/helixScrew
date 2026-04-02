package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/moby/moby/client"
)

func main() {

	containerName := "abc123"
	imageName := "ghcr.io/pterodactyl/yolks:java_25"

	os.MkdirAll("./containers/"+containerName, 0o777)

	// pull the latest jar image
	createJar("./containers/" + containerName + "/server.jar")
	// start container
	ctx := context.Background()
	apiClient, err := client.New(client.FromEnv, client.WithUserAgent("my-application/1.0.0"))
	if err != nil {
		log.Fatal(err)
	}
	defer apiClient.Close()

	out, err := apiClient.ImagePull(ctx, imageName, client.ImagePullOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	io.Copy(os.Stdout, out)

	containerId := checkOrCreateContainer(ctx, apiClient, containerName, imageName)

	_, err = apiClient.ContainerStart(ctx, containerId, client.ContainerStartOptions{})
	if err != nil {
		log.Fatal(err)
	}

	inspectCtr, err := apiClient.ContainerInspect(ctx, containerId, client.ContainerInspectOptions{})
	if err != nil {
		log.Fatal(err)
	}
	if !inspectCtr.Container.State.Running {
		log.Fatal("Container not running")
	}

}
