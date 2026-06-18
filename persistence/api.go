package persistence

// OpenOptions carries connection-level settings for a Connector.
type OpenOptions struct {
	Alias    string
	MaxIdle  int
	MaxOpen  int
	Debug    bool
	InitDB   bool
	Driver   string
	DSN        string
	DriverOpts any // driver-specific options set by *Config.Build()
}

// Connector opens database sessions.
type Connector interface {
	Open(opts OpenOptions) (Session, error)
	Migrate(alias string, opts SyncOptions, models ...any) error
	RegisterModels(models ...any)
}

// Session is the unit of work for database access, typically one per alias.
type Session interface {
	Transaction(fn func(Tx) error) error
	Model(value interface{}) Query
	SetDebug(enable bool)
	Alias() string
}

// Tx represents a transactional database session.
type Tx interface {
	Model(value interface{}) Query
	Create(value interface{}) error
	Save(value interface{}) error
	Delete(value interface{}, conds ...interface{}) error
	First(dest interface{}, conds ...interface{}) error
	Raw(sql string, values ...interface{}) Query
}

// QueryCond 查询条件构建。所有方法均返回 Query，可链式调用。
//
// column 参数使用数据库列名（如 "conn_type"），与 gorm tag 中 column 一致。
// 下文 SQL 以表名 users、列名与示例参数为例；实际表名由 Model 的 TableName/gorm 命名决定。
//
// 典型用法：
//
//	tx.Model(&User{}).
//	    EqualStr("code", "c1").
//	    ILikeStr("name", "abc").
//	    Equal("enabled", 1).Find(&rows)
//
// SQL：
//
//	SELECT * FROM users
//	WHERE code = 'c1' AND name ILIKE '%abc%' AND enabled = 1
type QueryCond interface {
	// Equal 追加 column = ? 条件。
	//
	// 例：tx.Model(&User{}).Equal("id", 1).First(&row)
	//
	// SQL：
	//
	//	SELECT * FROM users WHERE id = 1 ORDER BY users.id ASC LIMIT 1
	Equal(column string, value any) Query

	// NotEqual 追加 column <> ? 条件。
	//
	// 例：tx.Model(&User{}).NotEqual("status", 0).Find(&rows)
	//
	// SQL：
	//
	//	SELECT * FROM users WHERE status <> 0
	NotEqual(column string, value any) Query

	// In 追加 column IN ? 条件；values 为 slice 或数组。
	//
	// 例：tx.Model(&User{}).In("id", []int{1, 2, 3}).Find(&rows)
	//
	// SQL：
	//
	//	SELECT * FROM users WHERE id IN (1, 2, 3)
	In(column string, values any) Query

	// Like 追加 column LIKE %value%（两端模糊匹配）。
	//
	// 例：tx.Model(&User{}).Like("name", "张").Find(&rows)
	//
	// SQL：
	//
	//	SELECT * FROM users WHERE name LIKE '%张%'
	Like(column string, value string) Query

	// LLike 追加 column LIKE value%（前缀匹配）。
	//
	// 例：tx.Model(&User{}).LLike("code", "A").Find(&rows)
	//
	// SQL：
	//
	//	SELECT * FROM users WHERE code LIKE 'A%'
	LLike(column string, value string) Query

	// RLike 追加 column LIKE %value（后缀匹配）。
	//
	// 例：tx.Model(&User{}).RLike("email", "@test.com").Find(&rows)
	//
	// SQL：
	//
	//	SELECT * FROM users WHERE email LIKE '%@test.com'
	RLike(column string, value string) Query

	// ILike 追加 column ILIKE %value%（不区分大小写，PostgreSQL）。
	//
	// 例：tx.Model(&User{}).ILike("name", "abc").Find(&rows)
	//
	// SQL（PostgreSQL）：
	//
	//	SELECT * FROM users WHERE name ILIKE '%abc%'
	ILike(column string, value string) Query

	// EqualStr 当 value 非空时追加 Equal，空字符串跳过该条件（便于接 URL 查询参数）。
	//
	// 例：tx.Model(&User{}).EqualStr("code", "c1").EqualStr("conn_type", "").Find(&rows)
	//
	// SQL（conn_type 为空，该条件不出现）：
	//
	//	SELECT * FROM users WHERE code = 'c1'
	EqualStr(column string, value string) Query

	// LikeStr 当 value 非空时追加 Like，空字符串跳过。
	//
	// 例：tx.Model(&User{}).LikeStr("name", "张").Find(&rows)
	//
	// SQL：
	//
	//	SELECT * FROM users WHERE name LIKE '%张%'
	LikeStr(column string, value string) Query

	// LLikeStr 当 value 非空时追加 LLike，空字符串跳过。
	//
	// 例：tx.Model(&User{}).LLikeStr("code", "A").Find(&rows)
	//
	// SQL：
	//
	//	SELECT * FROM users WHERE code LIKE 'A%'
	LLikeStr(column string, value string) Query

	// RLikeStr 当 value 非空时追加 RLike，空字符串跳过。
	//
	// 例：tx.Model(&User{}).RLikeStr("email", "@test.com").Find(&rows)
	//
	// SQL：
	//
	//	SELECT * FROM users WHERE email LIKE '%@test.com'
	RLikeStr(column string, value string) Query

	// ILikeStr 当 value 非空时追加 ILike，空字符串跳过。
	//
	// 例：tx.Model(&User{}).ILikeStr("name", "abc").Find(&rows)
	//
	// SQL（PostgreSQL）：
	//
	//	SELECT * FROM users WHERE name ILIKE '%abc%'
	ILikeStr(column string, value string) Query

	// Where 追加原生 GORM 条件，query 可为 SQL 片段或 map。
	//
	// 例：tx.Model(&User{}).Where("age > ?", 18).Find(&rows)
	//
	// SQL：
	//
	//	SELECT * FROM users WHERE age > 18
	//
	// 例：tx.Model(&User{}).Where(map[string]any{"status": 1}).Find(&rows)
	//
	// SQL：
	//
	//	SELECT * FROM users WHERE status = 1
	Where(query interface{}, args ...interface{}) Query

	// Or 追加 OR 条件，与前面条件为或关系。
	//
	// 例：tx.Model(&User{}).Equal("status", 1).Or("status = ?", 2).Find(&rows)
	//
	// SQL：
	//
	//	SELECT * FROM users WHERE status = 1 OR status = 2
	Or(query interface{}, args ...interface{}) Query

	// Not 对条件取反（NOT）。
	//
	// 例：tx.Model(&User{}).Not("status = ?", 0).Find(&rows)
	//
	// SQL：
	//
	//	SELECT * FROM users WHERE NOT (status = 0)
	Not(query interface{}, args ...interface{}) Query
}

