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
	"fmt"

	"github.com/docker/docker/client"
	"github.com/theobori/dockerv/common"
	"github.com/theobori/dockerv/internal/file"
)

type DockerComposePoint struct {
	metadata *PointMetadata
	cli      *client.Client
}

var NewDockerComposePoint = func(cli *client.Client, metadata *PointMetadata) Point {
	return &DockerComposePoint{
		metadata,
		cli,
	}
}

func (d *DockerComposePoint) Metadata() *PointMetadata {
	return d.metadata
}

func (d *DockerComposePoint) resolveVolumeNames(volumes []string) []string {
	dir, err := common.PreviousDirName(d.metadata.value)

	if err != nil {
		return volumes
	}

	volumePrefix := dir + "_"

	for i, volume := range volumes {
		volumes[i] = volumePrefix + volume
	}

	return volumes
}

func (d *DockerComposePoint) Volumes() ([]string, error) {
	yamlData, err := file.ParseYAML(d.metadata.value)

	if err != nil {
		return []string{}, err
	}

	volumesField, ok := yamlData["volumes"].(map[string]any)

	if !ok {
		return []string{}, fmt.Errorf("missing the volumes field")
	}

	volumesKeys := common.MapKeys(volumesField)

	return d.resolveVolumeNames(volumesKeys), nil
}

func (d *DockerComposePoint) From(*[]string) error {
	return ErrOperation
}

func (d *DockerComposePoint) To(*[]string) error {
	return ErrOperation
}
