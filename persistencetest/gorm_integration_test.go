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
		count, err = tx.Model(&TestRecord{}).Equal("status", 1).Count()
		return err
	}))
	if count != 1 {
		t.Fatalf("count: got %d want 1", count)
	}

	var rows []TestRecord
	assertNoErr(t, session.Transaction(func(tx persistence.Tx) error {
		return tx.Model(&TestRecord{}).Equal("name", "alpha").Find(&rows)
	}))
	if len(rows) != 1 || rows[0].Name != "alpha" {
		t.Fatalf("unexpected rows: %+v", rows)
	}

	assertNoErr(t, session.Transaction(func(tx persistence.Tx) error {
		return tx.Model(&rows[0]).Update("status", 2)
	}))

	var updated TestRecord
	assertNoErr(t, session.Transaction(func(tx persistence.Tx) error {
		return tx.Model(&TestRecord{}).Equal("name", "alpha").First(&updated)
	}))
	if updated.Status != 2 {
		t.Fatalf("status after Update: got %d want 2", updated.Status)
	}

	updated.Name = "should-not-save"
	updated.Status = 5
	assertNoErr(t, session.Transaction(func(tx persistence.Tx) error {
		return tx.Model(&updated).Omit("Name").Updates(&updated)
	}))
	var omitted TestRecord
	assertNoErr(t, session.Transaction(func(tx persistence.Tx) error {
		return tx.Model(&TestRecord{}).Equal("name", "alpha").First(&omitted)
	}))
	if omitted.Name != "alpha" || omitted.Status != 5 {
		t.Fatalf("after Omit: got %+v want name=alpha status=5", omitted)
	}

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
