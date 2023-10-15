package point

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
)

type DirectoryPoint struct {
	value string
	cli   *client.Client
}

var NewDirectoryPoint = func(cli *client.Client, value string) Point {
	return &DirectoryPoint{
		value,
		cli,
	}
}

func (d *DirectoryPoint) Kind() PointKind {
	return Directory
}

func (d *DirectoryPoint) Volumes() ([]string, error) {
	volumes := []string{}

	err := filepath.Walk(d.value,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			p := PointFromValue(d.cli, path)

			if p == nil {
				return nil
			}

			if (*p).Kind() == Directory {
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

func (d *DirectoryPoint) Import([]string) error {
	return fmt.Errorf("unallowed point for this operation")
}

func (d *DirectoryPoint) Export(p *Point) error {
	return nil
}
