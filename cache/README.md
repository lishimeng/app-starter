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

业务**不要** `import "github.com/redis/go-redis/v9"` 判断 `redis.Nil`：

```go
if cache.IsNotFound(err) {
    // missing key
}
```

`cache.C` 的 `Get`/`Del` 同样会归一化 `redis.Nil`。
