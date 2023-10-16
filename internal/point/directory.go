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
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
)

type DirectoryPoint struct {
	metadata *PointMetadata
	cli      *client.Client
}

var NewDirectoryPoint = func(cli *client.Client, metadata *PointMetadata) Point {
	return &DirectoryPoint{
		metadata,
		cli,
	}
}

func (d *DirectoryPoint) Metadata() *PointMetadata {
	return d.metadata
}

func (d *DirectoryPoint) Volumes() ([]string, error) {
	volumes := []string{}

	err := filepath.Walk(d.metadata.value,
		func(path string, _ os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			p := PointFromValue(d.cli, path)

			if p == nil {
				return nil
			}

			if (*p).Metadata().Kind(d.cli) == Directory {
				return nil
			}

			pVolumes, err := (*p).Volumes()

			if err != nil {
				return nil
			}

			volumes = append(volumes, pVolumes...)

			return nil
		},
	)

	if err != nil {
		return []string{}, err
	}

	return volumes, nil
}

func (d *DirectoryPoint) From(*[]string) error {
	return ErrOperation
}

func (d *DirectoryPoint) To(*[]string) error {
	return ErrOperation
}
