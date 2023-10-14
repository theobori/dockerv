package point

import (
	"github.com/docker/docker/client"
)

type Point interface {
	Kind() PointKind
	Volumes() ([]string, error)
	Import([]string) error
	Export(*Point) error
}

type pointResolverFunc = *func(*client.Client, string) Point

var pointResolver = make(map[PointKind]pointResolverFunc)

func InitPointResolvers() {
	pointResolver[DockerCompose] = &NewDockerComposePoint
	pointResolver[Directory] = &NewDirectoryPoint
}

func Init() {
	InitPointKindFuncs()
	InitPointResolvers()
}

func PointFromMetadata(cli *client.Client, metadata *PointMetadata) *Point {
	var value string

	for k, v := range pointResolver {
		if metadata.Kind(cli) != k {
			continue
		}

		if k == DockerVolume {
			value = metadata.Value()
		} else {
			value = metadata.AbsValue()
		}

		p := (*v)(cli, value)

		return &p
	}

	return nil
}

func PointFromValue(cli *client.Client, value string) *Point {
	metadata := NewPointMetadata(value)

	return PointFromMetadata(cli, &metadata)
}
