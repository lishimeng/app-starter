package redis

import goredis "github.com/redis/go-redis/v9"

// Options mirrors go-redis client options; business uses this type only.
type Options goredis.Options
