package domain

import (
	"net/url"
	"testing"
)

type User struct {
	*url.URL
}

func TestUser(tb testing.TB) *User {
	tb.Helper()

	return &User{
		URL: &url.URL{
			Scheme: "https",
			Host:   "user.example.com",
			Path:   "/",
		},
	}
}
