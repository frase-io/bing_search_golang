package bing

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	apiKey = ""
)

type BingWebSearchResult struct {
	Type            string           `json:"_type,omitempty"`
	WebPages        *WebPages        `json:"webPages,omitempty"`
	Images          *Images          `json:"images,omitempty"`
	News            *News            `json:"news,omitempty"`
	Videos          *Videos          `json:"videos,omitempty"`
	RankingResponse *RankingResponse `json:"rankingResponse,omitempty"`
	Sidebar         *Sidebar         `json:"sidebar,omitempty"`
}

type BingNewsSearchResult struct {
	Type                  string        `json:"_type"`
	ReadLink              string        `json:"readLink"`
	TotalEstimatedMatches int64         `json:"totalEstimatedMatches"`
	Sort                  []*SortType   `json:"sort"`
	Value                 []*NewsResult `json:"value"`
}

type SortType struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	IsSelected bool   `json:"isSelected"`
	URL        string `json:"url"`
}

type WebPages struct {
	WebSearchURL          string     `json:"webSearchUrl,omitempty"`
	TotalEstimatedMatches int64      `json:"totalEstimatedMatches,omitempty"`
	Value                 []*WebPage `json:"value,omitempty"`
}

// WebPage is a value in the webpages section of the bing response
type WebPage struct {
	ID              string   `json:"id,omitempty"`
	Name            string   `json:"name,omitempty"`
	URL             string   `json:"url,omitempty"`
	About           []*About `json:"about,omitempty"`
	DisplayURL      string   `json:"displayUrl,omitempty"`
	Snippet         string   `json:"snippet,omitempty"`
	DateLastCrawled string   `json:"dateLastCrawled,omitempty"`
}

// About is the about section of the result
type About struct {
	Name string `json:"name,omitempty"`
}

type Images struct {
	ID                           string   `json:"id,omitempty"`
	ReadLink                     string   `json:"readLink,omitempty"`
	WebSearchURL                 string   `json:"webSearchUrl,omitempty"`
	IsFamilyFriendly             bool     `json:"isFamilyFriendly,omitempty"`
	Value                        []*Image `json:"value, omitempty"`
	DisplayShoppingSourcesBadges bool     `json:"displayShoppingSourcesBadges,omitempty"`
	DisplayRecipeSourcesBadges   bool     `json:"displayRecipeSourcesBadges,omitempty"`
}

type Image struct {
	Name               string     `json:"name,omitempty"`
	WebSearchURL       string     `json:"webSearchUrl,omitempty"`
	ThumbnailURL       string     `json:"thumbnailUrl,omitempty"`
	DatePublished      string     `json:"datePublished,omitempty"`
	ContentURL         string     `json:"contentUrl,omitempty"`
	HostPageURL        string     `json:"hostPageUrl,omitempty"`
	ContentSize        string     `json:"contentSize,omitempty"`
	EncodingFormat     string     `json:"encodingFormat,omitempty"`
	HostPageDisplayURL string     `json:"hostPageDisplayUrl,omitempty"`
	Width              int64      `json:"width,omitempty"`
	Height             int64      `json:"height,omitempty"`
	Thumbnail          *Thumbnail `json:"thumbnail,omitempty"`
}

type Thumbnail struct {
	Width      int64  `json:"width,omitempty"`
	Height     int64  `json:"height,omitempty"`
	ContentURL string `json:"contentUrl,omitempty"`
}

// News results from bing search
type News struct {
	ID       string        `json:"id,omitempty"`
	ReadLink string        `json:"readLink,omitempty"`
	Value    []*NewsResult `json:"value,omitempty"`
}

type NewsResult struct {
	Name          string     `json:"name,omitempty"`
	URL           string     `json:"url,omitempty"`
	Image         *NewsImage `json:"image,omitempty"`
	Description   string     `json:"description,omitempty"`
	About         *AboutNews `json:"about,omitempty"`
	Provider      *Provider  `json:"provider,omitempty"`
	DatePublished string     `json:"datePublished,omitempty"`
	Category      string     `json:"category,omitempty"`
}

