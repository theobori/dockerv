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
)

type ActionKind int

const (
	Import ActionKind = 0
	Export ActionKind = 1
	Copy   ActionKind = 2
	Move   ActionKind = 3
	List   ActionKind = 4
	Remove ActionKind = 5
)

type DockerVConfig struct {
	Kind             ActionKind
	PointSource      Point
	PointDestination Point
}

type Executes map[ActionKind]func() error

type DockerV struct {
	config DockerVConfig
	cli    *client.Client

	executes Executes
}

func NewDockerV(config *DockerVConfig) (*DockerV, error) {
	executes := make(Executes)

	dv := DockerV{
		*config,
		nil,
		executes,
	}

	dv.executes[Import] = dv._import
	dv.executes[Export] = dv.export
	dv.executes[Copy] = dv.copy
	dv.executes[Move] = dv.move
	dv.executes[List] = dv.list
	dv.executes[Remove] = dv.remove

	return &dv, nil
}

func NewDefaultDockerV() (*DockerV, error) {
	return NewDockerV(
		&DockerVConfig{
			Kind: Export,
		},
	)
}

func (dv *DockerV) SetDockerClient(cli *client.Client) {
	dv.cli = cli
}

func (dv *DockerV) SetConfig(config *DockerVConfig) {
	dv.config = *config
}

func (dv *DockerV) Config() DockerVConfig {
	return dv.config
}

func (dv *DockerV) _import() error {
	return nil
}

func (dv *DockerV) export() error {
	return nil
}

func (dv *DockerV) list() error {
	for _, volume := range dv.config.PointSource.Volumes() {
		fmt.Println(volume)
	}

	return nil
}

func (dv *DockerV) move() error {
	return nil
}

func (dv *DockerV) copy() error {
	return nil
}

func (dv *DockerV) remove() error {
	return nil
}

func (dv *DockerV) Execute() error {
	return dv.executes[dv.config.Kind]()
}
