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
	"context"
	"fmt"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/theobori/dockerv/internal/docker"
	"github.com/theobori/dockerv/internal/file"
)

type PointKind int
type PointPathKind int

const (
	DockerCompose PointKind = 0
	DockerVolume  PointKind = 1
	Directory     PointKind = 2
	Tarball       PointKind = 3
	Unknown       PointKind = 4
)

type PointPathFunc *func(string) (bool, error)
type PointKindFuncs map[PointPathFunc]PointKind

var pointKindFuncs = make(PointKindFuncs)

func InitPointKindFuncs() {
	pointKindFuncs[&file.IsDirectory] = Directory
	pointKindFuncs[&file.IsDockerCompose] = DockerCompose
	pointKindFuncs[&file.IsTarball] = Tarball
	pointKindFuncs[&file.HasTarballExt] = Tarball
}

type PointMetadata struct {
	value string
	kind  *PointKind
}

func NewPointMetadata(value string) PointMetadata {
	return PointMetadata{
		value,
		nil,
	}
}

func (p *PointMetadata) Value() string {
	return p.value
}

func (p *PointMetadata) AbsValue() string {
	abs, err := filepath.Abs(p.value)

	if err != nil {
		return p.value
	}

	return abs
}

func (p *PointMetadata) pathKind() (PointKind, error) {
	for f, pathKind := range pointKindFuncs {
		is, err := (*f)(p.value)

		if is && err == nil {
			return pathKind, nil
		}
	}

	return Unknown, fmt.Errorf("unable to determine path kind")
}

func (p *PointMetadata) guessKind(cli *client.Client) (PointKind, error) {
	if docker.DockerVolumeExists(context.Background(), cli, p.value) {
		return DockerVolume, nil
	}

	return p.pathKind()
}

func (p *PointMetadata) Kind(cli *client.Client) PointKind {
	if p.kind != nil {
		return *p.kind
	}

	kind, _ := p.guessKind(cli)

	p.kind = &kind

	return *p.kind
}
