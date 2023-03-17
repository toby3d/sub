package usecase

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"

	"source.toby3d.me/toby3d/sub/internal/channel"
	"source.toby3d.me/toby3d/sub/internal/common"
	"source.toby3d.me/toby3d/sub/internal/domain"
)

type (
	channelUseCase struct {
		channels channel.Repository
	}

	orderItem struct {
		uid   string
		index int
	}
)

func NewChannelUseCase(channels channel.Repository) channel.UseCase {
	return &channelUseCase{
		channels: channels,
	}
}

func (ucase *channelUseCase) Fetch(ctx context.Context, u domain.User) ([]domain.Channel, error) {
	channels, err := ucase.channels.Fetch(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch channels: %w", err)
	}

	return append([]domain.Channel{{
		UID:    common.ChannelNotifications,
		Name:   "Notifications",
		Weight: -1,
	}}, channels...), nil
}

func (ucase *channelUseCase) Create(ctx context.Context, u domain.User, name string) (*domain.Channel, error) {
	id := make([]byte, 16)
	if _, err := rand.Read(id); err != nil {
		return nil, fmt.Errorf("cannot generate UID for new channel: %w", err)
	}

	uid := hex.EncodeToString(id)
	if err := ucase.channels.Create(ctx, u, domain.Channel{
		UID:  uid,
		Name: name,
	}); err != nil {
		return nil, fmt.Errorf("cannot create channel: %w", err)
	}

	out, err := ucase.channels.Get(ctx, u, uid)
	if err != nil {
		return nil, fmt.Errorf("cannot return created channel: %w", err)
	}

	return out, nil
}

func (ucase *channelUseCase) Update(ctx context.Context, u domain.User, uid, name string) (*domain.Channel, error) {
	if err := ucase.channels.Update(ctx, u, uid, func(tx *domain.Channel) (*domain.Channel, error) {
		tx.Name = name

		return tx, nil
	}); err != nil {
		return nil, fmt.Errorf("cannot update channel: %w", err)
	}

	out, err := ucase.channels.Get(ctx, u, uid)
	if err != nil {
		return nil, fmt.Errorf("cannot return updated channel: %w", err)
	}

	return out, nil
}

func (ucase *channelUseCase) Order(ctx context.Context, u domain.User, uids []string) error {
	channels, err := ucase.channels.Fetch(ctx, u)
	if err != nil {
		return fmt.Errorf("cannot fetch channels for ordering: %w", err)
	}

	result := make([]string, len(channels))
	buf := make([]orderItem, 0)
	for i := range channels {
		result[i] = channels[i].UID

		for j := range uids {
			if channels[i].UID != uids[j] {
				continue
			}

			buf = append(buf, orderItem{index: i, uid: channels[i].UID})
		}
	}

	for i := range uids {
		buf[i].uid = uids[i]
	}

	for i := range buf {
		result[buf[i].index] = buf[i].uid
	}

	for i := range result {
		if err = ucase.channels.Update(ctx, u, result[i], func(tx *domain.Channel) (*domain.Channel, error) {
			tx.Weight = len(result[:i])

			return tx, err
		}); err != nil {
			return fmt.Errorf("cannot update order: %w", err)
		}
	}

	return nil
}

func (ucase *channelUseCase) Delete(ctx context.Context, u domain.User, uid string) error {
	if uid == common.ChannelNotifications {
		return channel.ErrNotifications
	}

	if err := ucase.channels.Delete(ctx, u, uid); err != nil {
		return fmt.Errorf("cannot delete channel: %w", err)
	}

	return nil
}
