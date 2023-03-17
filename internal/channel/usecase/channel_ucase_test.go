package usecase_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	channelmemoryrepo "source.toby3d.me/toby3d/sub/internal/channel/repository/memory"
	ucase "source.toby3d.me/toby3d/sub/internal/channel/usecase"
	"source.toby3d.me/toby3d/sub/internal/common"
	"source.toby3d.me/toby3d/sub/internal/domain"
)

func TestChannelUseCase_Create(t *testing.T) {
	t.Parallel()

	user := domain.TestUser(t)
	channels := channelmemoryrepo.NewMemoryChannelRepository()

	actual, err := ucase.NewChannelUseCase(channels).Create(context.Background(), *user, "Testing")
	if err != nil {
		t.Fatal(err)
	}

	if actual.UID == "" {
		t.Error("expect non empty UID, got empty")
	}

	if actual.Name != "Testing" {
		t.Errorf("expect %s, got %s", "Testing", actual.Name)
	}
}

func TestChannelUseCase_Update(t *testing.T) {
	t.Parallel()

	user := domain.TestUser(t)
	channels := channelmemoryrepo.NewMemoryChannelRepository()

	channel := domain.TestChannel(t)
	if err := channels.Create(context.Background(), *user, *channel); err != nil {
		t.Fatal(err)
	}

	actual, err := ucase.NewChannelUseCase(channels).
		Update(context.Background(), *user, channel.UID, "Testing")
	if err != nil {
		t.Fatal(err)
	}

	if actual.UID == "" {
		t.Error("expect non empty UID, got empty")
	}

	if actual.Name != "Testing" {
		t.Errorf("expect %s, got %s", "Testing", actual.Name)
	}
}

func TestChannelUseCase_Order(t *testing.T) {
	t.Parallel()

	user := domain.TestUser(t)
	channels := channelmemoryrepo.NewMemoryChannelRepository()

	for _, n := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
		if err := channels.Create(context.Background(), *user, domain.Channel{UID: n, Name: n}); err != nil {
			t.Fatal(err)
		}
	}

	if err := ucase.NewChannelUseCase(channels).
		Order(context.Background(), *user, []string{"d", "a", "c", "g"}); err != nil {
		t.Fatal(err)
	}

	result, err := channels.Fetch(context.Background(), *user)
	if err != nil {
		t.Fatal(err)
	}

	expect := []string{"d", "b", "a", "c", "e", "f", "g", "h"}
	actual := make([]string, len(expect))

	for i := range result {
		actual[i] = result[i].UID
	}

	if diff := cmp.Diff(expect, actual); diff != "" {
		t.Error(diff)
	}
}

func TestChannelUseCase_Fetch(t *testing.T) {
	t.Parallel()

	user := domain.TestUser(t)
	channels := channelmemoryrepo.NewMemoryChannelRepository()
	expect := []domain.Channel{
		{UID: common.ChannelNotifications, Name: "Notifications", Weight: -1},
		*domain.TestChannel(t),
		*domain.TestChannel(t),
	}
	expect[2].Weight = 1

	for _, c := range expect[1:] {
		if err := channels.Create(context.Background(), *user, c); err != nil {
			t.Fatal(err)
		}
	}

	actual, err := ucase.NewChannelUseCase(channels).Fetch(context.Background(), *user)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(actual, expect); diff != "" {
		t.Error(diff)
	}
}

func TestChannelUseCase_Delete(t *testing.T) {
	t.Parallel()

	user := domain.TestUser(t)
	channels := channelmemoryrepo.NewMemoryChannelRepository()

	channel := domain.TestChannel(t)
	if err := channels.Create(context.Background(), *user, *channel); err != nil {
		t.Fatal(err)
	}

	if err := ucase.NewChannelUseCase(channels).
		Delete(context.Background(), *user, channel.UID); err != nil {
		t.Fatal(err)
	}
}
