package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	channelmemoryrepo "source.toby3d.me/toby3d/sub/internal/channel/repository/memory"
	channelucase "source.toby3d.me/toby3d/sub/internal/channel/usecase"
	"source.toby3d.me/toby3d/sub/internal/common"
	"source.toby3d.me/toby3d/sub/internal/domain"
	delivery "source.toby3d.me/toby3d/sub/internal/microsub/delivery/http"
)

var update = flag.Bool("update", false, "update golden files")

func TestHandler_ServeHTTP_ChannelsCreate(t *testing.T) {
	t.Parallel()

	user := domain.TestUser(t)

	q := make(url.Values)
	q.Set("action", domain.ActionChannels.String())
	q.Set("name", "Testing")

	req := httptest.NewRequest(http.MethodPost, "https://example.com/", strings.NewReader(q.Encode()))
	req.Header.Set(common.HeaderContentType, common.MIMEApplicationFormCharsetUTF8)
	req = req.WithContext(context.WithValue(req.Context(), "user", user))

	w := httptest.NewRecorder()
	delivery.NewHandler(channelucase.NewChannelUseCase(channelmemoryrepo.NewMemoryChannelRepository())).
		ServeHTTP(w, req)

	resp := w.Result()
	if expect := http.StatusOK; resp.StatusCode != expect {
		t.Errorf("want %d, got %d", expect, resp.StatusCode)
	}

	actual := new(delivery.ResponseChannel)
	if err := json.NewDecoder(resp.Body).Decode(actual); err != nil {
		t.Fatal(err)
	}

	if expect := q.Get("name"); actual.Name != expect {
		t.Errorf("want '%s', got '%s'", expect, actual.Name)
	}
}

func TestHandler_ServeHTTP_ChannelsFetch(t *testing.T) {
	t.Parallel()

	user := domain.TestUser(t)
	req := httptest.NewRequest(http.MethodGet, "https://example.com/?action=channels", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user", user))

	channels := channelmemoryrepo.NewMemoryChannelRepository()

	for _, c := range [][2]string{
		{"31eccfe322d6c48c50dea2c84efc74ff", "IndieWeb"},
		{"1870e67e924856dc7e4c37732b303b45", "W3C"},
	} {
		if err := channels.Create(context.Background(), *user, domain.Channel{
			UID:  c[0],
			Name: c[1],
		}); err != nil {
			t.Fatal(err)
		}
	}

	w := httptest.NewRecorder()
	delivery.NewHandler(channelucase.NewChannelUseCase(channels)).
		ServeHTTP(w, req)

	resp := w.Result()
	if expect := http.StatusOK; resp.StatusCode != expect {
		t.Errorf("want %d, got %d", expect, resp.StatusCode)
	}

	actual, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	golden := filepath.Join("testdata", t.Name()+".golden")
	if *update {
		if err = ioutil.WriteFile(golden, actual, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	expected, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(actual, expected) {
		t.Error(cmp.Diff(actual, expected))
	}
}

func TestHandler_ServeHTTP_ChannelsDelete(t *testing.T) {
	t.Parallel()

	user := domain.TestUser(t)
	channel := domain.TestChannel(t)

	q := make(url.Values)
	q.Set("action", domain.ActionChannels.String())
	q.Set("method", "delete")
	q.Set("channel", channel.UID)

	req := httptest.NewRequest(http.MethodPost, "https://example.com/", strings.NewReader(q.Encode()))
	req.Header.Set(common.HeaderContentType, common.MIMEApplicationFormCharsetUTF8)
	req = req.WithContext(context.WithValue(req.Context(), "user", user))

	channels := channelmemoryrepo.NewMemoryChannelRepository()

	if err := channels.Create(context.Background(), *user, *channel); err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	delivery.NewHandler(channelucase.NewChannelUseCase(channels)).
		ServeHTTP(w, req)

	resp := w.Result()
	if expect := http.StatusNoContent; resp.StatusCode != expect {
		t.Errorf("want %d, got %d", expect, resp.StatusCode)
	}

	if !t.Failed() {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(body))
}

func TestHandler_ServeHTTP_ChannelsUpdate(t *testing.T) {
	t.Parallel()

	user := domain.TestUser(t)
	channel := domain.TestChannel(t)

	q := make(url.Values)
	q.Set("action", domain.ActionChannels.String())
	q.Set("channel", channel.UID)
	q.Set("name", "Testing")

	req := httptest.NewRequest(http.MethodPost, "https://example.com/", strings.NewReader(q.Encode()))
	req.Header.Set(common.HeaderContentType, common.MIMEApplicationFormCharsetUTF8)
	req = req.WithContext(context.WithValue(req.Context(), "user", user))

	channels := channelmemoryrepo.NewMemoryChannelRepository()

	if err := channels.Create(context.Background(), *user, *channel); err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	delivery.NewHandler(channelucase.NewChannelUseCase(channels)).
		ServeHTTP(w, req)

	resp := w.Result()
	if expect := http.StatusOK; resp.StatusCode != expect {
		t.Errorf("want %d, got %d", expect, resp.StatusCode)
	}

	actual, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	golden := filepath.Join("testdata", t.Name()+".golden")
	if *update {
		if err = ioutil.WriteFile(golden, actual, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	expected, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(actual, expected) {
		t.Error(cmp.Diff(actual, expected))
	}
}
