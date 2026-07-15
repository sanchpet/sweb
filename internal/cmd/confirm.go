package cmd

import (
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// confirmed prompts for a yes/no confirmation, returning true to proceed. It
// short-circuits to true when the command carries a --yes flag set to true; a
// prompt error (e.g. non-interactive) reads as "not confirmed". Commands that
// call it must register a bool --yes flag.
func confirmed(cmd *cobra.Command, title, affirmative string) bool {
	if yes, _ := cmd.Flags().GetBool("yes"); yes {
		return true
	}
	ok := false
	if err := huh.NewConfirm().
		Title(title).
		Affirmative(affirmative).
		Negative("Cancel").
		Value(&ok).
		Run(); err != nil {
		return false
	}
	return ok
}
