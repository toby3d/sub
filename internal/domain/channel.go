package domain

import (
	"encoding/hex"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"golang.org/x/exp/rand"

	"source.toby3d.me/toby3d/sub/internal/common"
)

type Channel struct {
	UID    string
	Name   string
	Weight int
}

func TestChannel(tb testing.TB) *Channel {
	tb.Helper()

	uid := make([]byte, 8+rand.Intn(16))
	if _, err := rand.Read(uid); err != nil {
		tb.Fatal(err)
	}

	return &Channel{
		UID:  hex.EncodeToString(uid),
		Name: gofakeit.InputName(),
	}
}

func (c Channel) IsNotifications() bool { return c.UID == common.ChannelNotifications }

func (c Channel) IsGlobal() bool { return c.UID == common.ChannelGlobal }
