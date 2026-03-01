package dmm

import "github.com/spf13/cobra"

// NewCmd returns the `dmm` subcommand group.
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dmm",
		Short: "DMM Affiliate API v3 commands",
	}

	cmd.AddCommand(
		newFloorsCmd(),
		newItemsCmd(),
		newActressCmd(),
		newGenreCmd(),
		newMakerCmd(),
		newSeriesCmd(),
		newAuthorCmd(),
	)

	return cmd
}
