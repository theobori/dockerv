package point

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/theobori/dockerv/common"
)

type TarballPoint struct {
	value string
	cli   *client.Client
}

var NewTarballPoint = func(cli *client.Client, value string) Point {
	return &TarballPoint{
		value,
		cli,
	}
}

func (t *TarballPoint) Kind() PointKind {
	return Tarball
}

func (t *TarballPoint) Volumes() ([]string, error) {
	volumes := []string{}

	return volumes, nil
}

func (t *TarballPoint) execContainer(
	ctx context.Context,
	command *strslice.StrSlice,
	mounts *[]mount.Mount,
) error {
	resp, err := t.cli.ContainerCreate(
		ctx,
		&container.Config{
			Cmd: *command,
			Image: "busybox",
		},
		&container.HostConfig{
			Mounts: *mounts,
		}, nil, nil, "",
	)

	if err != nil {
		return err
	}

	if err = t.containerStartAutoRemove(ctx, resp.ID); err != nil {
		return err
	}

	return nil
}

func (t *TarballPoint) containerStartAutoRemove(
	ctx context.Context,
	id string,
) error {
	var err error

	if err = t.cli.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
		return err
	}

	if err = t.cli.ContainerStop(ctx, id, container.StopOptions{}); err != nil {
		return err
	}

	if err = t.cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{}); err != nil {
		return err
	}

	return nil
}

func (t *TarballPoint) tarballFromVolume(
	ctx context.Context,
	volume string,
	tmpVolumeName string,
) error {
	return t.execContainer(
		ctx,
		&strslice.StrSlice{
			"tar",
			"-cvzf", "/dest/" + volume + ".tar.gz",
			"-C", "/src", ".",
		},
		&[]mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: volume,
				Target: "/src",
			},
			{
				Type:   mount.TypeVolume,
				Source: tmpVolumeName,
				Target: "/dest",
			},
		},
	)
}

func (t *TarballPoint) packTarballs(ctx context.Context, tmpVolumeName string) error {
	return t.execContainer(
		ctx,
		&strslice.StrSlice{
			"tar",
			"-cvzf", "/dest/" + filepath.Base(t.value),
			"-C", "/src", ".",
		},
		&[]mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: tmpVolumeName,
				Target: "/src",
			},
			{
				Type:   mount.TypeBind,
				Source: common.PopFilename(t.value),
				Target: "/dest",
			},
		},
	)
}

func (t *TarballPoint) Import(volumes []string) error {
	u, err := uuid.NewRandom()
	ctx := context.Background()

	if err != nil {
		return err
	}

	tmpVolumeName := strings.Replace(u.String(), "-", "", -1)

	// Create a tmp Docker volume to store very tarball
	if _, err = t.cli.VolumeCreate(
		ctx,
		volume.CreateOptions{Name: tmpVolumeName},
	); err != nil {
		return err
	}

	// Create tarballs
	for _, volume := range volumes {
		if err = t.tarballFromVolume(ctx, volume, tmpVolumeName); err != nil {
			return err
		}
	}

	// Pack every tarball together
	if err = t.packTarballs(ctx, tmpVolumeName); err != nil {
		return err
	}

	// Remove the Docker volume
	if err = t.cli.VolumeRemove(ctx, tmpVolumeName, true); err != nil {
		return err
	}

	return nil
}

func (t *TarballPoint) Export(p *Point) error {
	return nil
}
