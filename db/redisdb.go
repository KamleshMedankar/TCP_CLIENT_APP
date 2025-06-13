package db

import (
	"clientgo/models"
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"time"
	//"context"
)

var (
	Rdb *redis.Client
	//ctx := context.Background()
)

func ConnectRedisClient() {
	
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", 
		Password: "",               
		DB:       0,                // use default DB 0
	})

	
	_, err := Rdb.Ping().Result()
	if err != nil {
		log.Printf("Could not connect to Redis: %v\n", err)
		return
	}
	log.Println("Connected to Redis successfully")
}

func RedisSet(key string, data models.TenantData, expiration time.Duration) error {
	//ctx := context.Background()

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return Rdb.Set(key, jsonBytes, expiration).Err()
}

func RedisGet(key string) (models.TenantData, error) {
	var result models.TenantData

	val, err := Rdb.Get(key).Result()
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func CloseRedis() {
	err := Rdb.Close()
	if err != nil {
		log.Println("Error closing connection Redis: ", err)
		return
	}
	log.Println("Closed Redis connection successfully!")
}
