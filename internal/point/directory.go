package point

import (
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

func (dp *DirectoryPoint) Kind() PointKind {
	return Directory
}

func (dp *DirectoryPoint) Volumes() ([]string, error) {
	volumes := []string{}

	err := filepath.Walk(dp.value,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			p := PointFromValue(dp.cli, path)

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

func (dp *DirectoryPoint) Import([]string) error {
	return nil
}

func (dp *DirectoryPoint) Export(p *Point) error {
	return nil
}
