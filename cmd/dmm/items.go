package dmm

import (
	"github.com/spf13/cobra"
	"github.com/yourusername/av-cli/cmd"
	"github.com/yourusername/av-cli/internal/client"
	"github.com/yourusername/av-cli/internal/output"
)

func newItemsCmd() *cobra.Command {
	var (
		keyword   string
		site      string
		service   string
		floor     string
		hits      int
		offset    int
		sort      string
		article   []string
		articleID []string
		cid       string
	)

	c := &cobra.Command{
		Use:   "items",
		Short: "Search DMM items",
		RunE: func(c *cobra.Command, args []string) error {
			cl := cmd.LoadDMMClient()
			result, err := cl.ItemList(client.ItemListParams{
				Site:      site,
				Service:   service,
				Floor:     floor,
				Hits:      hits,
				Offset:    offset,
				Sort:      sort,
				Keyword:   keyword,
				CID:       cid,
				Article:   article,
				ArticleID: articleID,
			})
			if err != nil {
				cmd.HandleError(err)
			}
			output.Print(result.Result, cmd.Pretty())
			return nil
		},
	}

	f := c.Flags()
	f.StringVar(&keyword, "keyword", "", "Search keyword")
	f.StringVar(&site, "site", "FANZA", "Site: FANZA or DMM.com")
	f.StringVar(&service, "service", "", "Service code (from `dmm floors`)")
	f.StringVar(&floor, "floor", "", "Floor code (from `dmm floors`)")
	f.IntVar(&hits, "hits", 20, "Results per page (1-100)")
	f.IntVar(&offset, "offset", 1, "Start position")
	f.StringVar(&sort, "sort", "", "Sort: rank/price/-price/date/review/match")
	f.StringArrayVar(&article, "article", nil, "Filter type: actress/author/genre/series/maker (repeatable)")
	f.StringArrayVar(&articleID, "article-id", nil, "ID for article filter (repeatable)")
	f.StringVar(&cid, "cid", "", "Specific content ID")

	return c
}
