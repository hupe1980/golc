package integration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const baseURL = "https://%s.wikipedia.org/w/api.php"

type WikipediaOptions struct {
	LanguageCode string
	TopK         int
	DocMaxChars  int
	HTTPClient   HTTPClient
}

type Wikipedia struct {
	opts WikipediaOptions
}

func NewWikipedia(optFns ...func(o *WikipediaOptions)) *Wikipedia {
	opts := WikipediaOptions{
		LanguageCode: "en",
		TopK:         3,
		DocMaxChars:  4000,
		HTTPClient:   http.DefaultClient,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Wikipedia{
		opts: opts,
	}
}

func (w *Wikipedia) Run(ctx context.Context, query string) (string, error) {
	searchResult, err := w.search(ctx, query)
	if err != nil {
		return "", nil
	}

	summaries := []string{}

	for i, s := range searchResult.Query.Search {
		if i > w.opts.TopK {
			break
		}

		fetchPageResult, err := w.fetchPage(ctx, s.PageID)
		if err != nil {
			return "", err
		}

		page, ok := fetchPageResult.Query.Pages[fmt.Sprintf("%v", s.PageID)]
		if !ok {
			return "", errors.New("unexpected result from wikipedia api")
		}

		summaries = append(summaries, fmt.Sprintf("Page: %s\nSummary: %s", page.Title, page.Extract))
	}

	if len(summaries) == 0 {
		return "No good Wikipedia Search Result was found", nil
	}

	result := strings.Join(summaries, "\n\n")

	if len(result) >= w.opts.DocMaxChars {
		return result[:w.opts.DocMaxChars], nil
	}

	return result, nil
}

type searchResponse struct {
	Query struct {
		Search []struct {
			Ns        int       `json:"ns"`
			Title     string    `json:"title"`
			PageID    int       `json:"pageid"`
			Size      int       `json:"size"`
			WordCount int       `json:"wordcount"`
			Snippet   string    `json:"snippet"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"search"`
	} `json:"query"`
}

func (w *Wikipedia) search(ctx context.Context, query string) (*searchResponse, error) {
	params := make(url.Values)
	params.Add("format", "json")
	params.Add("action", "query")
	params.Add("list", "search")
	params.Add("srsearch", query)

	reqURL := fmt.Sprintf("%s?%s", fmt.Sprintf(baseURL, w.opts.LanguageCode), params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := w.opts.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result searchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (w *Wikipedia) fetchPage(ctx context.Context, pageID int) (*pageResult, error) {
	params := make(url.Values)
	params.Add("format", "json")
	params.Add("action", "query")
	params.Add("prop", "extracts")
	params.Add("pageids", fmt.Sprintf("%v", pageID))

	reqURL := fmt.Sprintf("%s?%s", fmt.Sprintf(baseURL, w.opts.LanguageCode), params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result pageResult

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type pageResult struct {
	Query struct {
		Pages map[string]struct {
			Title   string `json:"title"`
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}
