package redisdb

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/emersonary/go-authentication/config"
	"github.com/emersonary/go-authentication/pck"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func RedisDB(config *config.Conf) (*redis.Client, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + strconv.Itoa(config.RedisPort), // use default Addr
		Password: "Redis2019!",
		DB:       0, // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis on " + config.RedisHost + ":" + strconv.Itoa(config.RedisPort))

	return rdb, nil

}

func TestPush(redisClient *redis.Client) {

	id := pck.NewID()

	messages := []string{"1 Hello", "2 World", "3 This", "4 Is", "5 A", "6 Test"}

	// Save the array of messages in Redis
	for _, msg := range messages {
		err := redisClient.RPush(ctx, "inboxtest:"+id.String(), msg).Err()
		if err != nil {
			log.Fatalf("Could not push message: %v", err)
		}
	}

	retrievedMessages, err := redisClient.LRange(ctx, "inboxtest:"+id.String(), 0, -1).Result()
	if err != nil {
		log.Fatalf("Could not retrieve messages: %v", err)
	}

	fmt.Println("Retrieved messages:", retrievedMessages)

	err = redisClient.Del(ctx, "inboxtest:"+id.String()).Err()
	if err != nil {
		log.Fatalf("Could not delete key: %v", err)
	}

	retrievedMessages, err = redisClient.LRange(ctx, "inboxtest:"+id.String(), 0, -1).Result()
	if err != nil {
		log.Fatalf("Could not retrieve messages: %v", err)
	}

	fmt.Println("Retrieved messages:", retrievedMessages)

}
