package duga

import (
	"github.com/spf13/cobra"
	"github.com/yourusername/av-cli/cmd"
	"github.com/yourusername/av-cli/internal/client"
	"github.com/yourusername/av-cli/internal/output"
)

func newSearchCmd() *cobra.Command {
	var (
		keyword      string
		hits         int
		offset       int
		sort         string
		category     string
		performerID  string
		seriesID     string
		labelID      string
		adult        string
		target       string
		openStart    string
		openEnd      string
		releaseStart string
		releaseEnd   string
	)

	c := &cobra.Command{
		Use:   "search",
		Short: "Search DUGA items",
		RunE: func(c *cobra.Command, args []string) error {
			cl := cmd.LoadDUGAClient()
			result, err := cl.Search(client.DUGASearchParams{
				Keyword:      keyword,
				Hits:         hits,
				Offset:       offset,
				Sort:         sort,
				Category:     category,
				PerformerID:  performerID,
				SeriesID:     seriesID,
				LabelID:      labelID,
				Adult:        adult,
				Target:       target,
				OpenStart:    openStart,
				OpenEnd:      openEnd,
				ReleaseStart: releaseStart,
				ReleaseEnd:   releaseEnd,
			})
			if err != nil {
				cmd.HandleError(err)
			}
			output.Print(result, cmd.Pretty())
			return nil
		},
	}

	f := c.Flags()
	f.StringVar(&keyword, "keyword", "", "Search keyword")
	f.IntVar(&hits, "hits", 10, "Results per page (1-100)")
	f.IntVar(&offset, "offset", 1, "Start position")
	f.StringVar(&sort, "sort", "", "Sort: favorite/release/new/price/rating/mylist")
	f.StringVar(&category, "category", "", "Category ID")
	f.StringVar(&performerID, "performer-id", "", "Filter by performer ID")
	f.StringVar(&seriesID, "series-id", "", "Filter by series ID")
	f.StringVar(&labelID, "label-id", "", "Filter by label ID")
	f.StringVar(&adult, "adult", "1", "1=adult, 0=general")
	f.StringVar(&target, "target", "", "Sale type: ppv/sd/rental/hd/hdrental")
	f.StringVar(&openStart, "open-start", "", "Open date start (YYYYMMDD)")
	f.StringVar(&openEnd, "open-end", "", "Open date end (YYYYMMDD)")
	f.StringVar(&releaseStart, "release-start", "", "Release date start (YYYYMMDD)")
	f.StringVar(&releaseEnd, "release-end", "", "Release date end (YYYYMMDD)")

	return c
}
