package cmd

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

// TestHostingCronTree asserts the cron subgroup hangs off hosting with the full
// list/add/edit/remove lifecycle.
func TestHostingCronTree(t *testing.T) {
	var hosting *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "hosting" {
			hosting = c
		}
	}
	if hosting == nil {
		t.Fatal("hosting command not registered")
	}

	var cronGroup *cobra.Command
	for _, c := range hosting.Commands() {
		if c.Name() == "cron" {
			cronGroup = c
		}
	}
	if cronGroup == nil {
		t.Fatal("hosting cron command not registered")
	}

	sub := subNames(cronGroup)
	for _, n := range []string{"list", "add", "edit", "remove"} {
		if !sub[n] {
			t.Errorf("hosting cron is missing subcommand %q", n)
		}
	}
}

// TestScheduleFromFlags checks that the shared schedule builder rejects a
// missing --command and otherwise reflects the flag values.
func TestScheduleFromFlags(t *testing.T) {
	// Missing --command is an error.
	cmd := &cobra.Command{}
	addScheduleFlags(cmd)
	if _, err := scheduleFromFlags(cmd); !errors.Is(err, errMissingCommand) {
		t.Fatalf("scheduleFromFlags without --command: err = %v, want errMissingCommand", err)
	}

	// A populated set of flags maps onto the schedule.
	cmd = &cobra.Command{}
	addScheduleFlags(cmd)
	if err := cmd.Flags().Set("minute", "30"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("hour", "12"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("command", "backup.sh"); err != nil {
		t.Fatal(err)
	}
	sc, err := scheduleFromFlags(cmd)
	if err != nil {
		t.Fatalf("scheduleFromFlags: unexpected error %v", err)
	}
	if sc.Minute != 30 || sc.Hour != 12 || sc.Command != "backup.sh" {
		t.Errorf("schedule = %+v, want minute=30 hour=12 command=backup.sh", sc)
	}
}
