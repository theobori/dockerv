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

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use: "import",
	Short: "Import Docker volumes.",
	Long: `Import Docker volumes.
	
Usage examples:

dockerv import \
    --src volumes.tar.gz \
    --dest path/ \
    --force

dockerv import \
    --src volume.tar.gz \
    --dest new_volume_name \
    --create

dockerv import --src volume.tar.gz`,
	Run: func(cmd *cobra.Command, _ []string) {
		dvConfig.PointSource, _ = cmd.Flags().GetStringSlice("src")
		dvConfig.Force, _ = cmd.Flags().GetBool("force")

		dest, _ := cmd.Flags().GetString("dest")

		dvConfig.PointDest = &dest
		dvConfig.Kind = dockerv.Import

		if err := dockerVExecute(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	importCmd.PersistentFlags().String(
		"dest",
		".",
		"The destination point",
	)

	importCmd.PersistentFlags().BoolP(
		"force",
		"f",
		false,
		"Ignore the destination point volumes not fully matching the source point volumes",
	)

	importCmd.MarkPersistentFlagRequired("dest")

	rootCmd.AddCommand(importCmd)
}
