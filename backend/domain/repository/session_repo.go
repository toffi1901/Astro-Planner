package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Session struct {
	UserID    string    `json:"user_id"`
	Role      string    `json: "role"`
	IPadress  string    `json: "ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

type SessionRepository struct {
	rdb *redis.Client
}

func NewSessionRepository(rdb *redis.Client) *SessionRepository {
	return &SessionRepository{rdb: rdb}
}

func (r *SessionRepository) Create(ctx context.Context, jti string, s Session, ttl time.Duration) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, "session:"+jti, data, ttl).Err()
}

func (r *SessionRepository) Get(ctx context.Context, jti string) (*Session, error) {
	data, err := r.rdb.Get(ctx, "session:"+jti).Bytes()
	if err != nil {
		return nil, err
	}

	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepository) Delete(ctx context.Context, jti string) error {
	return r.rdb.Del(ctx, "session:"+jti).Err()
}
