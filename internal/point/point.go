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
	"reflect"

	"github.com/docker/docker/client"
)

var ErrOperation = fmt.Errorf("this operation with this point is not authorized")

type Point interface {
	Volumes() ([]string, error)
	From(*[]string) error
	To(*[]string) error
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

func PointFromValue(cli *client.Client, value string, user string) *Point {
	metadata := NewPointMetadata(value, user)

	return PointFromMetadata(cli, &metadata)
}

func PointEqual(pSrc *Point, pDest *Point) bool {
	vSrc, _ := (*pSrc).Volumes()
	vDest, _ := (*pDest).Volumes()

	return reflect.DeepEqual(vSrc, vDest)
}
