# Persistence（GORM）

app-starter 通过 `persistence` 包封装 GORM，业务侧使用 `app.EnableDatabase` 启用数据库，查询 API 对齐 [GORM 链式风格](https://gorm.io/zh_CN/docs/)。

## 包结构

```
persistence/                 # 核心：接口、GORM 实现、连接注册（不含任何 driver 依赖）
persistence/driver/postgres/ # postgres.Config + 硬编码 gorm.io/driver/postgres
persistence/driver/mysql/    # mysql.Config + 硬编码 gorm.io/driver/mysql
persistence/driver/sqlite/   # sqlite.Config + 硬编码 glebarez/sqlite
```

**编译隔离**：只有被 import 的 driver 子包会打进二进制。`persistence` 核心包本身不链接任何数据库驱动。

**自动注册**：import driver 子包时，`init()` 自动 `RegisterDialector`，无需 `RegisterDialector`、无需 `import gorm.io/driver/*`。

## 快速开始

```go
import "github.com/lishimeng/app-starter/persistence/driver/postgres"

builder.EnableDatabase(setup.PostgresConfig().Build(), new(model.YourModel))
```

`setup.go` 使用 `postgres.Config` 构建连接（import 该子包即完成 driver 注册）：

```go
import "github.com/lishimeng/app-starter/persistence/driver/postgres"

func PostgresConfig() *postgres.Config {
    return &postgres.Config{
        UserName: os.Getenv("DB_USER"),
        Host:     os.Getenv("DB_HOST"),
        DbName:   os.Getenv("DB_DATABASE"),
        TimeZone: "Asia/Shanghai",
    }
}
```

环境变量示例见 `examples/web-basic/setup/setup.go`。

---

## 连接配置

| 子包 | Config 类型 | Driver 名称 |
|------|-------------|-------------|
| `persistence/driver/postgres` | `postgres.Config` | `postgres` |
| `persistence/driver/mysql` | `mysql.Config` | `mysql` |
| `persistence/driver/sqlite` | `sqlite.Config` | `sqlite3` |

`Build()` 生成 `persistence.BaseConfig` 后传给 `EnableDatabase`。

`BaseConfig` 公共字段：

| 字段 | 说明 |
|------|------|
| `InitDb` | `true` 时启动后执行 `AutoMigrate` |
| `AliasName` | 连接别名，默认 `default` |
| `MaxIdleConns` / `MaxOpenConns` | 连接池（见下文） |
| `Debug` | 是否打印 SQL |
| `DriverOpts` | 驱动专属选项，由 `Config.Build()` 自动填充 |

### Build 内置的数据库特性

**postgres.Config**

| 字段 | 作用 |
|------|------|
| `TimeZone` | 追加 `TimeZone=...` 到 DSN |
| `AdvancedConfig` | 追加 sslcert、sslkey 等额外 DSN 参数 |
| `PreferSimpleProtocol` | `true` 时禁用 pgx 隐式 prepared statement 缓存 |

```go
postgres.Config{
    UserName:             "postgres",
    Host:                 "127.0.0.1",
    DbName:               "mydb",
    TimeZone:             "Asia/Shanghai",
    PreferSimpleProtocol: true,
}.Build()
```

**mysql.Config**

| 字段 | 作用 |
|------|------|
| `Charset` | 默认 `utf8mb4` |
| `DisableParseTime` | 默认启用 `parseTime=True` |
| `Loc` | 默认 `Local` |
| `DefaultStringSize` 等 | 非零/为 true 时使用 `mysql.New(Config{...})` |

```go
mysql.Config{
    UserName: "root",
    Host:     "127.0.0.1",
    Port:     3306,
    DbName:   "mydb",
}.Build()
```

---

## 各数据库连接注意事项

整理自 GORM 官方文档：[连接到数据库](https://gorm.io/zh_CN/docs/connecting_to_the_database.html)。

### PostgreSQL

- 时区：使用 `postgres.Config.TimeZone` 或 `AdvancedConfig`
- prepared statement：`PreferSimpleProtocol: true` 可禁用 pgx 缓存

### MySQL

- 默认 DSN：`charset=utf8mb4&parseTime=True&loc=Local`
- TiDB 兼容 MySQL 协议

### SQLite

- 默认纯 Go 驱动 `glebarez/sqlite`
- 内存库：`file::memory:?cache=shared`

### 其他数据库

SQL Server、Oracle、ClickHouse 等可通过 `persistence.RegisterDialector` 在业务侧扩展，或新增 `persistence/driver/xxx` 子包。

---

## 连接池

`BaseConfig.MaxIdleConns` / `MaxOpenConns` 在 `Open` 时设置 `sqlDB.SetMaxIdleConns` / `SetMaxOpenConns`。

---

## 自定义 dialector

```go
persistence.RegisterDialector("postgres", func(opts persistence.OpenOptions) gorm.Dialector {
    return pgdriver.New(pgdriver.Config{DSN: opts.DSN})
})
```

后注册会覆盖子包 `init()` 中的注册。常规场景使用 `postgres.Config` 等 Build 字段即可。

---

## 查询 API（GORM 风格）

```go
tx.Model(&model.User{}).
    Where("status = ?", 1).
    Order("id desc").
    Limit(10).
    Find(&list)
```

分页查询使用 `app.SimplePager` + `app.QueryPage`。

---

## 参考链接

- [GORM - 连接到数据库](https://gorm.io/zh_CN/docs/connecting_to_the_database.html)
- [GORM - 链式操作](https://gorm.io/zh_CN/docs/query.html)
