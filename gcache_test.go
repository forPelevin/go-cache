package go_cache

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func Test_gCache(t *testing.T) {
	r := require.New(t)

	gc := newGCache()

	// add a user in cache
	testUser := user{
		Id:    222,
		Email: "test@test.com",
	}
	err := gc.update(testUser, 1*time.Hour)
	r.NoError(err)

	// try to read
	u, err := gc.read(testUser.Id)
	r.NoError(err)
	r.Equal(testUser, u)

	// update
	testUser.Id = 333
	err = gc.update(testUser, 1*time.Hour)
	r.NoError(err)

	// read again
	u, err = gc.read(testUser.Id)
	r.NoError(err)
	r.Equal(testUser, u)

	// delete
	gc.delete(testUser.Id)

	// now there is no a user in cache
	u, err = gc.read(testUser.Id)
	r.EqualError(err, errUserNotInCache.Error())
	r.Equal(user{}, u)
}

func Benchmark_gCache(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	gc := newGCache()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 10 updates
			gc.update(user{
				Id:    rand.Int63(),
				Email: "test",
			}, 1*time.Hour)

			// 10 reads
			gc.read(rand.Int63())

			// 10 deletes
			gc.delete(rand.Int63())
		}
	})
}
