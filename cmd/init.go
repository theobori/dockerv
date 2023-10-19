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
	"github.com/theobori/dockerv/internal/docker"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init the needed resources needed for dockerv.",
	Long: `Init the needed resources needed for dockerv.

It can be Docker images or anything related.	`,
	Run: func(cmd *cobra.Command, args []string) {
		img := docker.DockerImage

		pulled, err := docker.DockerTryPull(
			cmd.Context(),
			cli,
			img,
		)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if pulled {
			fmt.Println("Pulled", img)
		} else {
			fmt.Println("You already have", img)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
