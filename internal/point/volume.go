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

	"github.com/docker/docker/client"
	"github.com/theobori/dockerv/internal/docker"
)

type DockerVolumePoint struct {
	metadata *PointMetadata
	cli      *client.Client
}

var NewDockerVolumePoint = func(cli *client.Client, metadata *PointMetadata) Point {
	return &DockerVolumePoint{
		metadata,
		cli,
	}
}

func (d *DockerVolumePoint) Metadata() *PointMetadata {
	return d.metadata
}

func (d *DockerVolumePoint) Volumes() ([]string, error) {
	return []string{d.metadata.value}, nil
}

func (d *DockerVolumePoint) copy(vSrc string, vDest string) error {
	return docker.DockerVolumeCopy(
		context.Background(),
		d.cli,
		vSrc,
		vDest,
	)
}

func (d *DockerVolumePoint) headVolume(volumes *[]string) string {
	if len(*volumes) < 1 {
		return ""
	}

	return (*volumes)[0]
}

func (d *DockerVolumePoint) From(vSrc *[]string) error {
	return d.copy(d.headVolume(vSrc), d.metadata.value)
}

func (d *DockerVolumePoint) To(vDest *[]string) error {
	return ErrOperation
}
