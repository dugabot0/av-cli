package dmm

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/av-cli/cmd"
	"github.com/yourusername/av-cli/internal/client"
	"github.com/yourusername/av-cli/internal/output"
)

func newMakerCmd() *cobra.Command {
	var (
		floorID string
		initial string
		hits    int
		offset  int
	)

	c := &cobra.Command{
		Use:   "maker",
		Short: "Search DMM makers",
		RunE: func(c *cobra.Command, args []string) error {
			if floorID == "" {
				cmd.Logf("--floor-id is required (get it from `av-cli dmm floors`)")
				os.Exit(cmd.ExitInput)
			}
			cl := cmd.LoadDMMClient()
			result, err := cl.MakerSearch(client.MakerSearchParams{
				FloorID: floorID,
				Initial: initial,
				Hits:    hits,
				Offset:  offset,
			})
			if err != nil {
				cmd.HandleError(err)
			}
			output.Print(result.Result, cmd.Pretty())
			return nil
		},
	}

	f := c.Flags()
	f.StringVar(&floorID, "floor-id", "", "Floor ID (required, from `dmm floors`)")
	f.StringVar(&initial, "initial", "", "Kana prefix filter")
	f.IntVar(&hits, "hits", 100, "Results per page (1-500)")
	f.IntVar(&offset, "offset", 1, "Start position")

	return c
}
