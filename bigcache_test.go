package go_cache

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func Test_bigCache(t *testing.T) {
	r := require.New(t)

	bc, err := newBigCache()
	r.NoError(err)

	// add a user in cache
	testUser := user{
		Id:    222,
		Email: "test@test.com",
	}
	err = bc.update(testUser)
	r.NoError(err)

	// try to read
	u, err := bc.read(testUser.Id)
	r.NoError(err)
	r.Equal(testUser, u)

	// update
	testUser.Id = 333
	err = bc.update(testUser)
	r.NoError(err)

	// read again
	u, err = bc.read(testUser.Id)
	r.NoError(err)
	r.Equal(testUser, u)

	// delete
	bc.delete(testUser.Id)

	// now there is no a user in cache
	u, err = bc.read(testUser.Id)
	r.EqualError(err, errUserNotInCache.Error())
	r.Equal(user{}, u)
}

func Benchmark_bigCache(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	bc, err := newBigCache()
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 10 updates
			bc.update(user{
				Id:    rand.Int63(),
				Email: "test",
			})

			// 10 reads
			bc.read(rand.Int63())

			// 10 deletes
			bc.delete(rand.Int63())
		}
	})
}
