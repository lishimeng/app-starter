# log

基于标准库 `slog` 的轻量封装，import 路径：

```go
import "github.com/lishimeng/app-starter/log"
```

**注意**：与标准库 `log` 同名，同一文件不要混用 `import "log"`。

## 快速开始

启动时（`application` 已内置，也可手动）：

```go
log.Config().LevelFromString("INFO").Text().Apply()
```

环境变量：`LOG_LEVEL`（默认 `INFO`）。Builder：`SetAppLogLevel("DEBUG")`（与 Iris `SetWebLogLevel` 无关）。

## 用法

API 与 `slog` 一致：`msg` + 可选 key-value 对；格式化用 `*f`。

```go
log.Info("server started")
log.Infof("listen %s", addr)
log.Info("query", "pageNum", pageNum, "pageSize", pageSize)
log.With("err", err).Error("verify failed")
```

`module` 从调用栈文件路径推断（`app-starter/` 后的包路径，如 `mqtt`、`application/api`）。

### 可选 `For`：显式模块名

```go
log.For("syncdb").Infof("create table %s", table)
log.For("").Info("same as default auto module")
```

### 链式

```go
log.With("err", err).Error("verify failed")
log.For("mqtt").With("client", id).Info("connected")
```

### 热路径：包内固定 logger

```go
var logger = log.For("mqtt")

func Connect() {
    logger.Info("connected")
}
```

## 级别

| 字符串 | slog |
|--------|------|
| FINEST, FINE, DEBUG | Debug |
| TRACE, INFO | Info |
| WARNING, WARN | Warn |
| ERROR, CRITICAL | Error |

运行时调整：`log.SetLevelFromString("ERROR")`（HTTP API：`application/api/log.level.go`）。

## 配置链

```go
log.Config().
    LevelFromString("DEBUG").
    Text().              // 或 JSON()
    Output(os.Stderr).
    Apply()
```
