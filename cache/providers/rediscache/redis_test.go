package rediscache_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"goeasy.dev/cache/providers/rediscache"
)

func TestRedisProvider(t *testing.T) {
	testString := "testing"
	testCases := []struct {
		desc       string
		value      interface{}
		dest       interface{}
		wantPutErr error
		wantErr    error
	}{
		{
			desc:  "HappyPath",
			value: testString,
			dest:  "",
		},
		{
			desc: "TestMap",
			value: map[string]interface{}{
				"testing": "a map",
			},
			dest: map[string]interface{}{},
		},
	}

	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{":6379"},
	})

	provider := rediscache.NewRedisProvider(client, time.Minute)
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := provider.Put(context.Background(), "test", tC.value)
			if tC.wantPutErr != nil && !assert.ErrorIs(t, err, tC.wantPutErr) {
				t.FailNow()
			} else if !assert.NoError(t, err) {
				t.FailNow()
			}

			err = provider.Get(context.Background(), "test", &tC.dest)
			if tC.wantErr != nil && !assert.ErrorIs(t, err, tC.wantErr) {
				t.FailNow()
			} else if !assert.NoError(t, err) {
				t.FailNow()
			}

			fmt.Println(tC.dest)

			assert.Equal(t, tC.value, tC.dest)
		})
	}
}
