package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetNamesLeaders(t *testing.T) {
	tCase := []struct {
		name           string
		url            string
		response       string
		expectedOutput []string
	}{
		{
			name:           "Completed Kata Leaderboard",
			url:            "https://www.codewars.com/users/leaderboard/kata",
			response:       `<table><tbody><tr><td><a href="/users/user1">user1</a></td></tr><tr><td><a href="/users/user2">user2</a></td></tr><tr><td><a href="/users/user3">user3</a></td></tr></tbody></table>`,
			expectedOutput: []string{"user1", "user2", "user3"},
		},
		{
			name:           "Authored Kata Leaderboard",
			url:            "https://www.codewars.com/users/leaderboard/authored",
			response:       `<table><tbody><tr><td><a href="/users/user4">user4</a></td></tr><tr><td><a href="/users/user5">user5</a></td></tr><tr><td><a href="/users/user6">user6</a></td></tr></tbody></table>`,
			expectedOutput: []string{"user4", "user5", "user6"},
		},
		{
			name:           "Ranks Leaderboard",
			url:            "https://www.codewars.com/users/leaderboard/ranks",
			response:       `<table><tbody><tr><td><a href="/users/user7">user7</a></td></tr><tr><td><a href="/users/user8">user8</a></td></tr><tr><td><a href="/users/user9">user9</a></td></tr></tbody></table>`,
			expectedOutput: []string{"user7", "user8", "user9"},
		},
		{
			name:           "Overall Leaderboard",
			url:            "https://www.codewars.com/users/leaderboard",
			response:       `<table><tbody><tr><td><a href="/users/user10">user10</a></td></tr><tr><td><a href="/users/user11">user11</a></td></tr><tr><td><a href="/users/user12">user12</a></td></tr></tbody></table>`,
			expectedOutput: []string{"user10", "user11", "user12"},
		},
	}

	for _, tCase := range tCase {
		t.Run(tCase.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(tCase.response))
			}))
			defer server.Close()

			p := New()
			actual, err := p.GetNamesLeaders(context.Background(), server.URL)
			if err != nil {
				t.Fail()
			}

			if !reflect.DeepEqual(actual, tCase.expectedOutput) {
				t.Fail()
			}
		})
	}
}
