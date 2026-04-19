package redis

import (
	"context"
	"log"
	"strings"

	"github.com/redis/go-redis/v9"
)

// StartSubscriber listens for expired session keys and decrements the
// per-event active counter. Requires Redis keyspace notifications enabled
// with: redis-server --notify-keyspace-events Ex
//
// onSessionEnd is called after a session expires and the counter is decremented,
// allowing the caller to admit the next queued user.
func StartSubscriber(ctx context.Context, client *redis.Client, onSessionEnd func(ctx context.Context, eventID string)) {
	pubsub := client.PSubscribe(ctx, "__keyevent@0__:expired")
	defer pubsub.Close()

	log.Println("session subscriber: listening for expired keys")

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			handleExpiredKey(ctx, client, msg.Payload, onSessionEnd)
		}
	}
}

// handleExpiredKey processes an expired key notification.
// Session keys have the format: session:{eventID}:{sessionID}
func handleExpiredKey(ctx context.Context, client *redis.Client, key string, onSessionEnd func(ctx context.Context, eventID string)) {
	if !strings.HasPrefix(key, "session:") {
		return
	}

	parts := strings.SplitN(key, ":", 3)
	if len(parts) < 3 {
		return
	}
	eventID := parts[1]

	// Decrement the active counter, floor at 0
	count, err := client.Decr(ctx, activeKey(eventID)).Result()
	if err != nil {
		log.Printf("session subscriber: failed to decr counter for event %s: %v", eventID, err)
		return
	}
	if count < 0 {
		client.Set(ctx, activeKey(eventID), 0, 0)
	}

	log.Printf("session subscriber: session expired for event %s, active count: %d", eventID, max(count, 0))

	if onSessionEnd != nil {
		onSessionEnd(ctx, eventID)
	}
}
