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
	"fmt"

	"github.com/docker/docker/client"
	"github.com/theobori/dockerv/internal/point"
)

type ActionKind int

const (
	Import ActionKind = 0
	Export ActionKind = 1
	Move   ActionKind = 3
	List   ActionKind = 4
	Remove ActionKind = 5
	Copy   ActionKind = 6
)

type DockerVConfig struct {
	Kind             ActionKind
	PointSource      []string
	PointDest *string
	Force            bool
	State            bool
}

type ExecutesValueField struct {
	needDestination bool
	callback        func(*[]string) error
}
type Executes map[ActionKind]ExecutesValueField

type DockerV struct {
	config DockerVConfig
	cli    *client.Client

	source      []*point.Point
	destination *point.Point

	executes Executes
}

func NewDockerV(cli *client.Client, config *DockerVConfig) (*DockerV, error) {
	var pointDest *point.Point
	
	if config.PointDest != nil {
		pointDest = point.PointFromValue(
			cli,
			*config.PointDest,
		)
	}
	
	pointsSrc := []*point.Point{}

	for _, pointSrc := range config.PointSource {
		pointsSrc = append(
			pointsSrc,
			point.PointFromValue(cli, pointSrc),
		)
	}

	dv := DockerV{
		*config,
		cli,
		pointsSrc,
		pointDest,
		make(Executes),
	}

	dv.executes[Import] = ExecutesValueField{true, dv._import}
	dv.executes[Export] = ExecutesValueField{true, dv.export}
	dv.executes[Move] = ExecutesValueField{true, dv.move}
	dv.executes[List] = ExecutesValueField{false, dv.list}
	dv.executes[Remove] = ExecutesValueField{false, dv.remove}
	dv.executes[Copy] = ExecutesValueField{true, dv.copy}

	return &dv, nil
}

func NewDefaultDockerV(cli *client.Client) (*DockerV, error) {
	return NewDockerV(
		cli,
		&DockerVConfig{
			Kind: Export,
		},
	)
}

func (dv *DockerV) SetConfig(config *DockerVConfig) {
	dv.config = *config
}

func (dv *DockerV) Config() DockerVConfig {
	return dv.config
}

func (dv *DockerV) Execute() error {
	if len(dv.source) == 0 {
		return fmt.Errorf("invalid source point")
	}

	exec := dv.executes[dv.config.Kind]

	if exec.needDestination && dv.destination == nil {
		return fmt.Errorf("invalid destination point")
	}

	volumesSrc := []string{}

	for _, pSrc := range dv.source {
		vSrc, _ := (*pSrc).Volumes()

		volumesSrc = append(volumesSrc, vSrc...)
	}

	if err := exec.callback(&volumesSrc); err != nil {
		return err
	}

	return nil
}
