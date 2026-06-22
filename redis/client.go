package redis

import goredis "github.com/redis/go-redis/v9"

// Client wraps go-redis; embeds the driver client so all commands (LPush, RPop, …) are available.
// Business imports github.com/lishimeng/app-starter/redis only, not go-redis.
type Client struct {
	*goredis.Client
}

// NewClient creates a Redis client from Options.
func NewClient(opts Options) *Client {
	return &Client{Client: goredis.NewClient((*goredis.Options)(&opts))}
}

// Close closes the client connection.
func (c *Client) Close() error {
	if c == nil || c.Client == nil {
		return nil
	}
	return c.Client.Close()
}
