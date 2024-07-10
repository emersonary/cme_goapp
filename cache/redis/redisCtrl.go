package rediscache

import (
	"context"
	"encoding/json"
	"log"

	"github.com/emersonary/go-authentication/model/message"
	"github.com/emersonary/go-authentication/pck"
	"github.com/go-redis/redis/v8"
)

var RedisCtrl *TRedisCtrl

type TRedisCtrl struct {
	redisSession *redis.Client
}

var ctx = context.Background()

func NewRedisControl(redisSession *redis.Client) *TRedisCtrl {
	return &TRedisCtrl{redisSession: redisSession}
}

func (r *TRedisCtrl) AddCache(message *message.TMessage) error {

	key := "inbox:" + message.ToUserId.String()

	jsonStr, err := json.Marshal(message)

	if err != nil {
		return err
	}

	err = r.redisSession.RPush(ctx, key, jsonStr).Err()

	return err

}

func (r *TRedisCtrl) DeleteFirstMessages(inboxId pck.UUID, qty int) error {

	key := "inbox:" + inboxId.String()

	for i := 1; i <= qty; i++ {

		err := r.redisSession.LPop(ctx, key).Err()

		if err != nil {

			return err

		}

	}

	return nil

}

func (r *TRedisCtrl) CachedMessages(inboxId pck.UUID) (*[]message.TMessage, error) {

	key := "inbox:" + inboxId.String()
	result := make([]message.TMessage, 0)

	retrievedMessages, err := r.redisSession.LRange(ctx, key, 0, -1).Result()

	if err != nil {
		return nil, err
	}

	for _, strJson := range retrievedMessages {

		var message message.TMessage

		err := json.Unmarshal([]byte(strJson), &message)
		if err != nil {
			log.Fatalf("Error unmarshaling JSON: %v", err)
		}

		result = append(result, message)

	}

	return &result, nil

}
