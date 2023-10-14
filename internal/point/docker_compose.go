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

func (dc *DockerComposePoint) Kind() PointKind {
	return DockerCompose
}

func (dv *DockerComposePoint) resolveVolumeNames(volumes []string) []string {
	return []string{}
}

func (dc *DockerComposePoint) Volumes() ([]string, error) {
	yamlData, err := file.ParseYAML(dc.value)

	if err != nil {
		return []string{}, err
	}

	volumes, ok := yamlData["volumes"].(map[string]any)

	if !ok {
		return []string{}, fmt.Errorf("missing the volumes field")
	}

	return common.MapKeys(volumes), nil
}

func (dc *DockerComposePoint) Import([]string) error {
	return nil
}

func (dc *DockerComposePoint) Export(p *Point) error {
	return nil
}
