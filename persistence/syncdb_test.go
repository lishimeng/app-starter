package persistence

import (
	"fmt"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	gormdb "gorm.io/gorm"
)

type syncFullRecord struct {
	Id     int    `gorm:"primaryKey;autoIncrement;column:id"`
	Name   string `gorm:"column:name"`
	Status int    `gorm:"column:status;default:0"`
}

func (syncFullRecord) TableName() string { return "sync_full_record" }

type syncIndexedRecord struct {
	Id   int    `gorm:"primaryKey;autoIncrement;column:id"`
	Code string `gorm:"column:code;uniqueIndex:idx_sync_code"`
	Name string `gorm:"column:name"`
}

func (syncIndexedRecord) TableName() string { return "sync_indexed_record" }

type syncStringStatusRecord struct {
	Id     int    `gorm:"primaryKey;column:id"`
	Status string `gorm:"column:status"`
}

func (syncStringStatusRecord) TableName() string { return "sync_alter_record" }

func openSyncDBTestDB(t *testing.T) *gormdb.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:syncdb_%s?mode=memory&cache=shared", t.Name())
	db, err := gormdb.Open(sqlite.Open(dsn), &gormdb.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestSyncDB_CreateTable(t *testing.T) {
	db := openSyncDBTestDB(t)
	m := db.Migrator()

	if err := SyncDB(db, SyncOptions{}, &syncFullRecord{}); err != nil {
		t.Fatal(err)
	}
	if !m.HasTable(&syncFullRecord{}) {
		t.Fatal("expected table")
	}
	for _, col := range []string{"id", "name", "status"} {
		if !m.HasColumn(&syncFullRecord{}, col) {
			t.Fatalf("expected column %s", col)
		}
	}
}

func TestSyncDB_AddColumn(t *testing.T) {
	db := openSyncDBTestDB(t)
	m := db.Migrator()

	if err := db.Exec(`CREATE TABLE sync_full_record (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT
	)`).Error; err != nil {
		t.Fatal(err)
	}
	if m.HasColumn(&syncFullRecord{}, "status") {
		t.Fatal("status should not exist before sync")
	}

	if err := SyncDB(db, SyncOptions{}, &syncFullRecord{}); err != nil {
		t.Fatal(err)
	}
	if !m.HasColumn(&syncFullRecord{}, "status") {
		t.Fatal("expected status column added")
	}
}

func TestSyncDB_AddIndex(t *testing.T) {
	db := openSyncDBTestDB(t)
	m := db.Migrator()

	if err := db.Exec(`CREATE TABLE sync_indexed_record (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT,
		name TEXT
	)`).Error; err != nil {
		t.Fatal(err)
	}

	model := &syncIndexedRecord{}
	if err := SyncDB(db, SyncOptions{}, model); err != nil {
		t.Fatal(err)
	}

	stmt := &gormdb.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, idx := range stmt.Schema.ParseIndexes() {
		if m.HasIndex(model, idx.Name) {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected at least one index created")
	}
}

func TestSyncDB_Force(t *testing.T) {
	db := openSyncDBTestDB(t)

	if err := SyncDB(db, SyncOptions{}, &syncFullRecord{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&syncFullRecord{Name: "keep"}).Error; err != nil {
		t.Fatal(err)
	}

	var count int64
	if err := db.Model(&syncFullRecord{}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("count before force: got %d want 1", count)
	}

	if err := SyncDB(db, SyncOptions{Force: true}, &syncFullRecord{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&syncFullRecord{}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("count after force: got %d want 0", count)
	}
	if !db.Migrator().HasTable(&syncFullRecord{}) {
		t.Fatal("table should exist after force recreate")
	}
}

func TestSyncDB_NoAlterColumn(t *testing.T) {
	db := openSyncDBTestDB(t)

	if err := db.Exec(`CREATE TABLE sync_alter_record (
		id INTEGER PRIMARY KEY,
		status INTEGER NOT NULL DEFAULT 0
	)`).Error; err != nil {
		t.Fatal(err)
	}

	if err := SyncDB(db, SyncOptions{}, &syncStringStatusRecord{}); err != nil {
		t.Fatal(err)
	}

	type pragmaCol struct {
		Name string `gorm:"column:name"`
		Type string `gorm:"column:type"`
	}
	var cols []pragmaCol
	if err := db.Raw("PRAGMA table_info('sync_alter_record')").Scan(&cols).Error; err != nil {
		t.Fatal(err)
	}
	for _, col := range cols {
		if col.Name == "status" && !strings.Contains(strings.ToUpper(col.Type), "INT") {
			t.Fatalf("status column type should remain integer-like, got %q", col.Type)
		}
	}
}

func TestSchemaSnapshot_LoadSqlite(t *testing.T) {
	db := openSyncDBTestDB(t)

	if err := db.Exec(`CREATE TABLE snap_meta_a (id INTEGER PRIMARY KEY, name TEXT)`).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`CREATE UNIQUE INDEX idx_snap_meta_a_code ON snap_meta_a (name)`).Error; err != nil {
		t.Fatal(err)
	}

	snap, err := loadSchemaSnapshot(db, []string{"snap_meta_a"})
	if err != nil {
		t.Fatal(err)
	}
	if !snap.HasTable("snap_meta_a") {
		t.Fatal("expected table in snapshot")
	}
	if !snap.HasColumn("snap_meta_a", "name") {
		t.Fatal("expected column in snapshot")
	}
	if !snap.HasIndex("snap_meta_a", "idx_snap_meta_a_code") {
		t.Fatal("expected index in snapshot")
	}
}

func ensureSqliteForTest(t *testing.T) {
	t.Helper()
	if err := Install(); err != nil {
		t.Fatal(err)
	}
	RegisterDialector(DriverSqlite.Name, func(opts OpenOptions) gormdb.Dialector {
		return sqlite.Open(opts.DSN)
	})
}

func TestRunSyncDB(t *testing.T) {
	ensureSqliteForTest(t)
	alias := "sync_" + strings.ReplaceAll(t.Name(), "/", "_")
	dsn := fmt.Sprintf("file:run_syncdb_%s?mode=memory&cache=shared", t.Name())
	if err := RegisterDatabase(BaseConfig{
		AliasName:  alias,
		Driver:     DriverSqlite,
		DataSource: dsn,
		InitDb:     false,
		Models:     []any{&syncFullRecord{}},
	}); err != nil {
		t.Fatal(err)
	}

	if err := RunSyncDB(alias, SyncOptions{Verbose: true}, &syncFullRecord{}); err != nil {
		t.Fatal(err)
	}

	c := getConnector().(*gormConnector)
	db, err := c.db(alias)
	if err != nil {
		t.Fatal(err)
	}
	if !db.Migrator().HasTable(&syncFullRecord{}) {
		t.Fatal("RunSyncDB should create table")
	}
}

func TestRegisterDatabaseInitDbUsesSyncDB(t *testing.T) {
	ensureSqliteForTest(t)
	alias := "initdb_" + strings.ReplaceAll(t.Name(), "/", "_")
	dsn := fmt.Sprintf("file:initdb_%s?mode=memory&cache=shared", t.Name())
	if err := RegisterDatabase(BaseConfig{
		AliasName:  alias,
		Driver:     DriverSqlite,
		DataSource: dsn,
		InitDb:     true,
		Models:     []any{&syncFullRecord{}},
	}); err != nil {
		t.Fatal(err)
	}

	c := getConnector().(*gormConnector)
	db, err := c.db(alias)
	if err != nil {
		t.Fatal(err)
	}
	if !db.Migrator().HasTable(&syncFullRecord{}) {
		t.Fatal("InitDb should run SyncDB and create table")
	}
}
