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

package dockerv

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
)

type PointKind int
type PointPathKind int
type Volumes []string

const (
	Path         PointKind = 0
	DockerVolume PointKind = 1
	Unknown      PointKind = 2

	DockerCompose PointPathKind = 0
	Directory     PointPathKind = 1
	Tarball       PointPathKind = 2
	UnknownPath   PointPathKind = 3
)

type PointPathKindFunc *func(string) (bool, error)
type PointPathKindFuncs map[PointPathKindFunc]PointPathKind

var pointPathKindFuncs = make(PointPathKindFuncs)

func InitPointKindFuncs() {
	pointPathKindFuncs[&isDirectory] = Directory
	pointPathKindFuncs[&isYAML] = DockerCompose
	pointPathKindFuncs[&isTarball] = Tarball
}

type Point struct {
	value string
	kind  *PointKind

	cli *client.Client
}

func NewPoint(value string) Point {
	return Point{
		value,
		nil,
		nil,
	}
}

func (p *Point) SetClient(cli *client.Client) *Point {
	p.cli = cli

	return p
}

func (p *Point) Value() string {
	return p.value
}

func (p *Point) guessKind() (PointKind, error) {
	if p.cli == nil {
		return Unknown, fmt.Errorf("missing the Docker client instance")
	}

	if DockerVolumeExists(context.Background(), p.cli, p.value) {
		return DockerVolume, nil
	}

	// It doesnt check if the filepath exist for obvious reasons.
	return Path, nil
}

func (p *Point) PathKind() (PointPathKind, error) {
	for f, pathKind := range pointPathKindFuncs {
		is, err := (*f)(p.value)

		if is && err == nil {
			return pathKind, nil
		}
	}

	return UnknownPath, fmt.Errorf("unable to determine path kind")
}

func (p *Point) Kind() PointKind {
	if p.kind != nil {
		return *p.kind
	}

	kind, err := p.guessKind()

	if err != nil {
		fmt.Println(err)
	}

	p.kind = &kind

	return *p.kind
}

func (p *Point) volumesFromPath() Volumes {
	pathKind, err := p.PathKind()

	if err != nil {
		return Volumes{}
	}

	switch pathKind {
	case DockerCompose:
		v, err := volumesFromDockerCompose(p.value)

		if err != nil {
			return Volumes{}
		}

		return v
	case Directory:
		v, err := volumesFromDirectory(p.value)

		if err != nil {
			return Volumes{}
		}

		return v

	case Tarball:
		fmt.Println("Tarball")
	}

	return Volumes{}
}

func (p *Point) Volumes() Volumes {
	if p.cli == nil {
		return Volumes{}
	}

	if p.Kind() == DockerVolume {
		return Volumes{p.value}
	}

	if p.Kind() == Path {
		return p.volumesFromPath()
	}

	return Volumes{}
}
