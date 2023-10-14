package point

import (
	"fmt"
	"strings"

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

func (dc *DockerComposePoint) Kind() PointKind {
	return DockerCompose
}

func (dc *DockerComposePoint) resolveVolumeNames(volumes []string) []string {
	splittedPath := strings.Split(dc.value, "/")

	if len(splittedPath) < 2 {
		return volumes
	}

	volumePrefix := splittedPath[len(splittedPath)-2] + "_"

	for i, volume := range volumes {
		volumes[i] = volumePrefix + volume
	}

	return volumes
}

func (dc *DockerComposePoint) Volumes() ([]string, error) {
	yamlData, err := file.ParseYAML(dc.value)

	if err != nil {
		return []string{}, err
	}

	volumesField, ok := yamlData["volumes"].(map[string]any)

	if !ok {
		return []string{}, fmt.Errorf("missing the volumes field")
	}

	volumesKeys := common.MapKeys(volumesField)

	return dc.resolveVolumeNames(volumesKeys), nil
}

func (dc *DockerComposePoint) Import([]string) error {
	return nil
}

func (dc *DockerComposePoint) Export(p *Point) error {
	return nil
}
