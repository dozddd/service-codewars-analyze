package kata_ids

import (
	"context"
	"fmt"
	"sync"
)

type Usecase struct{}

func New() *Usecase {
	return &Usecase{}
}

func (u *Usecase) GetKataIds(ctx context.Context, usernames []string) ([]string, error) {
	jobs := make(chan string, len(usernames))
	results := make(chan []string, len(usernames))

	var wg sync.WaitGroup
	for w := 1; w <= 5; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			u.worker(jobs, results)
		}(w)
	}

	for _, username := range usernames {
		jobs <- username
	}
	close(jobs)

	kataIds := make(map[string]bool)
	go func() {
		wg.Wait()
		close(results)
	}()
	for r := range results {
		for _, kataId := range r {
			kataIds[kataId] = true
		}
	}

	kataIdSlice := make([]string, 0, len(kataIds))
	for kataId := range kataIds {
		kataIdSlice = append(kataIdSlice, kataId)
	}

	return kataIdSlice, nil
}

func (u *Usecase) worker(jobs <-chan string, results chan<- []string) {
	for username := range jobs {
		kataIds, err := u.getKataIdsForUser(username)
		if err != nil {
			fmt.Printf("getKataIdsForUser: %w", err)
			continue
		}
		results <- kataIds
	}
}

func (u *Usecase) getKataIdsForUser(username string) ([]string, error) {
	kataIds := make(map[string]bool)
	page := 0
	for {
		resp, err := SendGetRequest(username, page)
		if err != nil {
			return nil, fmt.Errorf("SendGetRequest: %w", err)
		}

		for _, kata := range resp.Data {
			kataIds[kata.ID] = true
		}
		if page >= resp.TotalPages-1 {
			break
		}
		page++
	}

	kataIdSlice := make([]string, 0, len(kataIds))
	for kataId := range kataIds {
		kataIdSlice = append(kataIdSlice, kataId)
	}

	return kataIdSlice, nil
}