// QueryExec 排序、分页、查询执行与更新。
//
// 读操作：Count / Find / First / Take。
// 写操作：Update / Updates / Omit（配合 Updates 或读操作）。
//
// 典型用法：
//
//	var rows []User
//	err := tx.Model(&User{}).Equal("enabled", 1).
//	    Order("id desc").Offset(10).Limit(10).Find(&rows)
//
// SQL：
//
//	SELECT * FROM users
//	WHERE enabled = 1
//	ORDER BY id DESC
//	LIMIT 10 OFFSET 10
type QueryExec interface {
	// Select 指定参与本次操作的列，行为取决于后续链式方法：
	//
	//   - 接 Find/First/Take：只查询列出的字段（列名或 struct 字段名）。
	//   - 接 Updates：只更新列出的字段（白名单），不是先 SELECT 再 UPDATE。
	//
	// 例（只查部分列）：
	//
	//	tx.Model(&User{}).Select("id", "name").Find(&rows)
	//
	// SQL：
	//
	//	SELECT id, name FROM users
	//
	// 例（只更新 status；row.Id = 5, row.Status = 3）：
	//
	//	tx.Model(&row).Select("Status").Updates(&row)
	//
	// SQL：
	//
	//	UPDATE users SET status = 3 WHERE id = 5
	Select(query interface{}, args ...interface{}) Query

	// Omit 指定本次操作忽略的列（黑名单），可接 Updates / Find 等。
	//
	// 例（更新除 name 外的非零字段；row.Id = 5, row.Status = 2, row.Name = "x"）：
	//
	//	tx.Model(&row).Omit("Name").Updates(&row)
	//
	// SQL：
	//
	//	UPDATE users SET status = 2 WHERE id = 5
	//
	// 例（查询时排除大字段）：
	//
	//	tx.Model(&User{}).Omit("config").Find(&rows)
	//
	// SQL：
	//
	//	SELECT id, name, code, ... /* 不含 config */ FROM users
	Omit(columns ...string) Query

	// Order 追加排序，可多次调用或传入逗号分隔表达式。
	// Beego 兼容："id"→升序，"-Ctime"→降序；也支持 GORM 写法 "id desc"。
	//
	// 例：tx.Model(&User{}).Order("-Ctime").Order("id").Find(&rows)
	//
	// SQL：
	//
	//	SELECT * FROM users ORDER BY id DESC, name ASC
	Order(value interface{}) Query

	// Offset 跳过前 n 条记录（分页）。
	//
	// 例：tx.Model(&User{}).Offset(20).Limit(10).Find(&rows)
	//
	// SQL：
	//
	//	SELECT * FROM users LIMIT 10 OFFSET 20
	Offset(offset int) Query

	// Limit 限制返回条数。
	//
	// 例：tx.Model(&User{}).Limit(5).Find(&rows)
	//
	// SQL：
	//
	//	SELECT * FROM users LIMIT 5
	Limit(limit int) Query

	// Count 统计当前条件下记录数，不加载行数据。
	//
	// 例：n, err := tx.Model(&User{}).Equal("enabled", 1).Count()
	//
	// SQL：
	//
	//	SELECT count(*) FROM users WHERE enabled = 1
	Count() (int64, error)

	// Find 查询多条，结果写入 dest（通常为 *[]Model）。
	// 无匹配记录时返回 nil，dest 为空 slice，不报错。
	//
	// conds 可选：额外追加 Where 条件，一般通过链式 Equal/Where 已足够。
	//
	// 例：
	//
	//	tx.Model(&User{}).Equal("enabled", 1).Find(&rows)
	//
	// SQL：
	//
	//	SELECT * FROM users WHERE enabled = 1
	Find(dest interface{}, conds ...interface{}) error

	// First 查询符合条件的第一条记录，写入 dest（通常为 *Model）。
	//
	// 按主键或 Order 决定“第一条”；无匹配时返回 gorm.ErrRecordNotFound。
	// 与 Take 的区别：First 在无 Order 时倾向按主键升序取第一条。
	//
	// 例：
	//
	//	tx.Model(&User{}).Equal("id", 1).First(&row)
	//
	// SQL：
	//
	//	SELECT * FROM users WHERE id = 1 ORDER BY users.id ASC LIMIT 1
	First(dest interface{}, conds ...interface{}) error

	// Take 取一条记录，无匹配时返回 gorm.ErrRecordNotFound。
	// 与 First 类似，但不保证排序语义，任意匹配一行即可。
	//
	// 例：
	//
	//	tx.Model(&User{}).Equal("code", "c1").Take(&row)
	//
	// SQL：
	//
	//	SELECT * FROM users WHERE code = 'c1' LIMIT 1
	Take(dest interface{}, conds ...interface{}) error

	// Updates 按当前 Model/条件更新记录。
	//
	// 传入 struct 时：
	//   - 不更新主键；
	//   - 默认跳过零值字段（0、""、false、nil 指针等）；
	//   - 若需只更新部分字段，先链式 Select 指定列（白名单）。
	//
	// 传入 map[string]any 时：只更新 map 中的键，零值也会写入。
	//
	// 注意：若 row 由 First 加载且未 Select/Omit，Updates(&row) 会尝试更新
	// struct 中所有非零字段；仅改一列时可 Update、Select 或传 map。
	//
	// 例（只更新 status；row.Id = 5, row.Status = 2）：
	//
	//	tx.Model(&row).Select("Status").Updates(&row)
	//
	// SQL：
	//
	//	UPDATE users SET status = 2 WHERE id = 5
	//
	// 例（map 更新，含零值）：
	//
	//	tx.Model(&User{}).Equal("id", 1).Updates(map[string]any{"enabled": 0})
	//
	// SQL：
	//
	//	UPDATE users SET enabled = 0 WHERE id = 1
	Updates(value interface{}) error

	// Update 更新单列，column 为列名或 struct 字段名，value 为新值。
	// 零值也会写入（与 Updates(struct) 不同）。
	//
	// 例（status 自增后写回；row.Id = 5, row.Status = 3）：
	//
	//	tx.Model(&row).Update("status", row.Status)
	//
	// SQL：
	//
	//	UPDATE users SET status = 3 WHERE id = 5
	//
	// 例（按条件更新）：
	//
	//	tx.Model(&User{}).Equal("id", 1).Update("enabled", 0)
	//
	// SQL：
	//
	//	UPDATE users SET enabled = 0 WHERE id = 1
	Update(column string, value any) error
}

// Query 由 QueryCond 与 QueryExec 组成，对外链式用法不变。
//
// 通过 Session.Model 或 Tx.Model 获得：
//
//	q := tx.Model(&User{})          // 指定表/模型
//	q = q.Equal("id", 1)            // 条件
//	err := q.First(&row)            // 执行
//
// SQL：
//
//	SELECT * FROM users WHERE id = 1 ORDER BY users.id ASC LIMIT 1
//
// 或一步写完：
//
//	err := tx.Model(&User{}).Equal("id", 1).First(&row)
type Query interface {
	QueryCond
	QueryExec
}
