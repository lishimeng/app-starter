package persistence

func (q *gormQuery) Where(query interface{}, args ...interface{}) Query {
	if q == nil || q.db == nil {
		return q
	}
	return &gormQuery{db: q.db.Where(query, args...), model: q.model}
}

func (q *gormQuery) Or(query interface{}, args ...interface{}) Query {
	if q == nil || q.db == nil {
		return q
	}
	return &gormQuery{db: q.db.Or(query, args...), model: q.model}
}

func (q *gormQuery) Not(query interface{}, args ...interface{}) Query {
	if q == nil || q.db == nil {
		return q
	}
	return &gormQuery{db: q.db.Not(query, args...), model: q.model}
}

// Equal column = ?
func (q *gormQuery) Equal(column string, value any) Query {
	return q.Where(q.resolveColumn(column)+" = ?", value)
}

// NotEqual column <> ?
func (q *gormQuery) NotEqual(column string, value any) Query {
	return q.Where(q.resolveColumn(column)+" <> ?", value)
}

// In column IN ?
func (q *gormQuery) In(column string, values any) Query {
	return q.Where(q.resolveColumn(column)+" IN ?", values)
}

// Like column LIKE %value%
func (q *gormQuery) Like(column string, value string) Query {
	return q.Where(q.resolveColumn(column)+" LIKE ?", "%"+value+"%")
}

// LLike column LIKE value%（前缀匹配）
func (q *gormQuery) LLike(column string, value string) Query {
	return q.Where(q.resolveColumn(column)+" LIKE ?", value+"%")
}

// RLike column LIKE %value（后缀匹配）
func (q *gormQuery) RLike(column string, value string) Query {
	return q.Where(q.resolveColumn(column)+" LIKE ?", "%"+value)
}

// ILike column ILIKE %value%（不区分大小写，PostgreSQL）
func (q *gormQuery) ILike(column string, value string) Query {
	return q.Where(q.resolveColumn(column)+" ILIKE ?", "%"+value+"%")
}

// EqualStr value 非空时追加 Equal。
func (q *gormQuery) EqualStr(column string, value string) Query {
	if value == "" {
		return q
	}
	return q.Equal(column, value)
}

// LikeStr value 非空时追加 Like。
func (q *gormQuery) LikeStr(column string, value string) Query {
	if value == "" {
		return q
	}
	return q.Like(column, value)
}

// LLikeStr value 非空时追加 LLike。
func (q *gormQuery) LLikeStr(column string, value string) Query {
	if value == "" {
		return q
	}
	return q.LLike(column, value)
}

// RLikeStr value 非空时追加 RLike。
func (q *gormQuery) RLikeStr(column string, value string) Query {
	if value == "" {
		return q
	}
	return q.RLike(column, value)
}

// ILikeStr value 非空时追加 ILike。
func (q *gormQuery) ILikeStr(column string, value string) Query {
	if value == "" {
		return q
	}
	return q.ILike(column, value)
}
