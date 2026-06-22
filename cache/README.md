# cache

Redis 相关能力分两层：

| 组件 | 用途 |
|------|------|
| `cache.C` | 带本地 LFU 的对象缓存（`Get`/`Set`/`Exists`），如 token 校验 |
| `cache.RedisContext` | 裸 Redis 命令（KV、pub/sub），通过 `app.GetRedis()` 获取 |

## RedisContext

`EnableCache` 启动后会注册 `RedisSession`；业务使用：

```go
r := app.GetRedis()
val, err := r.Get(ctx, "key")
if cache.IsNotFound(err) {
    // key 不存在（内部已映射 redis.Nil）
}
_ = r.Set(ctx, "key", "value", time.Minute)
_ = r.Publish(ctx, "app:logs", line)
```

或使用对称 helper：

```go
app.Redis(func(r cache.RedisContext) error {
    return r.Set(ctx, "k", "v", time.Minute)
})
```

## 错误处理

### 两条路径，不要混用判断函数

| 路径 | 你拿到的 `err` | 该怎么判断 |
|------|----------------|------------|
| **框架 API**（`RedisContext` / `cache.C`） | 已在边界 `NormalizeErr`，为 `cache.ErrNotFound` | `cache.IsNotFound(err)` |
| **直接用 go-redis**（`client.Get` 等，未走本包封装） | 原始 `redis.Nil` | `cache.IsRedisNil(err)` |

框架 API 在 `redis_session.go`、`redis.go` 出口统一调用 `NormalizeErr`：**不会**把 `redis.Nil` 原样抛给业务。

`cache.ErrNotFound` 是包内归一化哨兵（`errors.New("cache: not found")`），供 `NormalizeErr` 返回、`IsNotFound` 比较。业务**不必**也**不应**手写 `errors.Is(err, cache.ErrNotFound)`，统一走下方工具函数。

### 工具函数（`cache/errors.go`）

| 函数 | 只适用于 |
|------|----------|
| `IsNotFound(err)` | 框架 API 返回的 err（已归一化） |
| `IsRedisNil(err)` | 直接 go-redis 返回的原始 err |
| `IsNotFoundAny(err)` | 来源不确定、或需同时兼容两种 err 时 |

### 示例

**走框架（推荐）：**

```go
val, err := app.GetRedis().Get(ctx, "key")
if cache.IsNotFound(err) {
    // key 不存在
}
```

**直接用 go-redis（少数场景）：**

```go
val, err := client.Get(ctx, "key").Result()
if cache.IsRedisNil(err) {
    // key 不存在
}
```

业务**不要** `import "github.com/redis/go-redis/v9"` 判断 `redis.Nil`；走框架用 `IsNotFound`，直连 driver 用 `IsRedisNil`。
