package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

type versionInfo struct {
	Release   string
	BuildDate string
	GitHash   string
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of calendar",
	RunE: func(cmd *cobra.Command, args []string) error {
		info := versionInfo{release, buildDate, gitHash}

		if err := json.NewEncoder(os.Stdout).Encode(info); err != nil {
			return fmt.Errorf("error while decode version info: %w", err)
		}

		return nil
	},
}
