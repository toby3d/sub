package memory

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"golang.org/x/exp/slices"

	"source.toby3d.me/toby3d/sub/internal/channel"
	"source.toby3d.me/toby3d/sub/internal/domain"
)

type memoryChannelRepository struct {
	mutex    *sync.RWMutex
	channels map[string][]domain.Channel
}

func NewMemoryChannelRepository() channel.Repository {
	return &memoryChannelRepository{
		mutex:    new(sync.RWMutex),
		channels: make(map[string][]domain.Channel, 0),
	}
}

func (repo *memoryChannelRepository) Create(ctx context.Context, u domain.User, c domain.Channel) error {
	result, err := repo.Get(ctx, u, c.UID)
	if err != nil && !errors.Is(err, channel.ErrNotExist) {
		return fmt.Errorf("cannot check creating channel: %w", err)
	}

	if result != nil {
		return channel.ErrExist
	}

	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	c.Weight = len(repo.channels[u.String()])
	repo.channels[u.String()] = append(repo.channels[u.String()], c)

	return nil
}

func (repo *memoryChannelRepository) Get(ctx context.Context, u domain.User, cid string) (*domain.Channel, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	for uid, channels := range repo.channels {
		if uid != u.String() {
			continue
		}

		for j := range channels {
			if channels[j].UID != cid {
				continue
			}

			return &channels[j], nil
		}

		break
	}

	return nil, channel.ErrNotExist
}

func (repo *memoryChannelRepository) Fetch(ctx context.Context, u domain.User) ([]domain.Channel, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	if out, ok := repo.channels[u.String()]; ok {
		sort.Slice(out, func(i, j int) bool {
			return out[i].Weight < out[j].Weight
		})

		return out, nil
	}

	return nil, channel.ErrNotExist
}

func (repo *memoryChannelRepository) Update(ctx context.Context, u domain.User, cid string, update channel.UpdateFunc) error {
	in, err := repo.Get(ctx, u, cid)
	if err != nil {
		return fmt.Errorf("cannot find updating channel: %w", err)
	}

	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	out, err := update(in)
	if err != nil {
		return fmt.Errorf("cannot update channel: %w", err)
	}

	channels := repo.channels[u.String()]

	for i := range channels {
		if channels[i].UID != cid {
			continue
		}

		repo.channels[u.String()][i] = *out
	}

	return nil
}

func (repo *memoryChannelRepository) Delete(ctx context.Context, u domain.User, cid string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for uid, channels := range repo.channels {
		if uid != u.String() {
			continue
		}

		for j := range channels {
			repo.channels[uid] = slices.Delete(channels, j, j+1)

			break
		}

		break
	}

	return nil
}
