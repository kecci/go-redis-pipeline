package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	totalKey = 10000
)

func main() {
	ctx := context.Background()

	rdc := NewRedisClient()

	nowSingle := time.Now()
	rdc.singleCommandRedis(ctx)
	durSingle := time.Since(nowSingle)

	nowPipe := time.Now()
	rdc.pipelineRedis(ctx)
	durPipe := time.Since(nowPipe)

	nowSinglePipe := time.Now()
	rdc.singlePipelineRedis(ctx)
	durSinglePipe := time.Since(nowSinglePipe)

	log.Printf("processing %v keys time singleCommandRedis: %+v or %+v/key\n", totalKey, durSingle, time.Duration(durSingle.Nanoseconds()/int64(totalKey)))
	log.Printf("processing %v keys time pipelineRedis: %+v or %+v/key\n", totalKey, durPipe, time.Duration(durPipe.Nanoseconds()/int64(totalKey)))
	log.Printf("processing %v keys time singlePipelineRedis: %+v or %+v/key\n", 1, durSinglePipe, time.Duration(durSinglePipe.Nanoseconds()/1))
}

type RedisClient struct {
	redisClient *redis.Client
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
	return RedisClient{
		redisClient: cl,
	}
}

func (r RedisClient) singleCommandRedis(ctx context.Context) {
	for i := 1; i <= totalKey; i++ {
		err := r.redisClient.HMSet(ctx, "redis-local:single", fmt.Sprintf("single_proccess_%d", i), i)
		if err != nil {
			log.Printf("%+v\n", err)
			continue
		}
	}
}

func (r RedisClient) pipelineRedis(ctx context.Context) {
	rdp := r.redisClient.Pipeline()

	redisData := make(map[string]string)
	for i := 1; i <= totalKey; i++ {
		key := fmt.Sprintf("testing_key_redis_%d", i)
		redisData[key] = fmt.Sprintf("%d", i)
	}

	err := rdp.HMSet(ctx, "redis-local:pipeline", redisData).Err()
	if err != nil {
		log.Printf("%+v\n", err)
	}

	_, err = rdp.Exec(ctx)
	if err != nil {
		log.Printf("%+v\n", err)
	}
}

func (r RedisClient) singlePipelineRedis(ctx context.Context) {
	rdp := r.redisClient.Pipeline()

	redisData := make(map[string]string)
	key := fmt.Sprintf("testing_key_redis_%d", 1)
	redisData[key] = fmt.Sprintf("%d", 1)

	err := rdp.HMSet(ctx, "redis-local:single-pipeline", redisData).Err()
	if err != nil {
		log.Printf("%+v\n", err)
	}

	_, err = rdp.Exec(ctx)
	if err != nil {
		log.Printf("%+v\n", err)
	}
}
