package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const dmmBaseURL = "https://api.dmm.com/affiliate/v3"

// DMMClient handles requests to the DMM Affiliate API v3.
type DMMClient struct {
	APIID       string
	AffiliateID string
	HTTP        *http.Client
}

// NewDMMClient constructs a DMMClient with the given timeout.
func NewDMMClient(apiID, affiliateID string, timeout time.Duration) *DMMClient {
	return &DMMClient{
		APIID:       apiID,
		AffiliateID: affiliateID,
		HTTP:        &http.Client{Timeout: timeout},
	}
}

func (c *DMMClient) baseParams() url.Values {
	p := url.Values{}
	p.Set("api_id", c.APIID)
	p.Set("affiliate_id", c.AffiliateID)
	p.Set("output", "json")
	return p
}

func (c *DMMClient) get(endpoint string, params url.Values) (json.RawMessage, error) {
	u := fmt.Sprintf("%s/%s?%s", dmmBaseURL, endpoint, params.Encode())
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

// --- Error types ---

type AuthError struct{ StatusCode int }

func (e *AuthError) Error() string { return fmt.Sprintf("authentication error (HTTP %d)", e.StatusCode) }

type APIError struct{ StatusCode int }

func (e *APIError) Error() string { return fmt.Sprintf("API error (HTTP %d)", e.StatusCode) }

// --- Floor API ---

type FloorListResponse struct {
	Result FloorResult `json:"result"`
}

type FloorResult struct {
	Site []Site `json:"site"`
}

type Site struct {
	Name    string    `json:"name"`
	Code    string    `json:"code"`
	Service []Service `json:"service"`
}

type Service struct {
	Name  string  `json:"name"`
	Code  string  `json:"code"`
	Floor []Floor `json:"floor"`
}

type Floor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

func (c *DMMClient) FloorList() (*FloorListResponse, error) {
	p := c.baseParams()
	raw, err := c.get("FloorList", p)
	if err != nil {
		return nil, err
	}
	var result FloorListResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parse floor list: %w", err)
	}
	return &result, nil
}

// --- Item API ---

type ItemListParams struct {
	Site      string
	Service   string
	Floor     string
	Hits      int
	Offset    int
	Sort      string
	Keyword   string
	CID       string
	Article   []string
	ArticleID []string
}

type ItemListResponse struct {
	Result ItemResult `json:"result"`
}

type ItemResult struct {
	Status        any    `json:"status"`
	ResultCount   any    `json:"result_count"`
	TotalCount    any    `json:"total_count"`
	FirstPosition any    `json:"first_position"`
	Items         []Item `json:"items"`
}

type Item struct {
	ServiceCode  string      `json:"service_code"`
	ServiceName  string      `json:"service_name"`
	FloorCode    string      `json:"floor_code"`
	FloorName    string      `json:"floor_name"`
	CategoryName string      `json:"category_name"`
	ContentID    string      `json:"content_id"`
	ProductID    string      `json:"product_id"`
	Title        string      `json:"title"`
	Volume       string      `json:"volume"`
	Number       string      `json:"number"`
	Review       *ItemReview `json:"review,omitempty"`
	AffiliateURL string      `json:"affiliateURL"`
	URL          string      `json:"URL"`
	ImageURL     *ImageURL   `json:"imageURL,omitempty"`
	SampleImages *SampleImages `json:"sampleImageURL,omitempty"`
	Prices       *ItemPrices `json:"prices,omitempty"`
	Date         string      `json:"date"`
}

type ItemReview struct {
	Count   any `json:"count"`
	Average any `json:"average"`
}

type ImageURL struct {
	List  string `json:"list"`
	Small string `json:"small"`
	Large string `json:"large"`
}

type SampleImages struct {
	SampleS *SampleImageList `json:"sample_s,omitempty"`
}

type SampleImageList struct {
	Image []string `json:"image"`
}

type ItemPrices struct {
	Price      any         `json:"price"`
	ListPrice  any         `json:"list_price"`
	Deliveries *Deliveries `json:"deliveries,omitempty"`
}

type Deliveries struct {
	Delivery []Delivery `json:"delivery"`
}

type Delivery struct {
	Type  string `json:"type"`
	Price any    `json:"price"`
}

