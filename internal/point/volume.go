package point

import (
	"github.com/docker/docker/client"
)

type DockerVolumePoint struct {
	value string
	cli   *client.Client
}

var NewDockerVolumePoint = func(cli *client.Client, value string) Point {
	return &DockerVolumePoint{
		value,
		cli,
	}
}

func (d *DockerVolumePoint) Kind() PointKind {
	return DockerVolume
}

func (d *DockerVolumePoint) Volumes() ([]string, error) {
	return []string{d.value}, nil
}

func (d *DockerVolumePoint) Import([]string) error {
	return nil
}

func (d *DockerVolumePoint) Export(p *Point) error {
	return nil
}
