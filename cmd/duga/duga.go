package duga

import "github.com/spf13/cobra"

// NewCmd returns the `duga` subcommand group.
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "duga",
		Short: "DUGA Affiliate API commands",
	}

	cmd.AddCommand(newSearchCmd())

	return cmd
}