func (c *DMMClient) ItemList(params ItemListParams) (*ItemListResponse, error) {
	p := c.baseParams()
	if params.Site != "" {
		p.Set("site", params.Site)
	} else {
		p.Set("site", "FANZA")
	}
	if params.Service != "" {
		p.Set("service", params.Service)
	}
	if params.Floor != "" {
		p.Set("floor", params.Floor)
	}
	if params.Hits > 0 {
		p.Set("hits", fmt.Sprint(params.Hits))
	}
	if params.Offset > 0 {
		p.Set("offset", fmt.Sprint(params.Offset))
	}
	if params.Sort != "" {
		p.Set("sort", params.Sort)
	}
	if params.Keyword != "" {
		p.Set("keyword", params.Keyword)
	}
	if params.CID != "" {
		p.Set("cid", params.CID)
	}
	for i, a := range params.Article {
		p.Set(fmt.Sprintf("article[%d]", i), a)
	}
	for i, id := range params.ArticleID {
		p.Set(fmt.Sprintf("article_id[%d]", i), id)
	}

	raw, err := c.get("ItemList", p)
	if err != nil {
		return nil, err
	}
	var result ItemListResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parse item list: %w", err)
	}
	return &result, nil
}

// --- Actress API ---

type ActressSearchParams struct {
	Keyword    string
	ActressID  string
	Initial    string
	GteBust    string
	LteBust    string
	GteWaist   string
	LteWaist   string
	GteHip     string
	LteHip     string
	GteHeight  string
	LteHeight  string
	GteBirthday string
	LteBirthday string
	Hits       int
	Offset     int
	Sort       string
}

type ActressSearchResponse struct {
	Result ActressResult `json:"result"`
}

type ActressResult struct {
	Status        any       `json:"status"`
	ResultCount   any       `json:"result_count"`
	TotalCount    any       `json:"total_count"`
	FirstPosition any       `json:"first_position"`
	Actress       []Actress `json:"actress"`
}

type Actress struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Ruby        string       `json:"ruby"`
	Bust        string       `json:"bust"`
	Cup         string       `json:"cup"`
	Waist       string       `json:"waist"`
	Hip         string       `json:"hip"`
	Height      string       `json:"height"`
	Birthday    *string      `json:"birthday"`
	BloodType   *string      `json:"blood_type"`
	Hobby       string       `json:"hobby"`
	Prefectures *string      `json:"prefectures"`
	ImageURL    *ActressImage `json:"imageURL,omitempty"`
	ListURL     *ActressListURL `json:"listURL,omitempty"`
}

type ActressImage struct {
	Small string `json:"small"`
	Large string `json:"large"`
}

type ActressListURL struct {
	Digital string `json:"digital"`
	Monthly string `json:"monthly"`
	Mono    string `json:"mono"`
}

func (c *DMMClient) ActressSearch(params ActressSearchParams) (*ActressSearchResponse, error) {
	p := c.baseParams()
	setIfNonEmpty(p, "keyword", params.Keyword)
	setIfNonEmpty(p, "actress_id", params.ActressID)
	setIfNonEmpty(p, "initial", params.Initial)
	setIfNonEmpty(p, "gte_bust", params.GteBust)
	setIfNonEmpty(p, "lte_bust", params.LteBust)
	setIfNonEmpty(p, "gte_waist", params.GteWaist)
	setIfNonEmpty(p, "lte_waist", params.LteWaist)
	setIfNonEmpty(p, "gte_hip", params.GteHip)
	setIfNonEmpty(p, "lte_hip", params.LteHip)
	setIfNonEmpty(p, "gte_height", params.GteHeight)
	setIfNonEmpty(p, "lte_height", params.LteHeight)
	setIfNonEmpty(p, "gte_birthday", params.GteBirthday)
	setIfNonEmpty(p, "lte_birthday", params.LteBirthday)
	if params.Hits > 0 {
		p.Set("hits", fmt.Sprint(params.Hits))
	}
	if params.Offset > 0 {
		p.Set("offset", fmt.Sprint(params.Offset))
	}
	setIfNonEmpty(p, "sort", params.Sort)

	raw, err := c.get("ActressSearch", p)
	if err != nil {
		return nil, err
	}
	var result ActressSearchResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parse actress search: %w", err)
	}
	return &result, nil
}

// --- Genre API ---

type GenreSearchParams struct {
	FloorID string
	Initial string
	Hits    int
	Offset  int
}

type GenreSearchResponse struct {
	Result GenreResult `json:"result"`
}

type GenreResult struct {
	Status        any     `json:"status"`
	ResultCount   any     `json:"result_count"`
	TotalCount    any     `json:"total_count"`
	FirstPosition any     `json:"first_position"`
	Genre         []Genre `json:"genre"`
}

type Genre struct {
	GenreID string `json:"genre_id"`
	Name    string `json:"name"`
	Ruby    string `json:"ruby"`
	ListURL string `json:"list_url"`
}

