package dmm

import (
	"github.com/spf13/cobra"
	"github.com/yourusername/av-cli/cmd"
	"github.com/yourusername/av-cli/internal/output"
)

func newFloorsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "floors",
		Short: "List all DMM sites → services → floors",
		Long:  "Returns the full hierarchy of sites, services, and floors. Floor IDs are required for other dmm subcommands.",
		RunE: func(c *cobra.Command, args []string) error {
			client := cmd.LoadDMMClient()
			result, err := client.FloorList()
			if err != nil {
				cmd.HandleError(err)
			}
			output.Print(result.Result, cmd.Pretty())
			return nil
		},
	}
}
