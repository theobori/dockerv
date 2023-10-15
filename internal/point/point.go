package point

import (
	"fmt"
	"reflect"

	"github.com/docker/docker/client"
)

var ErrOperation = fmt.Errorf("unallowed point for this operation")

type Point interface {
	Volumes() ([]string, error)
	From([]string) error
	To([]string) error
	Metadata() *PointMetadata
}

type pointResolverFunc = *func(*client.Client, *PointMetadata) Point

var pointResolver = make(map[PointKind]pointResolverFunc)

func InitPointResolvers() {
	pointResolver[DockerCompose] = &NewDockerComposePoint
	pointResolver[Directory] = &NewDirectoryPoint
	pointResolver[DockerVolume] = &NewDockerVolumePoint
	pointResolver[Tarball] = &NewTarballPoint
}

func Init() {
	InitPointKindFuncs()
	InitPointResolvers()
}

func PointFromMetadata(cli *client.Client, metadata *PointMetadata) *Point {
	for k, v := range pointResolver {
		if metadata.Kind(cli) != k {
			continue
		}

		if k == DockerVolume {
			metadata.value = metadata.Value()
		} else {
			metadata.value = metadata.AbsValue()
		}

		p := (*v)(cli, metadata)

		return &p
	}

	return nil
}

func PointFromValue(cli *client.Client, value string) *Point {
	metadata := NewPointMetadata(value)

	return PointFromMetadata(cli, &metadata)
}

func PointEqual(pSrc *Point, pDest *Point) bool {
	vSrc, _ := (*pSrc).Volumes()
	vDest, _ := (*pDest).Volumes()

	return reflect.DeepEqual(vSrc, vDest)
}
