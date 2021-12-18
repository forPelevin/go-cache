package go_cache

import (
	"errors"
	"fmt"
	"github.com/bluele/gcache"
	"time"
)

type gCache struct {
	users gcache.Cache
}

const (
	cacheSize = 1_000_000
	cacheTTL  = 1 * time.Hour // default expiration
)

func newGCache() *gCache {
	return &gCache{
		users: gcache.New(cacheSize).Expiration(cacheTTL).ARC().Build(),
	}
}

func (gc *gCache) update(u user, expireIn time.Duration) error {
	return gc.users.SetWithExpire(u.Id, u, expireIn)
}

func (gc *gCache) read(id int64) (user, error) {
	val, err := gc.users.Get(id)
	if err != nil {
		if errors.Is(err, gcache.KeyNotFoundError) {
			return user{}, errUserNotInCache
		}

		return user{}, fmt.Errorf("get: %w", err)
	}

	return val.(user), nil
}

func (gc *gCache) delete(id int64) {
	gc.users.Remove(id)
}
