package get_unique_names

import (
	"context"
	"fmt"

	"github.com/gocolly/colly"
	"golang.org/x/sync/errgroup"
)

const (
	KataURL     codeWarsUrl = "https://www.codewars.com/users/leaderboard/kata"
	AuthoredURL codeWarsUrl = "https://www.codewars.com/users/leaderboard/authored"
	RanksURL    codeWarsUrl = "https://www.codewars.com/users/leaderboard/ranks"
	LeadersURL  codeWarsUrl = "https://www.codewars.com/users/leaderboard"
)

type codeWarsUrl string

type Parser struct {
	urls []codeWarsUrl
}

func New(url ...codeWarsUrl) *Parser {
	return &Parser{
		urls: url,
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

	for _, url := range p.urls {
		url := url
		g.Go(func() error {
			names, err := p.getNamesLeaders(ctx, string(url))
			if err != nil {
				return fmt.Errorf("p.getNamesLeaders: %w", err)
			}
			for _, name := range names {
				uniqueNames[name] = struct{}{}
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("g.Wait: %w", err)
	}

	var uniqueNamesSlice []string
	for name := range uniqueNames {
		uniqueNamesSlice = append(uniqueNamesSlice, name)
	}

	return uniqueNamesSlice, nil
}
