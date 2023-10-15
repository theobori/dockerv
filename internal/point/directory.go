package point

import (
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
)

type DirectoryPoint struct {
	metadata *PointMetadata
	cli   *client.Client
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
		func(path string, info os.FileInfo, err error) error {
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

func (d *DirectoryPoint) From([]string) error {
	return ErrOperation
}

func (d *DirectoryPoint) To([]string) error {
	return ErrOperation
}
