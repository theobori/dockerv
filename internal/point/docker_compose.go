package point

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/theobori/dockerv/common"
	"github.com/theobori/dockerv/internal/file"
)

type DockerComposePoint struct {
	metadata *PointMetadata
	cli   *client.Client
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

func (d *DockerComposePoint) From([]string) error {
	return ErrOperation
}

func (d *DockerComposePoint) To([]string) error {
	return ErrOperation
}
