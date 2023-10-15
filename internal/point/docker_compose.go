package point

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/theobori/dockerv/common"
	"github.com/theobori/dockerv/internal/file"
)

type DockerComposePoint struct {
	value string
	cli   *client.Client
}

var NewDockerComposePoint = func(cli *client.Client, value string) Point {
	return &DockerComposePoint{
		value,
		cli,
	}
}

func (d *DockerComposePoint) Kind() PointKind {
	return DockerCompose
}

func (d *DockerComposePoint) resolveVolumeNames(volumes []string) []string {
	dir, err := common.PreviousDirName(d.value)

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
	yamlData, err := file.ParseYAML(d.value)

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

func (d *DockerComposePoint) Import([]string) error {
	return nil
}

func (d *DockerComposePoint) Export(p *Point) error {
	return nil
}
