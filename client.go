// Package hackernews is a simple HTTP client for Hacker News.
//
// Algolia graciously provided an API for working with Hacker News over at:
// https://hn.algolia.com/api
package hackernews

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const baseURL = `http://hn.algolia.com/api/v1`

// New HackerNews Client with defaults
func New() *Client {
	return &Client{http.DefaultClient}
}

// Client for HackerNews. The HTTP Client can be overriden with your own.
type Client struct {
	*http.Client
}

// FrontPage is a convenience function for getting the results on
// https://hackernews.com
func (c *Client) FrontPage(ctx context.Context) ([]*Story, error) {
	result, err := c.Search(ctx, &SearchRequest{
		Tags:           "front_page",
		ResultsPerPage: 34,
	})
	if err != nil {
		return nil, err
	}
	return result.Stories, nil
}

// Newest is a convenience function for getting the results on
// https://news.ycombinator.com/newest
func (c *Client) Newest(ctx context.Context) ([]*Story, error) {
	result, err := c.SearchRecent(ctx, &SearchRequest{
		Tags:           "story",
		ResultsPerPage: 34,
	})
	if err != nil {
		return nil, err
	}
	return result.Stories, nil
}

// AskHN is a convenience function for getting the results on
// https://news.ycombinator.com/ask
func (c *Client) AskHN(ctx context.Context) ([]*Story, error) {
	result, err := c.SearchRecent(ctx, &SearchRequest{
		Tags:           "ask_hn",
		ResultsPerPage: 34,
	})
	if err != nil {
		return nil, err
	}
	return result.Stories, nil
}

// ShowHN is a convenience function for getting the results on
// https://news.ycombinator.com/show
func (c *Client) ShowHN(ctx context.Context) ([]*Story, error) {
	result, err := c.SearchRecent(ctx, &SearchRequest{
		Tags:           "show_hn",
		ResultsPerPage: 34,
	})
	if err != nil {
		return nil, err
	}
	return result.Stories, nil
}

// Story is an individual entry on HackerNews.
type Story struct {
	ID          int        `json:"id,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	CreatedAtI  int        `json:"created_at_i,omitempty"`
	Type        string     `json:"type,omitempty"`
	Author      string     `json:"author,omitempty"`
	Title       string     `json:"title,omitempty"`
	URL         string     `json:"url,omitempty"`
	Text        *string    `json:"text,omitempty"`
	NumComments *int       `json:"num_comments,omitempty"`
	Points      int        `json:"points,omitempty"`
	ParentID    *int       `json:"parent_id,omitempty"`
	StoryID     *int       `json:"story_id,omitempty"`
	Children    []Children `json:"children"`
}

// Children are the comments.
type Children struct {
	ID         int        `json:"id,omitempty"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	CreatedAtI int        `json:"created_at_i,omitempty"`
	Type       string     `json:"type,omitempty"`
	Author     *string    `json:"author,omitempty"`
	Title      *string    `json:"title,omitempty"`
	URL        *string    `json:"url,omitempty"`
	Text       *string    `json:"text,omitempty"`
	Points     *int       `json:"points,omitempty"`
	ParentID   int        `json:"parent_id,omitempty"`
	StoryID    int        `json:"story_id,omitempty"`
	Children   []Children `json:"children"`
}

// Find a Story by its id.
func (c *Client) Find(ctx context.Context, id int) (*Story, error) {
	url := fmt.Sprintf("%s/items/%d", baseURL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status %d: %s", res.StatusCode, string(body))
	}
	story := new(Story)
	if err := json.Unmarshal(body, story); err != nil {
		return nil, err
	}
	story.Children = filterChildren(story.Children)
	recursivelySort(story.Children)
	return story, nil
}

// Some comments are nil for some reason (perhaps removed?)
func filterChildren(childs []Children) (children []Children) {
	for _, child := range childs {
		if child.Author == nil || child.Text == nil {
			continue
		}
		child.Children = filterChildren(child.Children)
		children = append(children, child)
	}
	return children
}

func recursivelySort(children []Children) {
	sort.Slice(children, func(a, b int) bool {
		return children[a].CreatedAtI < children[b].CreatedAtI
	})
	for _, child := range children {
		recursivelySort(child.Children)
	}
}

// SearchRequest query and filters
type SearchRequest struct {
	// Full-text query to search for (e.g. Duo)
	Query string

	// Tags filters the search on a specific tag.
	//
	// The available tags are:
	//   - story
	//   - comment
	//   - poll
	//   - pollopt
	//   - show_hn
	//   - ask_hn,
	//   - front_page
	//   - author_:USERNAME
	//   - story_:ID
	//
	// Tags are ANDed by default, can be ORed if between parenthesis. For example,
	// `author_pg,(story,poll)` filters on `author=pg AND (type=story OR type=poll)`
	Tags string

	// Filter by points. Points is a conditional query, so you can request stories
	// that have more than 500 points with "points > 500".
	Points string

	// Filter by date. CreatedAt is a conditional query, so you can request
	// stories between a time period wtih "created_at_i>X,created_at_i<Y", where
	// X and Y are timestamps in seconds.
	CreatedAt string

	// Filter by the number of comments. Comments is a conditional query, so you
	// can request stories that have more than 10 comments with "comments > 10".
	NumComments string

	// The page number
	Page int

	// ResultsPerPage is the number of results. Defaults to 34.
	ResultsPerPage int
}

