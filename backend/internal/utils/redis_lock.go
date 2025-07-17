package utils

import (
	"context"
	"sync"
	"time"

	"github.com/bsm/redislock"
	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

var (
	lockOnce sync.Once
	locker   *redislock.Client
)

// redisLocker returns a singleton instance of redislock.Client built on top of the
// shared utils.Redis() connection.
func redisLocker() *redislock.Client {
	lockOnce.Do(func() {
		// We purposefully reuse the singleton Redis connection that is already
		// configured through environment variables by utils.Redis().
		var client *goredis.Client = Redis()
		locker = redislock.New(client)
	})
	return locker
}

// AcquireChatLock obtains a distributed mutex for chat-related operations of a
// particular user. The key is scoped to the user so concurrent requests for
// **different** users will not block each other, while overlapping requests for
// the **same** user (e.g. /chat and /chat/finish arriving at the same time)
// will be serialized. The lock automatically expires after a short TTL so that
// it does not remain forever if the holder crashes.
func AcquireChatLock(ctx context.Context, userID uuid.UUID) (*redislock.Lock, error) {
	const ttl = 10 * time.Second // plenty for a single chat request
	return redisLocker().Obtain(ctx, "lock:chat:"+userID.String(), ttl, nil)
}
