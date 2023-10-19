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
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/volume"
	"github.com/spf13/cobra"
	"github.com/theobori/dockerv/internal/docker"
	dockerv "github.com/theobori/dockerv/internal/dockerv"
)

// moveCmd represents the export command
var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move a Docker volume to another.",
	Long: `Move a Docker volume to another.
	
The source volume is deleted after being copied.`,
	Run: func(cmd *cobra.Command, _ []string) {
		dvConfig.PointSource, _ = cmd.Flags().GetStringSlice("src")
		dest, _ := cmd.Flags().GetString("dest")
		force, _ := cmd.Flags().GetBool("force")

		dvConfig.PointDest = &dest
		dvConfig.Kind = dockerv.Move

		// Creates the destination if not exists
		ctx := context.Background()

		if force && !docker.DockerVolumeExists(ctx, cli, dest) {
			cli.VolumeCreate(
				ctx,
				volume.CreateOptions{Name: dest},
			)
		}

		if err := dockerVExecute(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	moveCmd.PersistentFlags().StringSlice(
		"src",
		[]string{},
		"The source point",
	)
	moveCmd.MarkPersistentFlagRequired("src")

	moveCmd.PersistentFlags().String(
		"dest",
		".",
		"The destination Docker volume",
	)

	moveCmd.PersistentFlags().BoolP(
		"force",
		"f",
		false,
		"Creates the Docker volume if it does not exist",
	)

	moveCmd.MarkPersistentFlagRequired("dest")

	rootCmd.AddCommand(moveCmd)
}
