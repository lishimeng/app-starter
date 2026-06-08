package persistence

import (
	"fmt"
	"strings"
)

func fieldFilter(field string, value any) (expr string, args []any) {
	parts := strings.Split(field, "__")
	col := parts[0]
	if len(parts) == 1 {
		return col + " = ?", []any{value}
	}
	switch parts[1] {
	case "gt":
		return col + " > ?", []any{value}
	case "gte":
		return col + " >= ?", []any{value}
	case "lt":
		return col + " < ?", []any{value}
	case "lte":
		return col + " <= ?", []any{value}
	case "in":
		return col + " IN ?", []any{value}
	case "exact":
		return col + " = ?", []any{value}
	case "icontains":
		return col + " ILIKE ?", []any{"%" + fmt.Sprint(value) + "%"}
	case "contains":
		return col + " LIKE ?", []any{"%" + fmt.Sprint(value) + "%"}
	default:
		return col + " = ?", []any{value}
	}
}

func isSQLExpr(expr string) bool {
	return strings.Contains(expr, "?") ||
		strings.Contains(expr, "=") ||
		strings.Contains(expr, " ") ||
		strings.Contains(expr, ">") ||
		strings.Contains(expr, "<")
}

func applyOrderExprs(expr ...string) []string {
	orders := make([]string, 0, len(expr))
	for _, e := range expr {
		if e == "" {
			continue
		}
		if strings.HasPrefix(e, "-") {
			orders = append(orders, strings.TrimPrefix(e, "-")+" DESC")
			continue
		}
		orders = append(orders, e+" ASC")
	}
	return orders
}
