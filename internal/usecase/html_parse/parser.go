package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

func main() {
	parser := New()
	ctx := context.Background()
	namesCompleted, err := parser.GetNamesLeaders(ctx, "https://www.codewars.com/users/leaderboard/kata")
	if err != nil {
		log.Fatalf("Error parsing names from leaders: %v", err)
	}
	fmt.Println("Уникальные ники из таблицы выполненных кат:", namesCompleted)

	namesAuthored, err := parser.GetNamesLeaders(ctx, "https://www.codewars.com/users/leaderboard/authored")
	if err != nil {
		log.Fatalf("Error parsing names from authored: %v", err)
	}
	fmt.Println("Уникальные ники из таблицы созданных кат:", namesAuthored)
	namesRanks, err := parser.GetNamesLeaders(ctx, "https://www.codewars.com/users/leaderboard/ranks")
	if err != nil {
		log.Fatalf("Error parsing names from ranks: %v", err)
	}
	fmt.Println("Уникальные ники из таблицы рангов:", namesRanks)
	namesLeaders, err := parser.GetNamesLeaders(ctx, "https://www.codewars.com/users/leaderboard")
	if err != nil {
		log.Fatalf("Error parsing names from ranks: %v", err)
	}
	fmt.Println("Уникальные ники из таблицы лидеров:", namesLeaders)
}

type Parser struct {
}

func New() *Parser {
	return &Parser{}
}

func (p Parser) GetNamesLeaders(ctx context.Context, url string) ([]string, error) {
	c := colly.NewCollector()

	var names []string

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		name := e.ChildText("td a")
		names = append(names, name)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("c.Visit: %w", err)
	}

	uniqueNames := make(map[string]struct{})
	var uniqueNamesSlice []string
	for _, name := range names {
		if _, exist := uniqueNames[name]; !exist {
			uniqueNames[name] = struct{}{}
			uniqueNamesSlice = append(uniqueNamesSlice, name)
		}
	}
	return uniqueNamesSlice, nil
	// return names, err
}
