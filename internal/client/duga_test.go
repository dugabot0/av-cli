package client

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDUGAClientBaseParams(t *testing.T) {
	c := NewDUGAClient("myapp", "myagent", "01", 10*time.Second)
	p := c.baseParams()

	checks := map[string]string{
		"version":  "1.2",
		"appid":    "myapp",
		"agentid":  "myagent",
		"bannerid": "01",
		"format":   "json",
	}
	for k, want := range checks {
		if got := p.Get(k); got != want {
			t.Errorf("param %s: want %q, got %q", k, want, got)
		}
	}
}

func TestDUGASearchResponseParsing(t *testing.T) {
	// Matches the actual API response structure: items is array of {"item":{...}}
	raw := `{
		"hits": "2",
		"count": 100,
		"offset": "1",
		"timestamp": "2026/02/28 12:00:00",
		"items": [
			{
				"item": {
					"productid": "bigmorkal-0712",
					"title": "テスト作品",
					"caption": "説明文",
					"makername": "テストメーカー",
					"url": "https://duga.jp/ppv/bigmorkal-0712/",
					"affiliateurl": "https://click.duga.jp/ppv/bigmorkal-0712/XX",
					"opendate": "2013/04/25",
					"price": "980円",
					"volume": 240,
					"label": [{"id":"bigmorkal","name":"ビッグモーカル","number":"0712"}],
					"category": [{"data":{"id":"07","name":"熟女"}}],
					"saletype": [{"data":{"type":"通常版","price":"980"}}],
					"ranking": [{"total":"55"}],
					"mylist": [{"total":"1051"}]
				}
			}
		]
	}`
	var resp DUGASearchResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("items: want 1, got %d", len(resp.Items))
	}
	item := resp.Items[0].Item
	if item.ProductID != "bigmorkal-0712" {
		t.Errorf("productid: want bigmorkal-0712, got %q", item.ProductID)
	}
	if item.Title != "テスト作品" {
		t.Errorf("title: want テスト作品, got %q", item.Title)
	}
	if item.Volume == nil {
		t.Error("volume: should not be nil")
	}
	if item.Label == nil {
		t.Error("label: RawMessage should not be nil")
	}
}
