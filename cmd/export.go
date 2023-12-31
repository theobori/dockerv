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

	"github.com/spf13/cobra"
	dockerv "github.com/theobori/dockerv/internal/dockerv"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export Docker volumes into a single file.",
	Long: `Export Docker volumes into a single file.
	
Usage examples:

dockerv export \
    --src docker_volume \
    --src path/docker-compose.yml \
    --src path/dir/ \
    --dest volumes.tar.gz \
    --force`,
	Run: func(cmd *cobra.Command, _ []string) {
		dvConfig.PointSource, _ = cmd.Flags().GetStringSlice("src")
		dvConfig.User, _ = cmd.Flags().GetString("user")
		dvConfig.Force, _ = cmd.Flags().GetBool("force")

		dest, _ := cmd.Flags().GetString("dest")

		dvConfig.PointDest = &dest
		dvConfig.Kind = dockerv.Export

		if err := dockerVExecute(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	exportCmd.PersistentFlags().StringSlice(
		"src",
		[]string{},
		"The source point",
	)
	exportCmd.MarkPersistentFlagRequired("src")

	exportCmd.PersistentFlags().String(
		"dest",
		".",
		"The destination point",
	)

	exportCmd.PersistentFlags().BoolP(
		"force",
		"f",
		false,
		"Ignore the Docker volumes that does not exist",
	)

	exportCmd.PersistentFlags().String(
		"user",
		"",
		"User that will be attached to the created container",
	)

	exportCmd.MarkPersistentFlagRequired("dest")

	rootCmd.AddCommand(exportCmd)
}
