package go_cache

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func Test_localCache(t *testing.T) {
	r := require.New(t)

	lc := newLocalCache(1 * time.Minute)

	// add a user in cache
	testUser := user{
		Id:    222,
		Email: "test@test.com",
	}
	lc.update(testUser, time.Now().Add(1*time.Hour).Unix())

	// try to read
	u, err := lc.read(testUser.Id)
	r.NoError(err)
	r.Equal(testUser, u)

	// update
	testUser.Id = 333
	lc.update(testUser, time.Now().Add(1*time.Hour).Unix())

	// read again
	u, err = lc.read(testUser.Id)
	r.NoError(err)
	r.Equal(testUser, u)

	// delete
	lc.delete(testUser.Id)

	// now there is no a user in cache
	u, err = lc.read(testUser.Id)
	r.EqualError(err, errUserNotInCache.Error())
	r.Equal(user{}, u)

	lc.stopCleanup()
}

func Benchmark_localCache(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	lc := newLocalCache(1 * time.Minute)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 10 updates
			lc.update(user{
				Id:    rand.Int63(),
				Email: "test",
			}, time.Now().Add(1*time.Hour).Unix())

			// 10 reads
			lc.read(rand.Int63())

			// 10 deletes
			lc.delete(rand.Int63())
		}
	})
}
