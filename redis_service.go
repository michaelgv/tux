package main

import (
	"github.com/go-redis/redis"
	"fmt"
	"time"
)

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	return client
}

func FlushRedisDB() {
	client := NewRedisClient()
	client.FlushDB()
}

func AccountListSetCache(values string) {
	client := NewRedisClient()
	err := client.Set("accountlist", values, 300 * time.Second).Err()
	checkErr(err)
}

func AccountListGetCache() (string, bool) {
	client := NewRedisClient()
	val, err := client.Get("accountlist").Result()
	if err != nil {
		return "", false
	}
	return val, true
}

func AccountListActiveGetCache() (string, bool) {
	client := NewRedisClient()
	val, err := client.Get("accountlistactive").Result()
	if err != nil {
		return "", false
	}
	return val, true
}

func AccountListActiveSetCache(values string) (string, bool) {
	client := NewRedisClient()
	err := client.Set("accountlistactive", values, 300 * time.Second).Err()
	checkErr(err)
	return values, true
}

func AccountListNewGetCache() (string, bool) {
	client := NewRedisClient()
	val, err := client.Get("accountlistnew").Result()
	if err != nil {
		return "", false
	}
	return val, true
}

func AccountListNewSetCache(values string) (string, bool) {
	client := NewRedisClient()
	err := client.Set("accountlistnew", values, 300 * time.Second).Err()
	checkErr(err)
	return values, true
}

func GenericRedisGet(key string) (string, bool) {
	client := NewRedisClient()
	val, err := client.Get("accountlist").Result()
	if err != nil {
		return "", false
	}
	return val, true
}

func GenericRedisSet(key string, values string) bool {
	client := NewRedisClient()
	err := client.Set(key, values, 300 * time.Second).Err()
	if err != nil {
		return false
	}
	return true
}

func GenericRedisUpdate(key, newValue string) (bool, error) {
	if GenericRedisDelete(key) {
		success := GenericRedisSet(key, newValue)
		if success {
			return true, nil
		}
	}
	return false, fmt.Errorf("unknown_error_for_key=%s", key)
}

func GenericRedisDelete(key string) (bool) {
	client := NewRedisClient()
	client.Del(key)
	return true
}