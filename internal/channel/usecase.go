package channel

import (
	"context"
	"errors"

	"source.toby3d.me/toby3d/sub/internal/common"
	"source.toby3d.me/toby3d/sub/internal/domain"
)

type UseCase interface {
	Fetch(ctx context.Context, u domain.User) ([]domain.Channel, error)
	Create(ctx context.Context, u domain.User, name string) (*domain.Channel, error)
	Update(ctx context.Context, u domain.User, uid, name string) (*domain.Channel, error)
	Order(ctx context.Context, u domain.User, uids []string) error
	Delete(ctx context.Context, u domain.User, uid string) error
}

var ErrNotifications = errors.New(common.ChannelNotifications + " channel cannot be deleted")
