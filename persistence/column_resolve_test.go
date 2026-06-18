package persistence

import (
	"fmt"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	gormdb "gorm.io/gorm"
)

type columnResolveModel struct {
	UserCode string `gorm:"column:user_code"`
	ConnType string `gorm:"column:conn_type"`
}

func (columnResolveModel) TableName() string { return "col_resolve_test" }

func openColumnResolveDB(t *testing.T) *gormdb.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:col_resolve_%s?mode=memory&cache=shared", t.Name())
	db, err := gormdb.Open(sqlite.Open(dsn), &gormdb.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestResolveColumn_FieldAndDBName(t *testing.T) {
	db := openColumnResolveDB(t)
	model := &columnResolveModel{}
	q := wrapGormQuery(db.Model(model), model).(*gormQuery)

	if got := q.resolveColumn("UserCode"); got != "user_code" {
		t.Fatalf("UserCode: got %q want user_code", got)
	}
	if got := q.resolveColumn("ConnType"); got != "conn_type" {
		t.Fatalf("ConnType: got %q want conn_type", got)
	}
	if got := q.resolveColumn("user_code"); got != "user_code" {
		t.Fatalf("user_code: got %q want user_code", got)
	}
}

func TestEqual_ResolvesStructFieldName(t *testing.T) {
	db := openColumnResolveDB(t)
	model := &columnResolveModel{}

	gq := wrapGormQuery(db.Model(model), model).Equal("ConnType", "http").(*gormQuery)
	var dest []columnResolveModel
	stmt := gq.db.Session(&gormdb.Session{DryRun: true}).Find(&dest).Statement
	if !strings.Contains(stmt.SQL.String(), "conn_type") {
		t.Fatalf("expected conn_type in SQL, got %q", stmt.SQL.String())
	}
}

func TestUpdate_ResolvesStructFieldName(t *testing.T) {
	db := openColumnResolveDB(t)
	model := &columnResolveModel{UserCode: "old"}

	gq := wrapGormQuery(db.Session(&gormdb.Session{DryRun: true}).Model(model), model).(*gormQuery)
	_ = gq.Update("UserCode", "new")
	sql := gq.db.Statement.SQL.String()
	if !strings.Contains(sql, "user_code") {
		t.Fatalf("expected user_code in SQL, got %q", sql)
	}
}
