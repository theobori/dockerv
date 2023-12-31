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

package common

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/theobori/dockerv/internal/file"
	"golang.org/x/exp/slices"
)

func MapKeys(m map[string]any) []string {
	keys := make([]string, len(m))
	i := 0

	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

func FilenameWithoutExt(path string) string {
	return strings.Split(filepath.Base(path), ".")[0]
}

func PreviousDirName(path string) (string, error) {
	if isDir, _ := file.IsDirectory(path); isDir {
		return path, nil
	}

	splittedPath := strings.Split(path, "/")

	if len(splittedPath) < 2 {
		return "", fmt.Errorf("unable to resolve the previous dir")
	}

	return splittedPath[len(splittedPath)-2], nil
}

func PopFilename(path string) string {
	return strings.Replace(
		path,
		filepath.Base(path),
		"",
		-1,
	)
}

func IsSliceinSlice[T comparable](src []T, dest []T) bool {
	for _, e := range dest {
		if !slices.Contains(src, e) {
			return false
		}
	}

	return true
}
