package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const dugaBaseURL = "http://affapi.duga.jp/search"

// DUGAClient handles requests to the DUGA Affiliate Web API.
type DUGAClient struct {
	AppID    string
	AgentID  string
	BannerID string
	HTTP     *http.Client
}

// NewDUGAClient constructs a DUGAClient with the given timeout.
func NewDUGAClient(appID, agentID, bannerID string, timeout time.Duration) *DUGAClient {
	return &DUGAClient{
		AppID:    appID,
		AgentID:  agentID,
		BannerID: bannerID,
		HTTP:     &http.Client{Timeout: timeout},
	}
}

func (c *DUGAClient) baseParams() url.Values {
	p := url.Values{}
	p.Set("version", "1.2")
	p.Set("appid", c.AppID)
	p.Set("agentid", c.AgentID)
	p.Set("bannerid", c.BannerID)
	p.Set("format", "json")
	return p
}

func (c *DUGAClient) get(params url.Values) (json.RawMessage, error) {
	u := dugaBaseURL + "?" + params.Encode()
	resp, err := c.HTTP.Get(u)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// ok
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, &AuthError{StatusCode: resp.StatusCode}
	default:
		return nil, &APIError{StatusCode: resp.StatusCode}
	}

	var raw json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return raw, nil
}

// --- Search ---

type DUGASearchParams struct {
	Keyword      string
	Hits         int
	Offset       int
	Sort         string
	Category     string
	PerformerID  string
	SeriesID     string
	LabelID      string
	Adult        string
	Target       string
	OpenStart    string
	OpenEnd      string
	ReleaseStart string
	ReleaseEnd   string
}

// DUGASearchResponse mirrors the actual API response.
// Sub-fields use json.RawMessage to pass through the API's array-of-objects
// format unchanged (e.g. posterimage, label, category, saletype, etc.).
type DUGASearchResponse struct {
	Hits      any              `json:"hits"`
	Count     any              `json:"count"`
	Offset    any              `json:"offset"`
	Timestamp string           `json:"timestamp"`
	Items     []DUGAItemWrapper `json:"items"`
}

// DUGAItemWrapper is the {"item": {...}} envelope the API wraps each result in.
type DUGAItemWrapper struct {
	Item DUGAItem `json:"item"`
}

// DUGAItem contains a single search result. Complex sub-fields that the API
// returns as unusual array-of-objects structures are left as RawMessage so
// agents receive the exact API payload.
type DUGAItem struct {
	ProductID    string          `json:"productid"`
	Title        string          `json:"title"`
	Caption      string          `json:"caption"`
	MakerName    string          `json:"makername"`
	URL          string          `json:"url"`
	AffiliateURL string          `json:"affiliateurl"`
	OpenDate     string          `json:"opendate"`
	ReleaseDate  string          `json:"releasedate,omitempty"`
	ItemNo       string          `json:"itemno,omitempty"`
	Price        string          `json:"price"`
	Volume       any             `json:"volume,omitempty"`
	PosterImage  json.RawMessage `json:"posterimage,omitempty"`
	JacketImage  json.RawMessage `json:"jacketimage,omitempty"`
	Thumbnail    json.RawMessage `json:"thumbnail,omitempty"`
	SampleMovie  json.RawMessage `json:"samplemovie,omitempty"`
	Label        json.RawMessage `json:"label,omitempty"`
	Category     json.RawMessage `json:"category,omitempty"`
	Series       json.RawMessage `json:"series,omitempty"`
	Performer    json.RawMessage `json:"performer,omitempty"`
	Director     json.RawMessage `json:"director,omitempty"`
	SaleType     json.RawMessage `json:"saletype,omitempty"`
	Ranking      json.RawMessage `json:"ranking,omitempty"`
	Review       json.RawMessage `json:"review,omitempty"`
	MyList       json.RawMessage `json:"mylist,omitempty"`
}

func (c *DUGAClient) Search(params DUGASearchParams) (*DUGASearchResponse, error) {
	p := c.baseParams()
	setIfNonEmpty(p, "keyword", params.Keyword)
	if params.Hits > 0 {
		p.Set("hits", fmt.Sprint(params.Hits))
	}
	if params.Offset > 0 {
		p.Set("offset", fmt.Sprint(params.Offset))
	}
	setIfNonEmpty(p, "sort", params.Sort)
	setIfNonEmpty(p, "category", params.Category)
	setIfNonEmpty(p, "performerid", params.PerformerID)
	setIfNonEmpty(p, "seriesid", params.SeriesID)
	setIfNonEmpty(p, "labelid", params.LabelID)
	setIfNonEmpty(p, "adult", params.Adult)
	setIfNonEmpty(p, "target", params.Target)
	setIfNonEmpty(p, "openstt", params.OpenStart)
	setIfNonEmpty(p, "openend", params.OpenEnd)
	setIfNonEmpty(p, "releasestt", params.ReleaseStart)
	setIfNonEmpty(p, "releaseend", params.ReleaseEnd)

	raw, err := c.get(p)
	if err != nil {
		return nil, err
	}
	var result DUGASearchResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parse DUGA search: %w", err)
	}
	return &result, nil
}
