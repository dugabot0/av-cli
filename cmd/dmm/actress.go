package dmm

import (
	"github.com/spf13/cobra"
	"github.com/yourusername/av-cli/cmd"
	"github.com/yourusername/av-cli/internal/client"
	"github.com/yourusername/av-cli/internal/output"
)

func newActressCmd() *cobra.Command {
	var (
		keyword    string
		actressID  string
		initial    string
		gteBust    string
		lteBust    string
		gteWaist   string
		lteWaist   string
		gteHip     string
		lteHip     string
		gteHeight  string
		lteHeight  string
		gteBirthday string
		lteBirthday string
		hits       int
		offset     int
		sort       string
	)

	c := &cobra.Command{
		Use:   "actress",
		Short: "Search DMM actresses",
		RunE: func(c *cobra.Command, args []string) error {
			cl := cmd.LoadDMMClient()
			result, err := cl.ActressSearch(client.ActressSearchParams{
				Keyword:     keyword,
				ActressID:   actressID,
				Initial:     initial,
				GteBust:     gteBust,
				LteBust:     lteBust,
				GteWaist:    gteWaist,
				LteWaist:    lteWaist,
				GteHip:      gteHip,
				LteHip:      lteHip,
				GteHeight:   gteHeight,
				LteHeight:   lteHeight,
				GteBirthday: gteBirthday,
				LteBirthday: lteBirthday,
				Hits:        hits,
				Offset:      offset,
				Sort:        sort,
			})
			if err != nil {
				cmd.HandleError(err)
			}
			output.Print(result.Result, cmd.Pretty())
			return nil
		},
	}

	f := c.Flags()
	f.StringVar(&keyword, "keyword", "", "Name keyword search")
	f.StringVar(&actressID, "actress-id", "", "Specific actress ID")
	f.StringVar(&initial, "initial", "", "Kana prefix (あ〜ん)")
	f.StringVar(&gteBust, "gte-bust", "", "Bust >= value (cm)")
	f.StringVar(&lteBust, "lte-bust", "", "Bust <= value (cm)")
	f.StringVar(&gteWaist, "gte-waist", "", "Waist >= value (cm)")
	f.StringVar(&lteWaist, "lte-waist", "", "Waist <= value (cm)")
	f.StringVar(&gteHip, "gte-hip", "", "Hip >= value (cm)")
	f.StringVar(&lteHip, "lte-hip", "", "Hip <= value (cm)")
	f.StringVar(&gteHeight, "gte-height", "", "Height >= value (cm)")
	f.StringVar(&lteHeight, "lte-height", "", "Height <= value (cm)")
	f.StringVar(&gteBirthday, "gte-birthday", "", "Birthday on/after (yyyy-mm-dd)")
	f.StringVar(&lteBirthday, "lte-birthday", "", "Birthday on/before (yyyy-mm-dd)")
	f.IntVar(&hits, "hits", 20, "Results per page (1-100)")
	f.IntVar(&offset, "offset", 1, "Start position")
	f.StringVar(&sort, "sort", "", "Sort: name/-name/bust/-bust/waist/-waist/hip/-hip/height/-height/birthday/-birthday/id/-id")

	return c
}
