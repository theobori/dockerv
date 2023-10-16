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

	"github.com/theobori/dockerv/common"
	"github.com/theobori/dockerv/internal/docker"
)

func (dv *DockerV) export(vSrc *[]string) error {
	vSrcExist := docker.DockerGetVolumesExist(
		context.Background(),
		dv.cli,
		*vSrc,
	)

	if !dv.config.Force && !common.IsSliceinSlice(vSrcExist, *vSrc){
		return fmt.Errorf("missing volumes, use --force or -f to skip")
	}

	if err := (*dv.destination).From(&vSrcExist); err != nil {
		return err
	}

	return nil
}
