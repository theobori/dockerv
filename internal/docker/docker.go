/*
Copyright © 2023 Théo Bori <nagi@tilde.team>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
)

// Without ext
const DockerComposeFilename = "docker-compose"

// Image used to operate with volumes
const DockerImage = "busybox"

func DockerVolumeExists(ctx context.Context, cli *client.Client, name string) bool {
	_, err := cli.VolumeInspect(ctx, name)

	return err == nil
}

func containerStartAutoRemove(ctx context.Context, cli *client.Client, id string) error {
	var err error

	if err = cli.ContainerStart(
		ctx,
		id,
		types.ContainerStartOptions{},
	); err != nil {
		return err
	}

	if err = cli.ContainerStop(
		ctx,
		id,
		container.StopOptions{},
	); err != nil {
		return err
	}

	if err = cli.ContainerRemove(
		ctx,
		id,
		types.ContainerRemoveOptions{},
	); err != nil {
		return err
	}

	return nil
}

func ExecContainer(ctx context.Context, cli *client.Client, command *strslice.StrSlice, mounts *[]mount.Mount) error {
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Cmd:   *command,
			Image: DockerImage,
		},
		&container.HostConfig{
			Mounts: *mounts,
		}, nil, nil, "",
	)

	if err != nil {
		return err
	}

	if err = containerStartAutoRemove(ctx, cli, resp.ID); err != nil {
		return err
	}

	return nil
}

func DockerVolumeCopy(ctx context.Context, cli *client.Client, vSrc string, vDest string) error {
	return ExecContainer(
		ctx,
		cli,
		&strslice.StrSlice{
			"cp", "-rT", "/src", "/dest",
		},
		&[]mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: vSrc,
				Target: "/src",
			},
			{
				Type:   mount.TypeVolume,
				Source: vDest,
				Target: "/dest",
			},
		},
	)
}
