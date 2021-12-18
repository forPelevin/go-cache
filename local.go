package go_cache

import (
	"errors"
	"sync"
	"time"
)

type user struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
}

type cachedUser struct {
	user
	expireAtTimestamp int64
}

type localCache struct {
	stop chan struct{}

	wg    sync.WaitGroup
	mu    sync.RWMutex
	users map[int64]cachedUser
}

func newLocalCache(cleanupInterval time.Duration) *localCache {
	lc := &localCache{
		users: make(map[int64]cachedUser),
		stop:  make(chan struct{}),
	}

	lc.wg.Add(1)
	go func(cleanupInterval time.Duration) {
		defer lc.wg.Done()
		lc.cleanupLoop(cleanupInterval)
	}(cleanupInterval)

	return lc
}

func (lc *localCache) cleanupLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-lc.stop:
			return
		case <-t.C:
			lc.mu.Lock()
			for uid, cu := range lc.users {
				if cu.expireAtTimestamp <= time.Now().Unix() {
					delete(lc.users, uid)
				}
			}
			lc.mu.Unlock()
		}
	}
}

func (lc *localCache) stopCleanup() {
	close(lc.stop)
	lc.wg.Wait()
}

func (lc *localCache) update(u user, expireAtTimestamp int64) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.users[u.Id] = cachedUser{
		user:              u,
		expireAtTimestamp: expireAtTimestamp,
	}
}

var (
	errUserNotInCache = errors.New("the user isn't in cache")
)

func (lc *localCache) read(id int64) (user, error) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	cu, ok := lc.users[id]
	if !ok {
		return user{}, errUserNotInCache
	}

	return cu.user, nil
}

func (lc *localCache) delete(id int64) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	delete(lc.users, id)
}
