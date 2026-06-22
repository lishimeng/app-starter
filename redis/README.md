# redis

业务通过 app-starter 使用 Redis，**不要** `import "github.com/redis/go-redis/v9"`。

## 初始化

```go
builder.EnableRedis(redis.Options{
    Addr:     addr,
    Password: pwd,
})
```

## 使用

```go
c := app.GetRedisClient()
if err := c.LPush(ctx, "queue", payload).Err(); err != nil { ... }
val, err := c.RPop(ctx, "queue").Result()
if redis.IsRedisNil(err) { ... }
```

`redis.Client` 嵌入 go-redis client，所有命令（KV、List、Pub/Sub 等）直接可用，无需逐条封装。

## 错误处理

| 路径 | 判断 |
|------|------|
| `app.GetRedisClient()` 命令返回的 err | `redis.IsRedisNil(err)` |

带本地 LFU 的对象缓存见 [`cache/README.md`](../cache/README.md)（`cache.IsNotFound`）。
