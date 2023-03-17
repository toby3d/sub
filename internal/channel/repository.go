package channel

import (
	"context"
	"errors"

	"source.toby3d.me/toby3d/sub/internal/domain"
)

type (
	UpdateFunc func(channel *domain.Channel) (*domain.Channel, error)

	Repository interface {
		Create(ctx context.Context, user domain.User, channel domain.Channel) error
		Get(ctx context.Context, user domain.User, uid string) (*domain.Channel, error)
		Fetch(ctx context.Context, user domain.User) ([]domain.Channel, error)
		Update(ctx context.Context, user domain.User, uid string, update UpdateFunc) error
		Delete(ctx context.Context, user domain.User, uid string) error
	}
)

var (
	ErrNotExist = errors.New("channel does not exist")
	ErrExist    = errors.New("channel already exists")
)
