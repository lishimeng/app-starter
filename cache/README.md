# cache

带本地 TinyLFU 的 **对象缓存**（基于 [go-redis/cache](https://github.com/go-redis/cache)），适用于 token 校验等需要「Redis + 进程内缓存」的场景。

裸 Redis 命令（KV、List、Pub/Sub 等）见 [`redis/README.md`](../redis/README.md)。

## 与 redis 的关系

| 包 | 用途 | 业务入口 |
|----|------|----------|
| [`redis`](../redis/) | 共享 Redis 连接，全命令集 | `app.GetRedisClient()` |
| `cache` | 本地 LFU + Redis 对象缓存 | `app.GetCache()` |

**初始化有依赖：** 必须先 `EnableRedis`，再 `EnableCache`（共用同一 `redis.Client`）。

```go
builder.
    EnableRedis(redis.Options{
        Addr:     addr,
        Password: pwd,
    }).
    EnableCache(cache.Options{
        MaxSize: 10000,
        Ttl:     time.Minute,
    })
```

仅 Redis、不需要本地 cache 时，只调用 `EnableRedis` 即可，不必 `EnableCache`。

## 使用

```go
c := app.GetCache()

var obj MyStruct
err := c.Get("key", &obj)
if cache.IsNotFound(err) {
    // key 不存在
}

_ = c.Set("key", &obj)
_ = c.SetTTL("key", &obj, time.Hour)
_ = c.Del("key")
ok := c.Exists("key")
```

`GetSkipLocal` 跳过本地 LFU，直接读 Redis。

## 错误处理

`cache.C` 在 `Get` / `Set` / `Del` 出口调用 `NormalizeErr`，返回 **`cache.ErrNotFound`**，不是原始 `redis.Nil`。

| 函数 | 说明 |
|------|------|
| `cache.IsNotFound(err)` | `cache.C` 返回的 err（已归一化） |
| `cache.IsNotFoundAny(err)` | 兼容原始 `redis.Nil` 或归一化 err |

业务**不要** `import "github.com/redis/go-redis/v9"`。  
使用 `app.GetRedisClient()` 时的 err 判断见 [`redis.IsRedisNil`](../redis/errors.go)。

```go
err := app.GetCache().Get(key, &obj)
if cache.IsNotFound(err) {
    // missing key
}
```
