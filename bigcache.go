package go_cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/allegro/bigcache"
	"strconv"
	"time"
)

type bigCache struct {
	users *bigcache.BigCache
}

func newBigCache() (*bigCache, error) {
	bCache, err := bigcache.NewBigCache(bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,

		// time after which entry can be evicted
		LifeWindow: 1 * time.Hour,

		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		CleanWindow: 5 * time.Minute,

		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: 1000 * 10 * 60,

		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: 500,

		// prints information about additional memory allocation
		Verbose: false,

		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: 256,

		// callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		OnRemove: nil,

		// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A constant representing the reason will be passed through.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// Ignored if OnRemove is specified.
		OnRemoveWithReason: nil,
	})
	if err != nil {
		return nil, fmt.Errorf("new big cache: %w", err)
	}

	return &bigCache{
		users: bCache,
	}, nil
}

func (bc *bigCache) update(u user) error {
	bs, err := json.Marshal(&u)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	return bc.users.Set(userKey(u.Id), bs)
}

func userKey(id int64) string {
	return strconv.FormatInt(id, 10)
}

func (bc *bigCache) read(id int64) (user, error) {
	bs, err := bc.users.Get(userKey(id))
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return user{}, errUserNotInCache
		}

		return user{}, fmt.Errorf("get: %w", err)
	}

	var u user
	err = json.Unmarshal(bs, &u)
	if err != nil {
		return user{}, fmt.Errorf("unmarshal: %w", err)
	}

	return u, nil
}

func (bc *bigCache) delete(id int64) {
	bc.users.Delete(userKey(id))
}
