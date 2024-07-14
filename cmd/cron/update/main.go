package main

import (
	"context"
	"fmt"
	"log"

	get_unique_names "github.com/dozddd/service-codewars-analyze.git/internal/usecase/get_uique_names"
)

func main() {
	parser := get_unique_names.New(
		get_unique_names.KataURL,
		get_unique_names.AuthoredURL,
		get_unique_names.RanksURL,
		get_unique_names.LeadersURL,
	)
	ctx := context.Background()

	names, err := parser.GetAllUniqueNames(ctx)
	if err != nil {
		log.Fatalf("Error parsing names from leaders: %v", err)
	}
	fmt.Println("Уникальные ники из всех таблиц:", names)
}
