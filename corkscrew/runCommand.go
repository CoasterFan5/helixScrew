package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/moby/moby/client"
)

func runCommand(apiClient *client.Client, ctx context.Context, containerId string, cmd []string) client.ExecAttachResult {
	r, err := apiClient.ExecCreate(ctx, containerId, client.ExecCreateOptions{
		Cmd:          cmd,
		WorkingDir:   "/container",
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	execAttachResp, err := apiClient.ExecAttach(ctx, r.ID, client.ExecAttachOptions{})
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(os.Stdout, execAttachResp.Reader)

	return execAttachResp
}
