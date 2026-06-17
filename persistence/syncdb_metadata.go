package persistence

import (
	"fmt"

	gormdb "gorm.io/gorm"
)

type tableNameRow struct {
	TableName string `gorm:"column:table_name"`
}

type tableColumnRow struct {
	TableName  string `gorm:"column:table_name"`
	ColumnName string `gorm:"column:column_name"`
}

type sqliteMasterRow struct {
	Type    string `gorm:"column:type"`
	Name    string `gorm:"column:name"`
	TblName string `gorm:"column:tbl_name"`
}

type pragmaColumnRow struct {
	Name string `gorm:"column:name"`
}

func loadSchemaSnapshot(db *gormdb.DB, tables []string) (*schemaSnapshot, error) {
	if db == nil {
		return nil, fmt.Errorf("persistence: syncdb: db nil")
	}
	switch db.Dialector.Name() {
	case "postgres":
		return loadPostgresSnapshot(db, tables)
	case "mysql":
		return loadMysqlSnapshot(db, tables)
	case "sqlite", "sqlite3":
		return loadSqliteSnapshot(db, tables)
	default:
		return loadInformationSchemaSnapshot(db, tables)
	}
}

func loadPostgresSnapshot(db *gormdb.DB, tables []string) (*schemaSnapshot, error) {
	snap := newSchemaSnapshot()

	tableSQL := `SELECT table_name FROM information_schema.tables
		WHERE table_schema = current_schema() AND table_type = 'BASE TABLE'`
	var tableRows []tableNameRow
	if len(tables) > 0 {
		if err := db.Raw(tableSQL+` AND table_name IN ?`, tables).Scan(&tableRows).Error; err != nil {
			return nil, err
		}
	} else if err := db.Raw(tableSQL).Scan(&tableRows).Error; err != nil {
		return nil, err
	}
	for _, r := range tableRows {
		snap.addTable(r.TableName)
	}

	columnSQL := `SELECT table_name, column_name FROM information_schema.columns
		WHERE table_schema = current_schema()`
	var colRows []tableColumnRow
	if len(tables) > 0 {
		if err := db.Raw(columnSQL+` AND table_name IN ?`, tables).Scan(&colRows).Error; err != nil {
			return nil, err
		}
	} else if err := db.Raw(columnSQL).Scan(&colRows).Error; err != nil {
		return nil, err
	}
	for _, r := range colRows {
		snap.addColumn(r.TableName, r.ColumnName)
	}

	indexSQL := `SELECT tablename, indexname FROM pg_indexes WHERE schemaname = current_schema()`
	var idxRows []struct {
		TableName string `gorm:"column:tablename"`
		IndexName string `gorm:"column:indexname"`
	}
	if len(tables) > 0 {
		if err := db.Raw(indexSQL+` AND tablename IN ?`, tables).Scan(&idxRows).Error; err != nil {
			return nil, err
		}
	} else if err := db.Raw(indexSQL).Scan(&idxRows).Error; err != nil {
		return nil, err
	}
	for _, r := range idxRows {
		snap.addIndex(r.TableName, r.IndexName)
	}

	return snap, nil
}

func loadMysqlSnapshot(db *gormdb.DB, tables []string) (*schemaSnapshot, error) {
	snap := newSchemaSnapshot()

	tableSQL := `SELECT TABLE_NAME AS table_name FROM information_schema.tables
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_TYPE = 'BASE TABLE'`
	var tableRows []tableNameRow
	if len(tables) > 0 {
		if err := db.Raw(tableSQL+` AND TABLE_NAME IN ?`, tables).Scan(&tableRows).Error; err != nil {
			return nil, err
		}
	} else if err := db.Raw(tableSQL).Scan(&tableRows).Error; err != nil {
		return nil, err
	}
	for _, r := range tableRows {
		snap.addTable(r.TableName)
	}

	columnSQL := `SELECT TABLE_NAME AS table_name, COLUMN_NAME AS column_name FROM information_schema.columns
		WHERE TABLE_SCHEMA = DATABASE()`
	var colRows []tableColumnRow
	if len(tables) > 0 {
		if err := db.Raw(columnSQL+` AND TABLE_NAME IN ?`, tables).Scan(&colRows).Error; err != nil {
			return nil, err
		}
	} else if err := db.Raw(columnSQL).Scan(&colRows).Error; err != nil {
		return nil, err
	}
	for _, r := range colRows {
		snap.addColumn(r.TableName, r.ColumnName)
	}

	indexSQL := `SELECT DISTINCT TABLE_NAME AS table_name, INDEX_NAME AS index_name FROM information_schema.statistics
		WHERE TABLE_SCHEMA = DATABASE()`
	var idxRows []struct {
		TableName string `gorm:"column:table_name"`
		IndexName string `gorm:"column:index_name"`
	}
	if len(tables) > 0 {
		if err := db.Raw(indexSQL+` AND TABLE_NAME IN ?`, tables).Scan(&idxRows).Error; err != nil {
			return nil, err
		}
	} else if err := db.Raw(indexSQL).Scan(&idxRows).Error; err != nil {
		return nil, err
	}
	for _, r := range idxRows {
		snap.addIndex(r.TableName, r.IndexName)
	}

	return snap, nil
}

