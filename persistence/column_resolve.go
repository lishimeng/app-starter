package persistence

import (
	"strings"

	gormdb "gorm.io/gorm"
)

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

// resolveOrder maps Beego OrderBy expr to GORM ORDER BY clause.
//   - "id" / "UserCode" → "<column> asc"
//   - "-Ctime" / "-ctime" → "<column> desc"
//   - "id desc" / "UserCode asc" → resolves column, keeps direction
func (q *gormQuery) resolveOrder(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return value
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if expr := q.resolveOrderExpr(strings.TrimSpace(part)); expr != "" {
			out = append(out, expr)
		}
	}
	return strings.Join(out, ", ")
}

func (q *gormQuery) resolveOrderExpr(expr string) string {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return expr
	}
	fields := strings.Fields(expr)
	if len(fields) >= 2 {
		dir := strings.ToLower(fields[len(fields)-1])
		if dir == "asc" || dir == "desc" {
			col := strings.Join(fields[:len(fields)-1], " ")
			return q.resolveColumn(col) + " " + dir
		}
	}
	if strings.HasPrefix(expr, "-") {
		col := strings.TrimSpace(expr[1:])
		if col == "" {
			return expr
		}
		return q.resolveColumn(col) + " desc"
	}
	return q.resolveColumn(expr) + " asc"
}
