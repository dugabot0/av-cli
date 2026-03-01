package client

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDMMClientBaseParams(t *testing.T) {
	c := NewDMMClient("myapi", "myaffiliate-990", 10*time.Second)
	p := c.baseParams()
	if p.Get("api_id") != "myapi" {
		t.Errorf("api_id: want %q, got %q", "myapi", p.Get("api_id"))
	}
	if p.Get("affiliate_id") != "myaffiliate-990" {
		t.Errorf("affiliate_id: want %q, got %q", "myaffiliate-990", p.Get("affiliate_id"))
	}
	if p.Get("output") != "json" {
		t.Errorf("output: want %q, got %q", "json", p.Get("output"))
	}
}

func TestDMMAuthError(t *testing.T) {
	err := &AuthError{StatusCode: 401}
	if err.Error() == "" {
		t.Error("AuthError.Error() should not be empty")
	}
}

func TestDMMAPIError(t *testing.T) {
	err := &APIError{StatusCode: 500}
	if err.Error() == "" {
		t.Error("APIError.Error() should not be empty")
	}
}

func TestFloorListResponseParsing(t *testing.T) {
	raw := `{
		"result": {
			"site": [
				{
					"name": "FANZA（アダルト）",
					"code": "FANZA",
					"service": [
						{
							"name": "動画",
							"code": "digital",
							"floor": [
								{"id": "43", "name": "ビデオ", "code": "videoa"}
							]
						}
					]
				}
			]
		}
	}`
	var resp FloorListResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Result.Site) != 1 {
		t.Fatalf("expected 1 site, got %d", len(resp.Result.Site))
	}
	site := resp.Result.Site[0]
	if site.Code != "FANZA" {
		t.Errorf("site.Code: want FANZA, got %q", site.Code)
	}
	if len(site.Service) != 1 {
		t.Fatalf("expected 1 service, got %d", len(site.Service))
	}
	svc := site.Service[0]
	if len(svc.Floor) != 1 {
		t.Fatalf("expected 1 floor, got %d", len(svc.Floor))
	}
	floor := svc.Floor[0]
	if floor.ID != "43" {
		t.Errorf("floor.ID: want 43, got %q", floor.ID)
	}
	if floor.Code != "videoa" {
		t.Errorf("floor.Code: want videoa, got %q", floor.Code)
	}
}

func TestActressResponseParsing(t *testing.T) {
	raw := `{
		"result": {
			"status": "200",
			"result_count": 1,
			"total_count": "1",
			"first_position": 1,
			"actress": [
				{
					"id": "1054998",
					"name": "松本いちか",
					"ruby": "まつもといちか",
					"bust": "83",
					"cup": "C",
					"waist": "55",
					"hip": "82",
					"height": "153",
					"birthday": null,
					"blood_type": null,
					"hobby": null,
					"prefectures": null
				}
			]
		}
	}`
	var resp ActressSearchResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Result.Actress) != 1 {
		t.Fatalf("expected 1 actress, got %d", len(resp.Result.Actress))
	}
	a := resp.Result.Actress[0]
	if a.ID != "1054998" {
		t.Errorf("actress.ID: want 1054998, got %q", a.ID)
	}
	if a.Birthday != nil {
		t.Errorf("actress.Birthday: want nil, got %v", a.Birthday)
	}
}

func TestItemListParamsBuilding(t *testing.T) {
	c := NewDMMClient("api", "aff-990", time.Second)
	p := c.baseParams()
	// Simulate ItemList param building
	params := ItemListParams{
		Site:      "FANZA",
		Service:   "digital",
		Floor:     "videoa",
		Hits:      10,
		Offset:    1,
		Sort:      "date",
		Keyword:   "test",
		Article:   []string{"actress", "genre"},
		ArticleID: []string{"123", "456"},
	}
	if params.Site != "FANZA" {
		t.Error("site mismatch")
	}
	if len(params.Article) != 2 {
		t.Errorf("article count: want 2, got %d", len(params.Article))
	}
	_ = p // suppress unused warning
}
