package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"go-redis"

	"time"
)

type IResetPasswdRepository interface {
	Create(ctx context.Context, userId string, ttl time.Duration) (string, error)
	Consume(ctx context.Context, token string) (string, error)
}

type ResetPasswdRepository struct {
	rdb *redis.Client
}

func NewpResetPasswdRepository(rdb *redis.Client) *ResetPasswdRepository {
	return &ResetPasswdRepository{rdb: rdb}
}

func (r *ResetPasswdRepository) Create(ctx context.Context, userId string, ttl time.Duration) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(bytes)
	if err := r.rdb.Set(ctx, "reset_password"+token, userId, ttl).Err(); err != nil {
		return "", nil
	}
	return token, nil

}
func (r *ResetPasswdRepository) Consume(ctx context.Context, token string) (string, error) {
	key := "password_reset" + token
	userID, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}
	r.rdb.Del(ctx, key).Err()
	return UserID, nil
}
