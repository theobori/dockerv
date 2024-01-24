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
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
)

// Without ext
const DockerComposeFilename = "docker-compose"

// Image used to operate with volumes
const DockerImage = "busybox:1.36"

func DockerVolumeExists(ctx context.Context, cli *client.Client, name string) bool {
	_, err := cli.VolumeInspect(ctx, name)

	return err == nil
}

func DockerGetVolumesExist(ctx context.Context, cli *client.Client, names []string) []string {
	ret := []string{}

	for _, name := range names {
		if DockerVolumeExists(ctx, cli, name) {
			ret = append(ret, name)
		}
	}

	return ret
}

func DockerRunSelfRemove(ctx context.Context, cli *client.Client, id string) error {
	if err := cli.ContainerStart(
		ctx,
		id,
		types.ContainerStartOptions{},
	); err != nil {
		return err
	}

	statusCh, errCh := cli.ContainerWait(
		context.Background(),
		id,
		container.WaitConditionNotRunning,
	)

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}

	case <-statusCh:
	}

	if err := cli.ContainerStop(
		ctx,
		id,
		container.StopOptions{},
	); err != nil {
		return err
	}

	if err := cli.ContainerRemove(
		ctx,
		id,
		types.ContainerRemoveOptions{},
	); err != nil {
		return err
	}

	return nil
}

func ExecContainer(
	ctx context.Context,
	cli *client.Client,
	command *strslice.StrSlice,
	mounts *[]mount.Mount,
	user string,
) error {
	config := container.Config{
		Cmd:   *command,
		Image: DockerImage,
	}

	if user != "" {
		config.User = user
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&config,
		&container.HostConfig{
			Mounts: *mounts,
		}, nil, nil, "",
	)

	if err != nil {
		return err
	}

	if err := DockerRunSelfRemove(ctx, cli, resp.ID); err != nil {
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
		}, "",
	)
}

func DockerImageExists(ctx context.Context, cli *client.Client, id string) bool {
	if _, _, err := cli.ImageInspectWithRaw(ctx, id); err != nil {
		return false
	}

	return true
}

func DockerTryPull(ctx context.Context, cli *client.Client, id string) (bool, error) {
	if DockerImageExists(ctx, cli, id) {
		return false, nil
	}

	out, err := cli.ImagePull(
		ctx,
		"docker.io/library/"+DockerImage,
		types.ImagePullOptions{},
	)

	if err != nil {
		return false, err
	}

	defer out.Close()

	if _, err := io.ReadAll(out); err != nil {
		return false, err
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
