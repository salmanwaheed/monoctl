package cmd

import (
	"fmt"

	"github.com/salmanwaheed/monoctl"
	"github.com/salmanwaheed/monoctl/internal/constants"
	"github.com/spf13/cobra"
)

var verCmd = &cobra.Command{
  Use: "version",
  Short: fmt.Sprintf("version for %s", rootCmd.Root().Use),
  Args: cobra.ExactArgs(0),
  RunE: func(cmd *cobra.Command, args []string) error {
    data := monoctl.BuildInfo{
      Version: "dev",
      Commit: "none",
      Repository: constants.Repository,
      Maintainer: constants.Maintainer,
    }

    return monoctl.Formatter.Execute(cmd, data)
  },
}

func init() {
  // list flags
  monoctl.Formatter.BindFlag(verCmd)

  rootCmd.AddCommand(verCmd)
}
