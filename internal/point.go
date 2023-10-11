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
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
)

type PointKind int
type PointPathKind int

const (
	Path         PointKind = 0
	DockerVolume PointKind = 1
	Unknown      PointKind = 1

	DockerCompose PointPathKind = 0
	Directory     PointPathKind = 1
	Tarball       PointPathKind = 2
	UnknownPath   PointPathKind = 3
)

type PointPathKindFunc *func(*os.File) (bool, error)
type PointPathKindFuncs map[PointPathKindFunc]PointPathKind

var pointPathKindFuncs = make(PointPathKindFuncs)

type Point struct {
	value string
	kind  *PointKind

	cli *client.Client
}

var isDirectory = func (file *os.File) (bool, error) {
	fileInfo, err := file.Stat()

	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

var isDockerCompose = func (file *os.File) (bool, error) {
	// TODO: use compose-go lib to parse as YAML file
	return true, nil
}

// Just checking file ext
var isTarball = func (file *os.File) (bool, error) {
	return strings.HasSuffix(file.Name(), ".tar.gz"), nil
}

func InitPointKindFuncs() {
	pointPathKindFuncs[&isDirectory] = Directory
	pointPathKindFuncs[&isDockerCompose] = Directory
	pointPathKindFuncs[&isTarball] = Directory
}

func NewPoint(value string) Point {
	return Point{
		value,
		nil,
		nil,
	}
}

func (p *Point) WithClient(cli *client.Client) *Point {
	p.cli = cli

	return p
}

func (p *Point) Value() string {
	return p.value
}

func (p *Point) Equal(pDest *Point) bool {
	kind := p.Kind()

	if kind != pDest.Kind() {
		return false
	}

	if kind == DockerVolume {
		return p.Value() == pDest.Value()
	}

	return filepath.Clean(p.Value()) == filepath.Clean(pDest.Value())
}

func (p *Point) guessKind() (PointKind, error) {
	if p.cli == nil {
		return Unknown, fmt.Errorf("missing the Docker client instance")
	}

	if IsDockerVolume(context.Background(), p.cli, p.value) {
		return DockerVolume, nil
	}

	// It doesnt check if the filepath exist for obvious reasons.
	return Path, nil
}

func (p *Point) PathKind() (PointPathKind, error) {
	var err error

	file, err := os.Open(p.value)

	if err != nil {
		return UnknownPath, err
	}

	defer file.Close()

	for f, pathKind := range pointPathKindFuncs {
		is, err := (*f)(file)

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

	*p.kind = kind

	return *p.kind
}
