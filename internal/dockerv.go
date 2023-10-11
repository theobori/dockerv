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

import "github.com/docker/docker/client"

type ActionKind int

const (
	Import ActionKind = 0
	Export ActionKind = 1
	Copy   ActionKind = 2
	Move   ActionKind = 3
)

type DockerVConfig struct {
	Kind             ActionKind
	PointSource      Point
	PointDestination Point
	Recursive        bool
}

type DockerV struct {
	config DockerVConfig
	cli    *client.Client
}

func NewDockerV(config *DockerVConfig) (*DockerV, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return nil, err
	}

	return &DockerV{
		*config,
		cli,
	}, nil
}

func NewDefaultDockerV() (*DockerV, error) {
	return NewDockerV(
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
