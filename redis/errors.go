package redis

import (
	"errors"

	goredis "github.com/redis/go-redis/v9"
)

// IsRedisNil reports missing-key errors from direct Client command results (.Result() / .Err()).
// Use this when calling app.GetRedisClient(); do not import go-redis in business code.
func IsRedisNil(err error) bool {
	return err != nil && errors.Is(err, goredis.Nil)
}
