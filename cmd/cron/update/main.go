package main

import (
	"context"
	"fmt"
	"log"

	get_unique_names "github.com/dozddd/service-codewars-analyze/internal/usecase/get_unique_names"
	kata_ids "github.com/dozddd/service-codewars-analyze/internal/usecase/kata_ids"
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

	resp, err := kata_ids.New().GetKataIds(ctx, names)
	if err != nil {
		log.Fatalf("GetKataIds: %v", err)
	}

	fmt.Println("Уникальные каты:", resp)
}
