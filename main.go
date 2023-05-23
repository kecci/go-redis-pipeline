package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	totalKeys := 1
	ctx := context.Background()
	rdc := NewRedisClient()

	nowSet := time.Now()
	rdc.setCommand(ctx, totalKeys)
	durSet := time.Since(nowSet)
	defer log.Printf("processing %v keys time setCommand: %+v or %+v/key\n", totalKeys, durSet, durSet/time.Duration(totalKeys))

	nowSetPipe := time.Now()
	rdc.setCommandPipeline(ctx, totalKeys)
	durSetPipe := time.Since(nowSetPipe)
	defer log.Printf("processing %v keys time setCommandPipeline: %+v or %+v/key\n", totalKeys, durSetPipe, durSetPipe/time.Duration(totalKeys))

	nowDel := time.Now()
	rdc.delCommand(ctx, totalKeys)
	durDel := time.Since(nowDel)
	defer log.Printf("processing %v keys time delCommand: %+v or %+v/key\n", totalKeys, durDel, durDel/time.Duration(totalKeys))

	nowDelPipe := time.Now()
	rdc.delCommandPipeline(ctx, totalKeys)
	durDelPipe := time.Since(nowDelPipe)
	defer log.Printf("processing %v keys time delCommandPipeline: %+v or %+v/key\n", totalKeys, durDelPipe, durDelPipe/time.Duration(totalKeys))

}

type RedisClient struct {
	redisClient *redis.Client
	expireAt    time.Time
}

func NewRedisClient() RedisClient {
	cl := redis.NewClient(&redis.Options{
		Addr:            "localhost:6379",
		Password:        "",  // no password set
		DB:              0,   // use default DB
		PoolSize:        100, // max-active
		MaxIdleConns:    200,
		MinIdleConns:    100,
		ConnMaxIdleTime: time.Duration(300) * time.Second,
	})

	expireAt := time.Now().Add(time.Hour)

	return RedisClient{
		redisClient: cl,
		expireAt:    expireAt,
	}
}

func (r RedisClient) setCommand(ctx context.Context, totalKeys int) {
	for i := 1; i <= totalKeys; i++ {
		args := map[string]string{
			fmt.Sprintf("%d", i): fmt.Sprintf("%d", i),
		}

		key := fmt.Sprintf("redis-local:single:%d", i)
		r.redisClient.HMSet(ctx, key, args)
		r.redisClient.ExpireAt(ctx, key, r.expireAt)
	}
}

func (r RedisClient) setCommandPipeline(ctx context.Context, totalKeys int) {
	rdp := r.redisClient.Pipeline()

	for i := 1; i <= totalKeys; i++ {
		args := map[string]string{
			fmt.Sprintf("%d", i): fmt.Sprintf("%d", i),
		}

		key := fmt.Sprintf("redis-local:pipeline:%d", i)
		rdp.HMSet(ctx, key, args)
		rdp.ExpireAt(ctx, key, r.expireAt)
	}

	cmdr, err := rdp.Exec(ctx)
	if err != nil {
		errs := []error{}
		for _, v := range cmdr {
			errs = append(errs, v.Err())
		}
		panic(fmt.Errorf("%v - %v", err, errs))
	}
}

func (r RedisClient) delCommand(ctx context.Context, totalKeys int) {
	keys := []string{}
	for i := 1; i <= totalKeys; i++ {
		key := fmt.Sprintf("redis-local:single:%d", i)
		keys = append(keys, key)
	}
	r.redisClient.Del(ctx, keys...)
}

func (r RedisClient) delCommandPipeline(ctx context.Context, totalKeys int) {
	rdp := r.redisClient.Pipeline()
	keys := []string{}

	for i := 1; i <= totalKeys; i++ {
		key := fmt.Sprintf("redis-local:pipeline:%d", i)
		keys = append(keys, key)
	}

	rdp.Del(ctx, keys...)

	cmdr, err := rdp.Exec(ctx)
	if err != nil {
		errs := []error{}
		for _, v := range cmdr {
			errs = append(errs, v.Err())
		}
		panic(fmt.Errorf("%v - %v", err, errs))
	}
}
