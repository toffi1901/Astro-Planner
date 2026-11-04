package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Session struct {
	JTI       string    `json:"jti"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"`
	IPadress  string    `json:"ip_address"`
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
	pipe := r.rdb.Pipeline()
	pipe.Set(ctx, "session:"+jti, data, ttl)
	pipe.SAdd(ctx, "user_sessions:"+s.UserID, jti)
	pipe.Expire(ctx, "user_sessions:"+s.UserID, ttl)

	_, err = pipe.Exec(ctx)
	return err
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
	sess, err := r.Get(ctx, jti)
	if err == nil && sess != nil {
		r.rdb.SRem(ctx, "user_sessions:"+sess.UserID, jti)
	}

	return r.rdb.Del(ctx, "session:"+jti).Err()
}

func (r *SessionRepository) DeleteAllUserSessions(ctx context.Context, userID string) error {
	jtis, err := r.rdb.SMembers(ctx, "user_sessions:"+userID).Result()
	if err != nil {
		return err
	}

	if len(jtis) == 0 {
		return nil
	}

	pipe := r.rdb.Pipeline()
	for _, jti := range jtis {
		pipe.Del(ctx, "session:"+jti)
	}
	pipe.Del(ctx, "user_sessions:"+userID)

	_, err = pipe.Exec(ctx)
	return err
}

func (r *SessionRepository) ListByUser(ctx context.Context, userId string) ([]*Session, error) {
	jtis, err := r.rdb.SMembers(ctx, "user_sessions:"+userId).Result()
	if err != nil {
		return nil, err
	}
	var sessions []*Session
	for _, jti := range jtis {
		sess, err := r.Get(ctx, jti)
		if err == nil && sess != nil {
			sess.JTI = jti
			sessions = append(sessions, sess)
		}
	}
	return sessions, nil
}