// Turns the search input into a query string.
func (s *SearchRequest) querystring() string {
	query := url.Values{}
	if s.Query != "" {
		query.Set("query", s.Query)
	}
	if s.Tags != "" {
		query.Set("tags", s.Tags)
	}
	if s.Page > 0 {
		query.Set("page", strconv.Itoa(s.Page))
	}
	var nfs []string
	if s.Points != "" {
		nfs = append(nfs, injectKey(s.Points, "points"))
	}
	if s.CreatedAt != "" {
		nfs = append(nfs, injectKey(s.CreatedAt, "created_at_i"))
	}
	if s.NumComments != "" {
		nfs = append(nfs, injectKey(s.NumComments, "num_comments"))
	}
	if len(nfs) > 0 {
		query.Set("numericFilters", strings.Join(nfs, ","))
	}
	// Set the number of results per page
	if s.ResultsPerPage > 0 {
		query.Set("hitsPerPage", strconv.Itoa(s.ResultsPerPage))
	}
	return query.Encode()
}

// Sugar on top to allow bot "points > 500" and "> 500", to reduce repetition
// with the key (e.g. Points: "points > 500")
func injectKey(query, key string) string {
	parts := strings.Split(query, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
		if !strings.HasPrefix(parts[i], key) {
			parts[i] = key + parts[i]
		}
	}
	return strings.Join(parts, ",")
}

// SearchResponse of a search
type SearchResponse struct {
	Stories              []*Story `json:"stories,omitempty"`
	Hits                 []*Hit   `json:"hits,omitempty"`
	NumResults           int      `json:"nbHits,omitempty"`
	Page                 int      `json:"page,omitempty"`
	NumPages             int      `json:"nbPages,omitempty"`
	ResultsPerPage       int      `json:"hitsPerPage,omitempty"`
	ExhaustiveNumResults bool     `json:"exhaustiveNbHits,omitempty"`
	Query                string   `json:"query,omitempty"`
	Params               string   `json:"params,omitempty"`
	ProcessingTimeMS     int      `json:"processingTimeMS,omitempty"`
}

func toStories(s *SearchResponse) ([]*Story, error) {
	stories := make([]*Story, len(s.Hits))
	for i, story := range s.Hits {
		id, err := strconv.Atoi(story.ID)
		if err != nil {
			return nil, err
		}
		stories[i] = &Story{
			Author:      story.Author,
			Children:    []Children{},
			CreatedAt:   story.CreatedAt,
			CreatedAtI:  story.CreatedAtI,
			ID:          id,
			NumComments: story.NumComments,
			ParentID:    story.ParentID,
			Points:      story.Points,
			StoryID:     story.StoryID,
			Title:       story.Title,
			Text:        nil,
			URL:         story.URL,
		}
	}
	return stories, nil
}

// Hit is an individual search result (story or comment)
type Hit struct {
	ID             string    `json:"objectID,omitempty"`
	Title          string    `json:"title,omitempty"`
	URL            string    `json:"url,omitempty"`
	Author         string    `json:"author,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	Points         int       `json:"points,omitempty"`
	StoryText      *string   `json:"story_text,omitempty"`
	CommentText    *string   `json:"comment_text,omitempty"`
	NumComments    *int      `json:"num_comments,omitempty"`
	StoryID        *int      `json:"story_id,omitempty"`
	StoryTitle     *string   `json:"story_title,omitempty"`
	StoryURL       *string   `json:"story_url,omitempty"`
	ParentID       *int      `json:"parent_id,omitempty"`
	CreatedAtI     int       `json:"created_at_i,omitempty"`
	RelevancyScore *int      `json:"relevancy_score,omitempty"`
	Tags           []string  `json:"_tags,omitempty"`
	Highlights     struct {
		Title     Highlight `json:"title,omitempty"`
		URL       Highlight `json:"url,omitempty"`
		Author    Highlight `json:"author,omitempty"`
		StoryText Highlight `json:"story_text,omitempty"`
	} `json:"_highlightResult,omitempty"`
	Children []int `json:"children"`
}

// Highlight indicates the words that matched the search query
type Highlight struct {
	Value        string   `json:"value,omitempty"`
	MatchLevel   string   `json:"matchLevel,omitempty"`
	MatchedWords []string `json:"matchedWords,omitempty"`
}

// Search for Stories. Sorted by relevance, then points, then number of comments.
func (c *Client) Search(ctx context.Context, search *SearchRequest) (*SearchResponse, error) {
	if search.Page >= 1 {
		search.Page = search.Page - 1
	}
	url := baseURL + "/search?" + search.querystring()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status %d: %s", res.StatusCode, string(body))
	}
	result := new(SearchResponse)
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}
	result.Page++
	// Convert the hits to stories
	stories, err := toStories(result)
	if err != nil {
		return nil, fmt.Errorf("failed to convert hits to stories: %w", err)
	}
	result.Stories = stories
	return result, nil
}

// Search for Stories. Sorted by date, more recent first.
func (c *Client) SearchRecent(ctx context.Context, search *SearchRequest) (*SearchResponse, error) {
	url := baseURL + "/search_by_date?" + search.querystring()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status %d: %s", res.StatusCode, string(body))
	}
	result := new(SearchResponse)
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}
	// Convert the hits to stories
	stories, err := toStories(result)
	if err != nil {
		return nil, fmt.Errorf("failed to convert hits to stories: %w", err)
	}
	result.Stories = stories
	return result, nil
}
