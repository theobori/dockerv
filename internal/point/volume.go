package point

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/theobori/dockerv/internal/docker"
)

type DockerVolumePoint struct {
	metadata *PointMetadata
	cli   *client.Client
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

func (d *DockerVolumePoint) headVolume(volumes []string) string {
	if len(volumes) < 1 {
		return ""
	}

	return volumes[0]
}

func (d *DockerVolumePoint) From(vSrc []string) error {
	return d.copy(d.headVolume(vSrc), d.metadata.value)
}

func (d *DockerVolumePoint) To(vDest []string) error {
	return d.copy(d.metadata.value, d.headVolume(vDest))
}
