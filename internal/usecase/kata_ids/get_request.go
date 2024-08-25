package kata_ids

import (
	"fmt"
	"io"
	"net/http"
)

func SendGetRequest(name string, page int) ([]byte, error) {
	url := fmt.Sprintf("http://www.codewars.com/api/v1/users/%s/code-challenges/completed?page=%d", name, page)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
