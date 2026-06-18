package persistence

import gormdb "gorm.io/gorm"

// resolveColumn maps struct field name to gorm column tag when Model is set.
// Accepts DB column names as-is (e.g. user_code). Beego Filter-style field names work too.
func (q *gormQuery) resolveColumn(column string) string {
	if q == nil || column == "" || q.model == nil || q.db == nil {
		return column
	}
	stmt := &gormdb.Statement{DB: q.db}
	if err := stmt.Parse(q.model); err != nil || stmt.Schema == nil {
		return column
	}
	if field := stmt.Schema.LookUpField(column); field != nil && field.DBName != "" {
		return field.DBName
	}
	return column
}
