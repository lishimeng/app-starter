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
// Catalog metadata is loaded in batch; Migrator Has* is not used for diff.
func SyncDB(db *gormdb.DB, opts SyncOptions, models ...any) error {
	if db == nil {
		return fmt.Errorf("persistence: syncdb: db nil")
	}
	if len(models) == 0 {
		return nil
	}

	tableNames, err := collectModelTables(db, models)
	if err != nil {
		return err
	}

	snap, err := loadSchemaSnapshot(db, tableNames)
	if err != nil {
		return err
	}

	m := db.Migrator()
	for _, model := range models {
		if model == nil {
			continue
		}
		if err := syncModel(m, db, snap, opts, model); err != nil {
			return err
		}
	}
	return nil
}

func syncModel(m gormdb.Migrator, db *gormdb.DB, snap *schemaSnapshot, opts SyncOptions, model any) error {
	stmt := &gormdb.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return err
	}
	if stmt.Schema == nil {
		return fmt.Errorf("persistence: syncdb: schema nil for %T", model)
	}
	table := stmt.Schema.Table

	if opts.Force && snap.HasTable(table) {
		syncLog(opts, "drop table %s", table)
		if err := m.DropTable(model); err != nil {
			return err
		}
		snap.removeTable(table)
	}

	if !snap.HasTable(table) {
		syncLog(opts, "create table %s", table)
		if err := m.CreateTable(model); err != nil {
			return err
		}
		snap.seedFromSchema(table, stmt.Schema.DBNames, indexNamesFromSchema(stmt))
		return nil
	}

	for _, dbName := range stmt.Schema.DBNames {
		field := stmt.Schema.FieldsByDBName[dbName]
		if field.IgnoreMigration {
			continue
		}
		if snap.HasColumn(table, dbName) {
			continue
		}
		syncLog(opts, "add column %s.%s", table, dbName)
		if err := m.AddColumn(model, dbName); err != nil {
			return err
		}
		snap.addColumn(table, dbName)
	}

	for _, idx := range stmt.Schema.ParseIndexes() {
		if snap.HasIndex(table, idx.Name) {
			continue
		}
		syncLog(opts, "create index %s on %s", idx.Name, table)
		if err := m.CreateIndex(model, idx.Name); err != nil {
			return err
		}
		snap.addIndex(table, idx.Name)
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
