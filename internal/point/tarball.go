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

package point

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/theobori/dockerv/common"
	"github.com/theobori/dockerv/internal/docker"
	"github.com/theobori/dockerv/internal/file"
	"golang.org/x/exp/slices"
)

type TarballPoint struct {
	metadata *PointMetadata
	cli      *client.Client
}

var NewTarballPoint = func(cli *client.Client, metadata *PointMetadata) Point {
	return &TarballPoint{
		metadata,
		cli,
	}
}

func (t *TarballPoint) Metadata() *PointMetadata {
	return t.metadata
}

func (t *TarballPoint) Volumes() ([]string, error) {
	volumes := []string{}

	tr, err := file.TarGzReader(t.metadata.value)

	if err != nil {
		return []string{}, err
	}

	for {
		header, err := tr.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return []string{}, err
		}

		volumeName := common.FilenameWithoutExt(header.Name)

		if volumeName == "" {
			continue
		}

		if !strings.HasSuffix(header.Name, ".tar.gz") {
			return []string{}, fmt.Errorf("invalid file format")
		}

		volumes = append(volumes, volumeName)
	}

	return volumes, nil
}

func (t *TarballPoint) From(vSrc *[]string) error {
	u, err := uuid.NewRandom()
	ctx := context.Background()

	if err != nil {
		return err
	}

	tmpVolumeName := strings.Replace(u.String(), "-", "", -1)

	// Create a tmp Docker volume to store every tarball
	if _, err = t.cli.VolumeCreate(
		ctx,
		volume.CreateOptions{Name: tmpVolumeName},
	); err != nil {
		return err
	}

	// Create tarballs
	for _, volume := range *vSrc {
		err = docker.ExecContainer(
			ctx,
			t.cli,
			&strslice.StrSlice{
				"tar",
				"-cvzf", "/dest/" + volume + ".tar.gz",
				"-C", "/src", ".",
			},
			&[]mount.Mount{
				{
					Consistency: mount.ConsistencyFull,
					Type:        mount.TypeVolume,
					Source:      volume,
					Target:      "/src",
				},
				{
					Consistency: mount.ConsistencyFull,
					Type:        mount.TypeVolume,
					Source:      tmpVolumeName,
					Target:      "/dest",
				},
			},
		)

		if err != nil {
			return err
		}

		fmt.Println("Exported", volume+".tar.gz")
	}

	filenameBase := filepath.Base(t.metadata.value)

	// Pack every tarball together
	err = docker.ExecContainer(
		ctx,
		t.cli,
		&strslice.StrSlice{
			"tar",
			"-cvzf", "/dest/" + filenameBase,
			"-C", "/src", ".",
		},
		&[]mount.Mount{
			{
				Consistency: mount.ConsistencyFull,
				Type:        mount.TypeVolume,
				Source:      tmpVolumeName,
				Target:      "/src",
			},
			{
				Consistency: mount.ConsistencyFull,
				Type:        mount.TypeBind,
				Source:      common.PopFilename(t.metadata.value),
				Target:      "/dest",
			},
		},
	)

	if err != nil {
		return err
	}

	// Remove the Docker volume
	if err = t.cli.VolumeRemove(ctx, tmpVolumeName, true); err != nil {
		return err
	}

	return nil
}

func (t *TarballPoint) To(vDest *[]string) error {
	ctx := context.Background()
	vSrc, _ := t.Volumes()

	filenameBase := filepath.Base(t.metadata.value)

	for _, volume := range *vDest {
		if !slices.Contains(vSrc, volume) {
			fmt.Println("Skip", volume)
			continue
		}

		err := docker.ExecContainer(
			ctx,
			t.cli,
			&strslice.StrSlice{
				"/bin/sh", "-c",
				"tar -xzvf " +
					"/" + filenameBase + " ./" + volume +
					".tar.gz " + "-C / && tar -xzvf " +
					"/" + volume + ".tar.gz" + " -C /dest",
			},
			&[]mount.Mount{
				{
					Consistency: mount.ConsistencyFull,
					Type:        mount.TypeBind,
					Source:      t.metadata.value,
					Target:      "/" + filenameBase,
				},
				{
					Consistency: mount.ConsistencyFull,
					Type:        mount.TypeVolume,
					Source:      volume,
					Target:      "/dest",
				},
			},
		)

		if err != nil {
			return err
		}

		fmt.Println("Imported", volume)
	}

	return nil
}