func loadSqliteSnapshot(db *gormdb.DB, tables []string) (*schemaSnapshot, error) {
	snap := newSchemaSnapshot()

	var masterRows []sqliteMasterRow
	if err := db.Raw(`SELECT type, name, tbl_name FROM sqlite_master
		WHERE type IN ('table', 'index') AND name NOT LIKE 'sqlite_%'`).Scan(&masterRows).Error; err != nil {
		return nil, err
	}

	tableSet := make(map[string]struct{})
	if len(tables) > 0 {
		for _, t := range tables {
			tableSet[normIdent(t)] = struct{}{}
		}
	}

	existingTables := make([]string, 0)
	for _, r := range masterRows {
		switch r.Type {
		case "table":
			if len(tableSet) > 0 {
				if _, ok := tableSet[normIdent(r.Name)]; !ok {
					continue
				}
			}
			snap.addTable(r.Name)
			existingTables = append(existingTables, r.Name)
		case "index":
			if r.TblName == "" {
				continue
			}
			if len(tableSet) > 0 {
				if _, ok := tableSet[normIdent(r.TblName)]; !ok {
					continue
				}
			}
			snap.addIndex(r.TblName, r.Name)
		}
	}

	columnTables := existingTables
	if len(tables) > 0 {
		columnTables = make([]string, 0, len(tables))
		for _, t := range tables {
			if snap.HasTable(t) {
				columnTables = append(columnTables, t)
			}
		}
	}

	for _, table := range columnTables {
		var cols []pragmaColumnRow
		if err := db.Raw(fmt.Sprintf("PRAGMA table_info(%q)", table)).Scan(&cols).Error; err != nil {
			return nil, err
		}
		for _, c := range cols {
			snap.addColumn(table, c.Name)
		}
	}

	return snap, nil
}

func loadInformationSchemaSnapshot(db *gormdb.DB, tables []string) (*schemaSnapshot, error) {
	var dbName string
	if err := db.Raw("SELECT DATABASE()").Scan(&dbName).Error; err != nil {
		return loadSqliteSnapshot(db, tables)
	}
	if dbName == "" {
		return loadPostgresSnapshot(db, tables)
	}
	return loadMysqlSnapshot(db, tables)
}

func collectModelTables(db *gormdb.DB, models []any) ([]string, error) {
	names := make([]string, 0, len(models))
	seen := make(map[string]struct{})
	for _, model := range models {
		if model == nil {
			continue
		}
		stmt := &gormdb.Statement{DB: db}
		if err := stmt.Parse(model); err != nil {
			return nil, err
		}
		if stmt.Schema == nil {
			return nil, fmt.Errorf("persistence: syncdb: schema nil for %T", model)
		}
		table := stmt.Schema.Table
		if _, ok := seen[table]; ok {
			continue
		}
		seen[table] = struct{}{}
		names = append(names, table)
	}
	return names, nil
}

func indexNamesFromSchema(stmt *gormdb.Statement) []string {
	if stmt == nil || stmt.Schema == nil {
		return nil
	}
	indexes := stmt.Schema.ParseIndexes()
	names := make([]string, 0, len(indexes))
	for _, idx := range indexes {
		if idx.Name != "" {
			names = append(names, idx.Name)
		}
	}
	return names
}
