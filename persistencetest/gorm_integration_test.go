package persistencetest

import (
	"testing"

	"github.com/lishimeng/app-starter/persistence"
)

func TestGormSqliteCRUD(t *testing.T) {
	assertNoErr(t, persistence.Install())
	registerTestModels()
	assertNoErr(t, persistence.RegisterDatabase(sqliteConfig(t, true)))

	session := persistence.GetSession(persistence.DefaultAlias)
	if session == nil {
		t.Fatal("session nil")
	}

	assertNoErr(t, session.Transaction(func(tx persistence.Tx) error {
		return tx.Create(&TestRecord{Name: "alpha", Status: 1})
	}))

	var count int64
	assertNoErr(t, session.Transaction(func(tx persistence.Tx) error {
		var err error
		count, err = tx.Model(&TestRecord{}).Where("status = ?", 1).Count()
		return err
	}))
	if count != 1 {
		t.Fatalf("count: got %d want 1", count)
	}

	var rows []TestRecord
	assertNoErr(t, session.Transaction(func(tx persistence.Tx) error {
		return tx.Model(&TestRecord{}).Where("name = ?", "alpha").Find(&rows)
	}))
	if len(rows) != 1 || rows[0].Name != "alpha" {
		t.Fatalf("unexpected rows: %+v", rows)
	}

	rows[0].Status = 2
	assertNoErr(t, session.Transaction(func(tx persistence.Tx) error {
		return tx.Model(&rows[0]).Select("Status").Updates(&rows[0])
	}))

	assertNoErr(t, session.Transaction(func(tx persistence.Tx) error {
		return tx.Delete(&rows[0])
	}))

	assertNoErr(t, session.Transaction(func(tx persistence.Tx) error {
		var err error
		count, err = tx.Model(&TestRecord{}).Count()
		return err
	}))
	if count != 0 {
		t.Fatalf("count after delete: got %d want 0", count)
	}
}
