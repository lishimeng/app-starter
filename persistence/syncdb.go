package persistence

import (
	"fmt"
	"log"

	gormdb "gorm.io/gorm"
)

// SyncOptions controls lightweight schema sync (Beego RunSyncdb semantics).
type SyncOptions struct {
	Force   bool // drop and recreate tables before sync (data loss)
	Verbose bool // log each DDL action
}

func syncLog(opts SyncOptions, format string, args ...any) {
	if opts.Verbose {
		log.Printf("syncdb: "+format, args...)
	}
}

// SyncDB creates missing tables, adds missing columns and indexes.
// It does not alter or drop existing columns or indexes.
func SyncDB(db *gormdb.DB, opts SyncOptions, models ...any) error {
	if db == nil {
		return fmt.Errorf("persistence: syncdb: db nil")
	}
	if len(models) == 0 {
		return nil
	}
	m := db.Migrator()
	for _, model := range models {
		if model == nil {
			continue
		}
		if err := syncModel(m, db, opts, model); err != nil {
			return err
		}
	}
	return nil
}

func syncModel(m gormdb.Migrator, db *gormdb.DB, opts SyncOptions, model any) error {
	stmt := &gormdb.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return err
	}
	if stmt.Schema == nil {
		return fmt.Errorf("persistence: syncdb: schema nil for %T", model)
	}
	table := stmt.Schema.Table

	if opts.Force && m.HasTable(model) {
		syncLog(opts, "drop table %s", table)
		if err := m.DropTable(model); err != nil {
			return err
		}
	}

	if !m.HasTable(model) {
		syncLog(opts, "create table %s", table)
		return m.CreateTable(model)
	}

	for _, dbName := range stmt.Schema.DBNames {
		field := stmt.Schema.FieldsByDBName[dbName]
		if field.IgnoreMigration {
			continue
		}
		if !m.HasColumn(model, dbName) {
			syncLog(opts, "add column %s.%s", table, dbName)
			if err := m.AddColumn(model, dbName); err != nil {
				return err
			}
		}
	}

	for _, idx := range stmt.Schema.ParseIndexes() {
		if !m.HasIndex(model, idx.Name) {
			syncLog(opts, "create index %s on %s", idx.Name, table)
			if err := m.CreateIndex(model, idx.Name); err != nil {
				return err
			}
		}
	}

	return nil
}

// RunSyncDB runs SyncDB on a registered database alias (Beego orm.RunSyncdb equivalent).
func RunSyncDB(alias string, opts SyncOptions, models ...any) error {
	c := getConnector()
	if c == nil {
		return fmt.Errorf("persistence: no connector registered; call Install() and register dialectors")
	}
	if alias == "" {
		alias = DefaultAlias
	}
	return c.Migrate(alias, opts, models...)
}
