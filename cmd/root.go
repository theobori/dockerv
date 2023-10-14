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

package cmd

import (
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	dockerv "github.com/theobori/dockerv/internal"
	"github.com/theobori/dockerv/internal/point"
)

var (
	dvConfig = dockerv.DockerVConfig{}

	cli *client.Client

	rootCmd = &cobra.Command{
		Use:   "dockerv",
		Short: "Manage your Docker volumes",
		Long: `dockerv is a tool to import, export,
move and copy your Docker volumes.

It also supports docker compose
files and has recursive features.

A point can be:
	- A tarball
	- A Docker volume
	- A Docker compose file
	- A directory
`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String(
		"src",
		".",
		"The source point",
	)
	rootCmd.MarkPersistentFlagRequired("src")

	// DockerV setup
	point.Init()

	var err error

	// Setup the Docker Client
	cli, err = client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
