package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
)

func initRedis() error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Username: "",
		Password: "",
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		fmt.Println("connect to redis failed: ", err)
		return err
	} else {
		fmt.Println("connect to redis")
		return nil
	}
}

func TestSetValue(t *testing.T) {
	err := initRedis()
	if err != nil {
		return
	}
	ctx := context.Background()
	key := "name"
	value := "John Doe"

	err = redisClient.Set(ctx, key, value, 86400*time.Second).Err()
	if err != nil {
		t.Fatalf("set value failed: %v", err)
	} else {
		t.Log("set value success")
	}
}

func TestGetValue(t *testing.T) {
	err := initRedis()
	if err != nil {
		return
	}
	ctx := context.Background()
	value, err := redisClient.Get(ctx, "name").Result()
	if err != nil {
		t.Fatalf("get value failed: %v", err)
	} else {
		t.Log("get value success: ", value)
	}
}

func TestSetListValueAndGet(t *testing.T) {
	err := initRedis()
	if err != nil {
		return
	}
	ctx := context.Background()
	key := "zqf_hobbies"
	value := []interface{}{"basketball", "football", "photography"}
	err = redisClient.Expire(ctx, key, 86400*time.Second).Err()
	if err != nil {
		t.Fatalf("set expire failed: %v", err)
	}

	err = redisClient.RPush(ctx, key, value...).Err()
	if err != nil {
		t.Fatalf("set list value failed: %v", err)
	} else {
		t.Log("set list value success")
	}

	time.Sleep(5 * time.Second)

	var v []string
	v, err = redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		t.Fatalf("get list value failed: %v", err)
	} else {
		t.Log("get list value success: ", v)
	}

}
