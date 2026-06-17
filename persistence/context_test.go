package persistence

import "testing"

type mockQuery struct {
	filters int
}

func (q *mockQuery) Where(query interface{}, args ...interface{}) Query {
	q.filters++
	return q
}

func (q *mockQuery) Equal(column string, value any) Query       { q.filters++; return q }
func (q *mockQuery) NotEqual(column string, value any) Query    { q.filters++; return q }
func (q *mockQuery) In(column string, values any) Query         { q.filters++; return q }
func (q *mockQuery) Like(column string, value string) Query     { q.filters++; return q }
func (q *mockQuery) LLike(column string, value string) Query    { q.filters++; return q }
func (q *mockQuery) RLike(column string, value string) Query    { q.filters++; return q }
func (q *mockQuery) ILike(column string, value string) Query    { q.filters++; return q }
func (q *mockQuery) EqualStr(column string, value string) Query { if value != "" { q.filters++ }; return q }
func (q *mockQuery) LikeStr(column string, value string) Query  { if value != "" { q.filters++ }; return q }
func (q *mockQuery) LLikeStr(column string, value string) Query { if value != "" { q.filters++ }; return q }
func (q *mockQuery) RLikeStr(column string, value string) Query { if value != "" { q.filters++ }; return q }
func (q *mockQuery) ILikeStr(column string, value string) Query { if value != "" { q.filters++ }; return q }

func (q *mockQuery) Or(query interface{}, args ...interface{}) Query     { return q }
func (q *mockQuery) Not(query interface{}, args ...interface{}) Query    { return q }
func (q *mockQuery) Select(query interface{}, args ...interface{}) Query { return q }
func (q *mockQuery) Omit(columns ...string) Query                        { return q }
func (q *mockQuery) Order(value interface{}) Query                       { return q }
func (q *mockQuery) Offset(offset int) Query                             { return q }
func (q *mockQuery) Limit(limit int) Query                               { return q }
func (q *mockQuery) Count() (int64, error)                               { return 0, nil }
func (q *mockQuery) Find(dest interface{}, conds ...interface{}) error     { return nil }
func (q *mockQuery) First(dest interface{}, conds ...interface{}) error  { return nil }
func (q *mockQuery) Take(dest interface{}, conds ...interface{}) error    { return nil }
func (q *mockQuery) Updates(value interface{}) error                     { return nil }
func (q *mockQuery) Update(column string, value any) error               { return nil }

type mockTx struct {
	created bool
}

func (t *mockTx) Model(value interface{}) Query { return &mockQuery{} }
func (t *mockTx) Create(value interface{}) error {
	t.created = true
	return nil
}
func (t *mockTx) Save(value interface{}) error                         { return nil }
func (t *mockTx) Delete(value interface{}, conds ...interface{}) error { return nil }
func (t *mockTx) First(dest interface{}, conds ...interface{}) error   { return nil }
func (t *mockTx) Raw(sql string, values ...interface{}) Query          { return &mockQuery{} }

type mockSession struct {
	debug bool
	tx    *mockTx
}

func (s *mockSession) Transaction(fn func(Tx) error) error {
	if fn == nil {
		return nil
	}
	s.tx = &mockTx{}
	return fn(s.tx)
}

func (s *mockSession) Model(value interface{}) Query { return &mockQuery{} }
func (s *mockSession) SetDebug(enable bool)          { s.debug = enable }
func (s *mockSession) Alias() string                 { return "mock" }

func TestOrmContextFacade(t *testing.T) {
	s := &mockSession{}
	ctx := WrapSession(s)

	if ctx.Model(struct{}{}) == nil {
		t.Fatal("expected query from facade")
	}

	ctx.SetLogEnable(true)
	if !s.debug {
		t.Fatal("expected debug enabled on session")
	}

	err := ctx.Transaction(func(tx TxContext) error {
		if tx.Model(struct{}{}) == nil {
			t.Fatal("expected query from tx facade")
		}
		return tx.Create(struct{}{})
	})
	if err != nil {
		t.Fatal(err)
	}
	if s.tx == nil || !s.tx.created {
		t.Fatal("expected transaction to run create")
	}
}

func TestWrapTxExposesTx(t *testing.T) {
	tx := &mockTx{}
	wrapped := WrapTx(tx)
	if wrapped.Tx == nil {
		t.Fatal("expected tx on context")
	}
	if err := wrapped.Create(struct{}{}); err != nil {
		t.Fatal(err)
	}
	if !tx.created {
		t.Fatal("expected create through facade")
	}
}
