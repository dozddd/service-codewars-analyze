package html_parse

import "context"

type Parser struct {
}

func New() *Parser {
	return &Parser{}
}

// GetNamesLeaders возращает список уникальных ников из таблиц лидеров
func (p Parser) GetNamesLeaders(ctx context.Context) ([]string, error) {
	// TODO написать парсер
	return nil, nil
}
