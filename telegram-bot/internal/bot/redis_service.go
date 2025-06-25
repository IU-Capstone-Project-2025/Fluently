package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// RedisService provides helper methods for Redis operations
type RedisService struct{}

// NewRedisService creates a new Redis service instance
func NewRedisService() *RedisService {
	return &RedisService{}
}

// SetUserData stores user data in Redis with expiration
func (rs *RedisService) SetUserData(userID int64, data interface{}, expiration time.Duration) error {
	ctx := context.Background()
	key := getUserKey(userID)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return RedisClient.Set(ctx, key, jsonData, expiration).Err()
}

// GetUserData retrieves user data from Redis
func (rs *RedisService) GetUserData(userID int64, dest interface{}) error {
	ctx := context.Background()
	key := getUserKey(userID)

	jsonData, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(jsonData), dest)
}

// DeleteUserData removes user data from Redis
func (rs *RedisService) DeleteUserData(userID int64) error {
	ctx := context.Background()
	key := getUserKey(userID)

	return RedisClient.Del(ctx, key).Err()
}

// SetUserState sets user state in Redis
func (rs *RedisService) SetUserState(userID int64, state string, expiration time.Duration) error {
	ctx := context.Background()
	key := getUserStateKey(userID)

	return RedisClient.Set(ctx, key, state, expiration).Err()
}

// GetUserState gets user state from Redis
func (rs *RedisService) GetUserState(userID int64) (string, error) {
	ctx := context.Background()
	key := getUserStateKey(userID)

	return RedisClient.Get(ctx, key).Result()
}

// Helper functions for Redis keys
func getUserKey(userID int64) string {
	return fmt.Sprintf("user:%d", userID)
}

func getUserStateKey(userID int64) string {
	return fmt.Sprintf("user_state:%d", userID)
}
