package persistence

import "testing"

type mockQuery struct {
	filters int
}

func (q *mockQuery) Filter(expr string, args ...any) Query {
	q.filters++
	return q
}

func (q *mockQuery) FilterCond(cond Condition) Query { return q }
func (q *mockQuery) OrderBy(expr ...string) Query   { return q }
func (q *mockQuery) Offset(n int) Query               { return q }
func (q *mockQuery) Limit(n int) Query                { return q }
func (q *mockQuery) Count() (int64, error)            { return 0, nil }
func (q *mockQuery) All(dest any) (int64, error)      { return 0, nil }
func (q *mockQuery) One(dest any) error               { return nil }

type mockTx struct {
	inserted bool
}

func (t *mockTx) Query(model any) Query { return &mockQuery{} }
func (t *mockTx) Insert(model any) error {
	t.inserted = true
	return nil
}
func (t *mockTx) Update(model any, cols ...string) error { return nil }
func (t *mockTx) Delete(model any, cols ...string) error { return nil }
func (t *mockTx) Get(model any, cols ...string) error    { return nil }
func (t *mockTx) Raw(sql string, args ...any) Query      { return &mockQuery{} }

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

func (s *mockSession) Query(model any) Query { return &mockQuery{} }
func (s *mockSession) SetDebug(enable bool)  { s.debug = enable }
func (s *mockSession) Alias() string         { return "mock" }

func TestOrmContextFacade(t *testing.T) {
	s := &mockSession{}
	ctx := WrapSession(s)

	if ctx.Query(struct{}{}) == nil {
		t.Fatal("expected query from facade")
	}

	ctx.SetLogEnable(true)
	if !s.debug {
		t.Fatal("expected debug enabled on session")
	}

	err := ctx.Transaction(func(tx TxContext) error {
		if tx.Query(struct{}{}) == nil {
			t.Fatal("expected query from tx facade")
		}
		return tx.Insert(struct{}{})
	})
	if err != nil {
		t.Fatal(err)
	}
	if s.tx == nil || !s.tx.inserted {
		t.Fatal("expected transaction to run insert")
	}
}

func TestWrapTxExposesTx(t *testing.T) {
	tx := &mockTx{}
	wrapped := WrapTx(tx)
	if wrapped.Tx == nil {
		t.Fatal("expected tx on context")
	}
	if err := wrapped.Insert(struct{}{}); err != nil {
		t.Fatal(err)
	}
	if !tx.inserted {
		t.Fatal("expected insert through facade")
	}
}
