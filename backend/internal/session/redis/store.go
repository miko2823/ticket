package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/KaoriNakajima/sturdyticket/backend/internal/session"
)

// Store is the Redis adapter for the session Store port.
type Store struct {
	client *redis.Client
}

func NewStore(client *redis.Client) *Store {
	return &Store{client: client}
}

func sessionKey(eventID, sessionID string) string {
	return fmt.Sprintf("session:%s:%s", eventID, sessionID)
}

func activeKey(eventID string) string {
	return fmt.Sprintf("event:%s:active", eventID)
}

func userEventKey(userID, eventID string) string {
	return fmt.Sprintf("user:%s:event:%s:session", userID, eventID)
}

// Create atomically creates a session, sets the user-event mapping, and increments the counter.
// Returns existing sessionID if the user already has a session for this event.
var createScript = redis.NewScript(`
local existing = redis.call('GET', KEYS[3])
if existing then return existing end
redis.call('HSET', KEYS[1], 'userID', ARGV[1], 'createdAt', ARGV[2])
redis.call('EXPIRE', KEYS[1], ARGV[3])
redis.call('SET', KEYS[3], ARGV[4], 'EX', ARGV[3])
redis.call('INCR', KEYS[2])
return 'ok'
`)

func (s *Store) Create(ctx context.Context, sess *session.Session, ttl time.Duration) error {
	ttlSec := int(ttl.Seconds())
	result, err := createScript.Run(ctx, s.client,
		[]string{
			sessionKey(sess.EventID, sess.ID),
			activeKey(sess.EventID),
			userEventKey(sess.UserID, sess.EventID),
		},
		sess.UserID,
		sess.CreatedAt.Unix(),
		ttlSec,
		sess.ID,
	).Result()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	// If the script returned an existing sessionID, update the session ID
	if str, ok := result.(string); ok && str != "ok" {
		sess.ID = str
	}
	return nil
}

func (s *Store) Refresh(ctx context.Context, sessionID, eventID string, ttl time.Duration) error {
	key := sessionKey(eventID, sessionID)
	exists, err := s.client.Expire(ctx, key, ttl).Result()
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("session not found")
	}
	// Also refresh the user-event mapping
	userID, err := s.client.HGet(ctx, key, "userID").Result()
	if err != nil {
		return err
	}
	s.client.Expire(ctx, userEventKey(userID, eventID), ttl)
	return nil
}

// End atomically removes the session and decrements the counter.
var endScript = redis.NewScript(`
local userID = redis.call('HGET', KEYS[1], 'userID')
if not userID then return 0 end
redis.call('DEL', KEYS[1])
if KEYS[3] then redis.call('DEL', KEYS[3]) end
local count = redis.call('DECR', KEYS[2])
if count < 0 then redis.call('SET', KEYS[2], 0) end
return 1
`)

func (s *Store) End(ctx context.Context, sessionID, eventID string) error {
	key := sessionKey(eventID, sessionID)
	// Get userID before deletion to build user-event key
	userID, _ := s.client.HGet(ctx, key, "userID").Result()
	ueKey := userEventKey(userID, eventID)

	_, err := endScript.Run(ctx, s.client,
		[]string{key, activeKey(eventID), ueKey},
	).Result()
	return err
}

func (s *Store) Find(ctx context.Context, sessionID, eventID string) (*session.Session, error) {
	key := sessionKey(eventID, sessionID)
	vals, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(vals) == 0 {
		return nil, fmt.Errorf("session not found")
	}
	ts, _ := strconv.ParseInt(vals["createdAt"], 10, 64)
	return &session.Session{
		ID:        sessionID,
		EventID:   eventID,
		UserID:    vals["userID"],
		CreatedAt: time.Unix(ts, 0),
	}, nil
}

func (s *Store) FindByUserEvent(ctx context.Context, userID, eventID string) (string, error) {
	sessionID, err := s.client.Get(ctx, userEventKey(userID, eventID)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return sessionID, err
}

func (s *Store) GetActiveCount(ctx context.Context, eventID string) (int, error) {
	val, err := s.client.Get(ctx, activeKey(eventID)).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	count, _ := strconv.Atoi(val)
	if count < 0 {
		return 0, nil
	}
	return count, nil
}
