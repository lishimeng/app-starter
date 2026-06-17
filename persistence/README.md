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
| `InitDb` | `true` 时启动后执行 **SyncDB**（轻量级建表/补列/补索引；批量加载元数据后内存 diff） |
| `SyncForce` | `true` 时 SyncDB 先删表再重建（**会丢数据**，仅开发环境使用） |
| `SyncVerbose` | `true` 时 SyncDB 打印每一步 DDL 动作 |
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

## SyncDB（轻量级表结构同步）

`InitDb=true` 时执行 **SyncDB**，语义对齐 Beego `orm.RunSyncdb`，**不调用** GORM `AutoMigrate`。

实现上先**批量加载**库表/列/索引元数据（Postgres/MySQL 各 3 条 SQL；SQLite 2 条 + 按表 PRAGMA），在内存 diff 后仅对缺失项执行 DDL，避免逐列 `Has*` 查询。

| 操作 | SyncDB |
|------|--------|
| 表不存在 → 创建（含主键与索引） | 是 |
| 表已存在 → 补缺失列 | 是 |
| 表已存在 → 补缺失索引 | 是 |
| `SyncForce=true` → 删表重建 | 是（有数据丢失风险） |
| `SyncVerbose=true` → 打印 DDL 动作 | 是 |
| 修改已有列类型/约束 | **否** |
| 删除模型中已移除的列/索引 | **否** |

手动调用（等价 Beego `RunSyncdb`）：

```go
import "github.com/lishimeng/app-starter/persistence"

err := persistence.RunSyncDB(
    persistence.DefaultAlias,
    persistence.SyncOptions{Verbose: true},
    &model.YourModel{},
)
```

driver Config 示例：

```go
postgres.Config{
    InitDb:      true,
    SyncVerbose: true,
    // SyncForce: true, // 仅开发：删表重建
}.Build()
```

模型字段/索引从模型中删除后，数据库残留需**人工处理**（与 Beego 一致）。

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

## 查询 API

`Query` 由 `QueryCond`（条件）与 `QueryExec`（排序/分页/执行）嵌入组成，对外仍使用 `persistence.Query`。

### 条件封装（QueryCond，推荐）

减少手写 SQL 表达式，链式调用：

| 方法 | 语义 |
|------|------|
| `Equal(col, val)` | `col = ?` |
| `NotEqual(col, val)` | `col <> ?` |
| `In(col, vals)` | `col IN ?` |
| `Like(col, s)` | `LIKE %s%` |
| `LLike(col, s)` | `LIKE s%`（前缀） |
| `RLike(col, s)` | `LIKE %s`（后缀） |
| `ILike(col, s)` | `ILIKE %s%`（PostgreSQL） |
| `EqualStr` / `LikeStr` / `ILikeStr` 等 | 值为空时跳过条件 |

```go
tx.Model(&model.User{}).
    Equal("status", 1).
    ILikeStr("name", keyword).
    EqualStr("code", code).
    Order("id desc").
    Limit(10).
    Find(&list)
```

复杂条件仍可使用 `Where("a = ? AND b > ?", x, y)`。

### 执行与分页（QueryExec）

`Select`、`Omit`、`Order`、`Offset`、`Limit`、`Count`、`Find`、`First`、`Take`、`Update`、`Updates`。

### GORM 原生

```go
tx.Model(&model.User{}).Where("status = ?", 1).Find(&list)
```

分页查询使用 `app.SimplePager` + `app.QueryPage`。

---

## 参考链接

- [GORM - 连接到数据库](https://gorm.io/zh_CN/docs/connecting_to_the_database.html)
- [GORM - 链式操作](https://gorm.io/zh_CN/docs/query.html)
