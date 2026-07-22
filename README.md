# App Starter

[![Go project version](https://badge.fury.io/go/github.com%2Flishimeng%2Fapp-starter.svg)](https://badge.fury.io/go/github.com%2Flishimeng%2Fapp-starter)
[![issues](https://img.shields.io/github/issues/lishimeng/app-starter)](https://github.com/lishimeng/app-starter)
[![Go Report Card](https://goreportcard.com/badge/github.com/lishimeng%2Fapp-starter?style=flat-square)](https://goreportcard.com/report/github.com/lishimeng%2Fapp-starter)
[![Source graph](https://sourcegraph.com/github.com/lishimeng/app-starter/-/badge.svg)](https://sourcegraph.com/github.com/lishimeng/app-starter?badge)
[![License](https://img.shields.io/github/license/lishimeng/app-starter)](https://github.com/lishimeng/app-starter)

## Admin 端口（默认 `:6060`）

应用启用 Web 后，默认会在独立端口启动 admin listener（与业务 Web 端口分离）：

| 路径 | 说明 |
|------|------|
| `GET /pprof/*` | Go pprof |
| `GET /metrics` | Prometheus metrics |
| `POST /cl` | 运行时改日志级别，body 示例：`{"level":"debug"}` |

关闭 admin：`builder.SetPprofListen("")`  
改监听地址：`builder.SetPprofListen(":7070")`  
`SetWebLogLevel("DEBUG")` 时，admin 启动也会打印 Gin 路由列表（与业务 Web 一致）。

### 注册业务 Admin API

通过 `EnableAdminRoutes` 挂自定义路由。参数是 `server.Router`（**不要**直接依赖 `gin`）。

框架会先注册 `/pprof`、`/metrics`、`POST /cl`，再执行业务 Setup，业务路由不会覆盖框架 API。

```go
import (
    "github.com/lishimeng/app-starter"
    "github.com/lishimeng/app-starter/server"
)

func RegisterAdmin(root server.Router) {
    root.Get("/demo/ping", func(ctx server.Context) {
        ctx.JSON(map[string]any{"ok": true})
    })
}

func main() {
    app.New().Start(func(ctx context.Context, builder *app.ApplicationBuilder) error {
        builder.
            EnableAdminRoutes(RegisterAdmin).
            EnableWeb(":9527", YourAPIRouter)
        return nil
    })
}
```

访问示例：

```text
GET  http://localhost:6060/demo/ping
GET  http://localhost:6060/metrics
POST http://localhost:6060/cl   {"level":"debug"}
```

完整示例见 [`examples/web-basic`](examples/web-basic)（含 `admin.Register` 与测试）。

Delete tag
---

```shell
$ git tag -d v0.4.0
Deleted tag 'v0.4.0' (was f74dcae)

$ git push origin :v0.4.0
To https://github.com/lishimeng/app-starter.git
 - [deleted]         3.3.0.1492
```
