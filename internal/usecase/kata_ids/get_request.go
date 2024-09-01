package kata_ids

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CodewarsResponse struct {
	TotalPages int `json:"totalPages"`
	TotalItems int `json:"totalItems"`
	Data       []struct {
		ID                 string    `json:"id"`
		Name               string    `json:"name"`
		Slug               string    `json:"slug"`
		CompletedAt        time.Time `json:"completedAt"`
		CompletedLanguages []string  `json:"completedLanguages"`
	} `json:"data"`
}

func SendGetRequest(name string, page int) (*CodewarsResponse, error) {
	url := fmt.Sprintf("http://www.codewars.com/api/v1/users/%s/code-challenges/completed?page=%d", name, page)
	var resp *http.Response
	var err error

	delay := 1 * time.Second
	maxDelay := 8 * time.Second

	for {
		resp, err = http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			break
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			select {
			case <-time.After(delay):
				if delay < maxDelay {
					delay += 3 * time.Second
				}
			}
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var codewarsResponse CodewarsResponse
	err = json.Unmarshal(body, &codewarsResponse)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return &codewarsResponse, nil
}