func (c *DMMClient) GenreSearch(params GenreSearchParams) (*GenreSearchResponse, error) {
	p := c.baseParams()
	p.Set("floor_id", params.FloorID)
	setIfNonEmpty(p, "initial", params.Initial)
	if params.Hits > 0 {
		p.Set("hits", fmt.Sprint(params.Hits))
	}
	if params.Offset > 0 {
		p.Set("offset", fmt.Sprint(params.Offset))
	}

	raw, err := c.get("GenreSearch", p)
	if err != nil {
		return nil, err
	}
	var result GenreSearchResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parse genre search: %w", err)
	}
	return &result, nil
}

// --- Maker API ---

type MakerSearchParams struct {
	FloorID string
	Initial string
	Hits    int
	Offset  int
}

type MakerSearchResponse struct {
	Result MakerResult `json:"result"`
}

type MakerResult struct {
	Status        any     `json:"status"`
	ResultCount   any     `json:"result_count"`
	TotalCount    any     `json:"total_count"`
	FirstPosition any     `json:"first_position"`
	Maker         []Maker `json:"maker"`
}

type Maker struct {
	MakerID string `json:"maker_id"`
	Name    string `json:"name"`
	Ruby    string `json:"ruby"`
	ListURL string `json:"list_url"`
}

func (c *DMMClient) MakerSearch(params MakerSearchParams) (*MakerSearchResponse, error) {
	p := c.baseParams()
	p.Set("floor_id", params.FloorID)
	setIfNonEmpty(p, "initial", params.Initial)
	if params.Hits > 0 {
		p.Set("hits", fmt.Sprint(params.Hits))
	}
	if params.Offset > 0 {
		p.Set("offset", fmt.Sprint(params.Offset))
	}

	raw, err := c.get("MakerSearch", p)
	if err != nil {
		return nil, err
	}
	var result MakerSearchResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parse maker search: %w", err)
	}
	return &result, nil
}

// --- Series API ---

type SeriesSearchParams struct {
	FloorID string
	Initial string
	Hits    int
	Offset  int
}

type SeriesSearchResponse struct {
	Result SeriesResult `json:"result"`
}

type SeriesResult struct {
	Status        any      `json:"status"`
	ResultCount   any      `json:"result_count"`
	TotalCount    any      `json:"total_count"`
	FirstPosition any      `json:"first_position"`
	Series        []Series `json:"series"`
}

type Series struct {
	SeriesID string `json:"series_id"`
	Name     string `json:"name"`
	Ruby     string `json:"ruby"`
	ListURL  string `json:"list_url"`
}

func (c *DMMClient) SeriesSearch(params SeriesSearchParams) (*SeriesSearchResponse, error) {
	p := c.baseParams()
	p.Set("floor_id", params.FloorID)
	setIfNonEmpty(p, "initial", params.Initial)
	if params.Hits > 0 {
		p.Set("hits", fmt.Sprint(params.Hits))
	}
	if params.Offset > 0 {
		p.Set("offset", fmt.Sprint(params.Offset))
	}

	raw, err := c.get("SeriesSearch", p)
	if err != nil {
		return nil, err
	}
	var result SeriesSearchResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parse series search: %w", err)
	}
	return &result, nil
}

// --- Author API ---

type AuthorSearchParams struct {
	FloorID string
	Initial string
	Hits    int
	Offset  int
}

type AuthorSearchResponse struct {
	Result AuthorResult `json:"result"`
}

type AuthorResult struct {
	Status        any      `json:"status"`
	ResultCount   any      `json:"result_count"`
	TotalCount    any      `json:"total_count"`
	FirstPosition any      `json:"first_position"`
	Author        []Author `json:"author"`
}

type Author struct {
	AuthorID string `json:"author_id"`
	Name     string `json:"name"`
	Ruby     string `json:"ruby"`
	ListURL  string `json:"list_url"`
}

func (c *DMMClient) AuthorSearch(params AuthorSearchParams) (*AuthorSearchResponse, error) {
	p := c.baseParams()
	p.Set("floor_id", params.FloorID)
	setIfNonEmpty(p, "initial", params.Initial)
	if params.Hits > 0 {
		p.Set("hits", fmt.Sprint(params.Hits))
	}
	if params.Offset > 0 {
		p.Set("offset", fmt.Sprint(params.Offset))
	}

	raw, err := c.get("AuthorSearch", p)
	if err != nil {
		return nil, err
	}
	var result AuthorSearchResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parse author search: %w", err)
	}
	return &result, nil
}

func setIfNonEmpty(p url.Values, key, val string) {
	if val != "" {
		p.Set(key, val)
	}
}
