package repository

import (
	"context"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

type ILimitRepository interface {
	isBlocked(ctx context.Context, ip net.IP) (bool, error)
	Increment(ctx context.Context, ip net.IP) error
	Reset(ctx context.Context, ip net.IP) error
}

type LimitsRepository struct {
	rdb *redis.Client
}

const (
	maxAttempts = 5
	lockoutTime = 15 * time.Minute
)

func NewLimiterRepository(rdb *redis.Client) *LimitsRepository {
	return &LimitsRepository{rdb: rdb}
}

func (r *LimitsRepository) IsBlocked(ctx context.Context, ip string) (bool, error) {
	count, err := r.rdb.Get(ctx, "login_attempts:"+ip).Int()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return count >= maxAttempts, nil
}

func (r *LimitsRepository) Increment(ctx context.Context, ip net.IP) error {
	key := "login_attempts:" + ip.String()
	pipe := r.rdb.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, lockoutTime)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *LimitsRepository) Reset(ctx context.Context, ip net.IP) error {
	return r.rdb.Del(ctx, "login_attempts:"+ip.String()).Err()
}
