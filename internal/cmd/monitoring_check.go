package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// checkCmd groups the monitoring-check operations (endpoint /monitoring/checks):
// listing checks, the reference dictionaries, the create/edit lifecycle, the
// activate/deactivate toggle, removal, and check history. It hangs off the
// monitoring group.
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Manage monitoring checks",
}

// checkID parses a check's <id> positional argument.
func checkID(arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err != nil {
		return 0, fmt.Errorf("check id must be an integer: %q", arg)
	}
	return id, nil
}

func init() {
	monitoringCmd.AddCommand(checkCmd)
}
