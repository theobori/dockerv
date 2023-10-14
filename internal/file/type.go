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

package file

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/compose-spec/compose-go/loader"
	"github.com/theobori/dockerv/internal/docker"
)

var IsDirectory = func(path string) (bool, error) {
	file, err := os.Open(path)

	if err != nil {
		return false, err
	}

	fileInfo, err := file.Stat()

	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

func ParseYAML(path string) (map[string]any, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	buffer, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	yamlData, err := loader.ParseYAML(buffer)

	if err != nil {
		return nil, err
	}

	return yamlData, nil
}

var isYAML = func(path string) (bool, error) {
	_, err := ParseYAML(path)

	if err != nil {
		return false, err
	}

	return true, nil
}

var IsDockerCompose = func(path string) (bool, error) {
	isYaml, err := isYAML(path)

	if err != nil {
		return false, err
	}

	nameFull := strings.Split(filepath.Base(path), ".")
	name := nameFull[0]
	ext := nameFull[len(nameFull) - 1]

	return name == docker.DockerComposeFilename &&
		isYaml &&
		(ext == "yml" || ext == "yaml"),
		nil
}

var IsTarball = func(path string) (bool, error) {
	file, err := os.Open(path)

	if err != nil {
		return false, err
	}

	gr, err := gzip.NewReader(file)

	if err != nil {
		return false, err
	}

	tr := tar.NewReader(gr)

	if _, err := tr.Next(); err != nil {
		return false, err
	}

	return true, nil
}
