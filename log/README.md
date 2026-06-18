# log

基于标准库 `slog` 的轻量封装，import 路径：

```go
import "github.com/lishimeng/app-starter/log"
```

**注意**：与标准库 `log` 同名，同一文件不要混用 `import "log"`。

## 默认配置

包加载时（`init`）已生效，**不调用** `log.Config().Apply()` 也使用下表默认值：

| 项 | 默认值 | 说明 |
|----|--------|------|
| 级别 | `INFO` | `slog.LevelInfo` |
| 格式 | Text | `slog` 文本 handler（`time=... level=... msg=...`） |
| 输出 | `os.Stdout` | 标准输出 |
| `source=` | 关 | 无 `source=file:line` 字段 |
| `module` | 自动 | 每条业务日志推断包路径（`app-starter/` 后），非 Config 项 |
| 时间格式 | slog 内置 | RFC3339 纳秒，如 `2026-06-18T17:43:01.524+08:00`，不可配置 |

`log.Config()` 未显式设置的链式项与上表一致；`Apply()` 前若已调过 `SetLevelFromString`，级别取当前全局值。

示例（均为显式写出默认值，等价于仅 `Apply()`）：

```go
log.Config().
    LevelFromString("INFO").
    Text().
    Out(os.Stdout).
    Caller(false).
    Apply()
```

## 快速开始

启动时在 `main` 中配置（application 不再代为初始化）：

```go
log.Config().LevelFromString("INFO").Text().Apply()
```

可选源码位置（默认关闭）：`log.Config().Caller(true).Apply()`。业务 log 入口栈过滤（跳过 `log/` 封装）；GORM SQL 由 `persistence` 自行过滤（跳过 `gorm.io/`）。

```text
time=... level=ERROR msg="transaction fail" module=examples/web-basic/router source=examples/web-basic/router/transaction_fail.api.go:42 err="..."
```

## 用法

API 与 `slog` 一致：`msg` + 可选 key-value 对；格式化用 `*f`。

```go
log.Info("server started")
log.Infof("listen %s", addr)
log.Info("query", "pageNum", pageNum, "pageSize", pageSize)
log.With("err", err).Error("verify failed")
```

`module` 从调用栈推断。`source`（可选）由各入口栈过滤后写入（跳过 `log/`、`gorm.io/` 等封装层），相对路径 `file:line`。

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

| 方法 | 作用 | 默认 |
|------|------|------|
| `Level` / `LevelFromString` | 日志级别 | `INFO` |
| `Text` / `JSON` | 输出格式 | Text |
| `Out(w)` | 输出目标 | `os.Stdout` |
| `Caller(bool)` | `source=file:line` | `false` |
| `Apply()` | 生效并设为 `slog.Default` | — |

```go
log.Config().
    LevelFromString("DEBUG").
    Caller(true).
    Text().              // 或 JSON()
    Out(os.Stderr).      // 不写则 stdout
    Apply()
```
