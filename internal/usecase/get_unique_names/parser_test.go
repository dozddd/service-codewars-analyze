package get_unique_names

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNamesLeaders(t *testing.T) {
	tCases := []struct {
		name        string
		htmlFile    string
		expected    []string
		expectedErr error
	}{
		{
			name:        "Success",
			htmlFile:    "test_data/success.html",
			expected:    []string{"user1", "user2", "user3"},
			expectedErr: nil,
		},
		{
			name:        "Duplicates",
			htmlFile:    "test_data/duplicates.html",
			expected:    []string{"user1", "user2"},
			expectedErr: nil,
		},
		{
			name:        "NoUsernames",
			htmlFile:    "test_data/nousernames.html",
			expected:    []string{},
			expectedErr: nil,
		},
		{
			name:        "InvalidHTML",
			htmlFile:    "test_data/invalid.html",
			expected:    []string{},
			expectedErr: nil,
		},
	}

	for _, tCase := range tCases {
		t.Run(tCase.name, func(t *testing.T) {
			htmlData, err := os.ReadFile(tCase.htmlFile)
			assert.NoError(t, err)

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write(htmlData)
			}))
			defer server.Close()

			parser := &Parser{urls: []codeWarsUrl{codeWarsUrl(server.URL)}}
			names, err := parser.getNamesLeaders(context.Background(), server.URL)

			if tCase.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tCase.expected, names)
			}
		})
	}
}

func TestGetAllUniqueNames(t *testing.T) {
	tCases := []struct {
		name        string
		htmlFiles   []string
		expected    []string
		expectedErr error
	}{
		{
			name:        "Success",
			htmlFiles:   []string{"test_data/success.html"},
			expected:    []string{"user1", "user2", "user3"},
			expectedErr: nil,
		},
		{
			name:        "Duplicates",
			htmlFiles:   []string{"test_data/duplicates.html"},
			expected:    []string{"user1", "user2"},
			expectedErr: nil,
		},
		{
			name:        "NoUsernames",
			htmlFiles:   []string{"test_data/nousernames.html"},
			expected:    []string{},
			expectedErr: nil,
		},
		{
			name:        "InvalidHTML",
			htmlFiles:   []string{"test_data/invalid.html"},
			expected:    []string{},
			expectedErr: nil,
		},
		{
			name:        "MultipleFiles",
			htmlFiles:   []string{"test_data/success.html", "test_data/duplicates.html"},
			expected:    []string{"user1", "user2", "user3"},
			expectedErr: nil,
		},
	}

	for _, tCase := range tCases {
		t.Run(tCase.name, func(t *testing.T) {
			var servers []*httptest.Server
			for _, htmlFile := range tCase.htmlFiles {
				htmlData, err := os.ReadFile(htmlFile)
				assert.NoError(t, err)

				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write(htmlData)
				}))
				defer server.Close()
				servers = append(servers, server)
			}

			var urls []codeWarsUrl
			for _, server := range servers {
				urls = append(urls, codeWarsUrl(server.URL))
			}

			parser := &Parser{urls: urls}
			names, err := parser.GetAllUniqueNames(context.Background())

			if tCase.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tCase.expected, names)
			}
		})
	}
}
