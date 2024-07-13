package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/gocolly/colly"
	"golang.org/x/sync/errgroup"
)

const (
	kataURL     = "https://www.codewars.com/users/leaderboard/kata"
	authoredURL = "https://www.codewars.com/users/leaderboard/authored"
	ranksURL    = "https://www.codewars.com/users/leaderboard/ranks"
	leadersURL  = "https://www.codewars.com/users/leaderboard"
)

func main() {
	parser := New()
	ctx := context.Background()

	names, err := parser.GetAllUniqueNames(ctx)
	if err != nil {
		log.Fatalf("Error parsing names from leaders: %v", err)
	}
	fmt.Println("Уникальные ники из всех таблиц:", names)
}

type Parser struct {
	urls []string
}

func New() *Parser {
	return &Parser{
		urls: []string{kataURL, authoredURL, ranksURL, leadersURL},
	}
}

func (p *Parser) getNamesLeaders(_ context.Context, url string) ([]string, error) {
	c := colly.NewCollector()

	names := make(map[string]struct{})

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		names[e.ChildText("td a")] = struct{}{}
	})

	if err := c.Visit(url); err != nil {
		return nil, fmt.Errorf("c.Visit: %w", err)
	}

	var uniqueNamesSlice []string
	for name := range names {
		uniqueNamesSlice = append(uniqueNamesSlice, name)
	}

	return uniqueNamesSlice, nil
}

func (p *Parser) GetAllUniqueNames(ctx context.Context) ([]string, error) {
	var g errgroup.Group
	uniqueNames := make(map[string]struct{})
	mt := sync.Mutex{}

	for _, url := range p.urls {
		url := url
		g.Go(func() error {
			names, err := p.getNamesLeaders(ctx, url)
			if err != nil {
				return err
			}
			mt.Lock()
			for _, name := range names {
				uniqueNames[name] = struct{}{}
			}
			mt.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	var uniqueNamesSlice []string
	for name := range uniqueNames {
		uniqueNamesSlice = append(uniqueNamesSlice, name)
	}

	return uniqueNamesSlice, nil
}
