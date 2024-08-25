// package kata_ids

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"time"
// )

// type Usecase struct{}

// func New() *Usecase {
// 	return &Usecase{}
// }

// type CodewarsResponse struct {
// 	TotalPages int `json:"totalPages"`
// 	TotalItems int `json:"totalItems"`
// 	Data       []struct {
// 		ID                 string    `json:"id"`
// 		Name               string    `json:"name"`
// 		Slug               string    `json:"slug"`
// 		CompletedAt        time.Time `json:"completedAt"`
// 		CompletedLanguages []string  `json:"completedLanguages"`
// 	} `json:"data"`
// }

// func (u *Usecase) GetKataIds(ctx context.Context, usernames []string) ([]string, error) {
// 	kataIds := make(map[string]bool)
// 	for _, username := range usernames {
// 		page := 0
// 		for {
// 			resp, err := SendGetRequest(username, page)
// 			if err != nil {
// 				return nil, fmt.Errorf("SendGetRequest: %w", err)
// 			}
// 			var codewarsResponse CodewarsResponse
// 			err = json.Unmarshal(resp, &codewarsResponse)
// 			if err != nil {
// 				return nil, fmt.Errorf("Unmarshal: %w", err)
// 			}
// 			for _, kata := range codewarsResponse.Data {
// 				kataIds[kata.ID] = true
// 			}
// 			if page >= codewarsResponse.TotalPages-1 {
// 				break
// 			}
// 			page++
// 		}
// 	}
// 	kataIdSlice := make([]string, 0, len(kataIds))
// 	for kataId := range kataIds {
// 		kataIdSlice = append(kataIdSlice, kataId)
// 	}
// 	return kataIdSlice, nil
// }

package kata_ids

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

type Usecase struct{}

func New() *Usecase {
	return &Usecase{}
}

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

func (u *Usecase) GetKataIds(ctx context.Context, usernames []string) ([]string, error) {
	kataIds := make(chan string)
	var wg sync.WaitGroup

	for _, username := range usernames {
		wg.Add(1)
		go func(username string) {
			defer wg.Done()
			page := 0
			for {
				resp, err := SendGetRequest(username, page)
				if err != nil {
					return
				}
				var codewarsResponse CodewarsResponse
				err = json.Unmarshal(resp, &codewarsResponse)
				if err != nil {
					return
				}
				for _, kata := range codewarsResponse.Data {
					kataIds <- kata.ID
				}
				if page >= codewarsResponse.TotalPages-1 {
					break
				}
				page++
			}
		}(username)
	}

	go func() {
		wg.Wait()
		close(kataIds)
	}()

	kataIdSlice := make([]string, 0, len(usernames))
	for kataId := range kataIds {
		kataIdSlice = append(kataIdSlice, kataId)
	}
	return kataIdSlice, nil
}
