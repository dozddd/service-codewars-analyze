//go:build integration

package get_unique_names

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser_GetAllUniqueNames(t *testing.T) {
	parser := New(
		KataURL,
		AuthoredURL,
		RanksURL,
		LeadersURL,
	)
	ctx := context.Background()

	names, err := parser.GetAllUniqueNames(ctx)

	require.NoError(t, err)
	require.Greater(t, len(names), 0)
}
