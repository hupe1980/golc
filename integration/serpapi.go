package integration

import (
	"context"

	"github.com/hupe1980/golc/internal/util"
	g "github.com/serpapi/google-search-results-golang"
)

type SerpAPIOptions struct {
	Parameter map[string]string
}

type SerpAPI struct {
	engine    string
	apiKey    string
	parameter map[string]string
}

func NewSerpAPI(apiKey string) (*SerpAPI, error) {
	opts := SerpAPIOptions{
		Parameter: map[string]string{
			"engine":        "google",
			"google_domain": "google.com",
			"gl":            "us",
			"hl":            "en",
		},
	}

	engine := opts.Parameter["engine"]

	return &SerpAPI{
		engine:    engine,
		apiKey:    apiKey,
		parameter: opts.Parameter,
	}, nil
}

func (s *SerpAPI) Run(ctx context.Context, query string) (string, error) {
	params := util.CopyMap(s.parameter)
	params["q"] = query
	params["api_key"] = s.apiKey

	search := g.NewSearch(s.engine, params, s.apiKey)

	res, err := search.GetJSON()
	if err != nil {
		return "", err
	}

	return s.processResponse(res), nil
}

func (s *SerpAPI) processResponse(res map[string]any) string {
	if _, ok := res["error"]; ok {
		return res["error"].(string)
	}

	if answerBox, ok := res["answer_box"]; ok {
		if answer, ok := answerBox.(map[string]any)["answer"]; ok {
			return answer.(string)
		}
	}

	if answerBox, ok := res["answer_box"]; ok {
		if answer, ok := answerBox.(map[string]any)["snippet"]; ok {
			return answer.(string)
		}
	}

	if answerBox, ok := res["answer_box"]; ok {
		if answer, ok := answerBox.(map[string]any)["snippet_highlighted_words"]; ok {
			return answer.([]string)[0]
		}
	}

	// TODO

	return "No good search result found"
}