type NewsImage struct {
	ContentUrl string     `json:"contentUrl,omitempty"`
	Thumbnail  *Thumbnail `json:"thumbnail,omitempty"`
}

type AboutNews struct {
	Name     string `json:"name,omitempty"`
	ReadLink string `json:"readLink,omitempty"`
}

type Provider struct {
	Type string `json:"_type,omitempty"`
	Name string `json:"name,omitempty"`
}

type Videos struct {
	ID               string   `json:"id,omitempty"`
	ReadLink         string   `json:"readLink,omitempty"`
	WebSearchURL     string   `json:"webSearchUrl,omitempty"`
	IsFamilyFriendly bool     `json:"isFamilyFriendly,omitempty"`
	Value            []*Video `json:"value,omitempty"`
}

type Video struct {
	Name               string     `json:"name,omitempty"`
	Description        string     `json:"description,omitempty"`
	WebSearchURL       string     `json:"webSearchUrl,omitempty"`
	ThumbnailURL       string     `json:"thumbnailUrl,omitempty"`
	DatePublished      string     `json:"datePublished,omitempty"`
	Publisher          *Publisher `json:"publisher,omitempty"`
	ContentURL         string     `json:"contentUrl,omitempty"`
	HostPageURL        string     `json:"hostPageUrl,omitempty"`
	EncodingFormat     string     `json:"encodingFormat,omitempty"`
	HostPageDisplayURL string     `json:"hostPageDisplayUrl,omitempty"`
	Width              int64      `json:"width,omitempty"`
	Height             int64      `json:"height,omitempty"`
	Duration           string     `json:"duration,omitempty"`
	MotionThumbnailURL string     `json:"motionThumbnailUrl,omitempty"`
	EmbedHTML          string     `json:"embedHtml,omitempty"`
	AllowHTTPSEmbed    bool       `json:"allowHttpsEmbed,omitempty"`
	ViewCount          int64      `json:"viewCount,omitempty"`
	Thumbnail          *Thumbnail `json:"thumbnail,omitempty"`
	AllowMobileEmbed   bool       `json:"allowMobileEmbed,omitempty"`
	IsSuperfresh       bool       `json:"isSuperfresh,omitempty"`
}

type Publisher struct {
	Name string `json:"publisher,omitempty"`
}

type RankingResponse struct {
	Mainline []*Item `json:"mainline,omitempty"`
}

type Item struct {
	AnswerType  string `json:"answerType,omitempty"`
	ResultIndex int64  `json:"resultIndex,omitempty"`
	Value       *Value `json:"value,omitempty"`
}

type Value struct {
	ID string `json:"id,omitempty"`
}

type Sidebar struct {
	Items []*SidebarItem `json:"items,omitempty"`
}

type SidebarItem struct {
	AnswerType string `json:"answerType,omitempty"`
	Value      *Value `json:"value,omitempty"`
}

func (bingSearchResult *BingWebSearchResult) MakeBingRequest(query string, resultCount string) error {
	client := &http.Client{}
	newQuery := strings.Replace(query, " ", "+", -1)
	reqURL := "https://api.cognitive.microsoft.com/bing/v5.0/search?q=" + newQuery + "&offset=0&count=" + resultCount + "&mkt=" + config["bingMkt"]

	req, err := http.NewRequest("GET", reqURL, nil)
	req.Header.Set("Ocp-Apim-Subscription-Key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error: " + err.Error())
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &bingSearchResult)

	return nil
}

func (bingSearchResult *BingNewsSearchResult) MakeBingRequest(query string, resultCount string) error {
	client := &http.Client{}
	newQuery := strings.Replace(query, " ", "+", -1)
	reqURL := "https://api.cognitive.microsoft.com/bing/v5.0/news/search?q=" +
		newQuery + "&offset=0&count=" + resultCount + "&freshness=Month&mkt=" + config["bingMkt"]

	req, err := http.NewRequest("GET", reqURL, nil)
	req.Header.Set("Ocp-Apim-Subscription-Key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error: " + err.Error())
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &bingSearchResult)

	return nil
}
