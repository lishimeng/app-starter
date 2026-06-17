package persistence

import "strings"

// schemaSnapshot holds database catalog metadata loaded in batch for SyncDB diff.
type schemaSnapshot struct {
	tables  map[string]struct{}
	columns map[string]map[string]struct{}
	indexes map[string]map[string]struct{}
}

func newSchemaSnapshot() *schemaSnapshot {
	return &schemaSnapshot{
		tables:  make(map[string]struct{}),
		columns: make(map[string]map[string]struct{}),
		indexes: make(map[string]map[string]struct{}),
	}
}

func normIdent(name string) string {
	return strings.ToLower(name)
}

func (s *schemaSnapshot) HasTable(table string) bool {
	_, ok := s.tables[normIdent(table)]
	return ok
}

func (s *schemaSnapshot) HasColumn(table, column string) bool {
	cols, ok := s.columns[normIdent(table)]
	if !ok {
		return false
	}
	_, ok = cols[normIdent(column)]
	return ok
}

func (s *schemaSnapshot) HasIndex(table, index string) bool {
	idx, ok := s.indexes[normIdent(table)]
	if !ok {
		return false
	}
	_, ok = idx[normIdent(index)]
	return ok
}

func (s *schemaSnapshot) addTable(table string) {
	s.tables[normIdent(table)] = struct{}{}
}

func (s *schemaSnapshot) removeTable(table string) {
	delete(s.tables, normIdent(table))
	delete(s.columns, normIdent(table))
	delete(s.indexes, normIdent(table))
}

func (s *schemaSnapshot) addColumn(table, column string) {
	t := normIdent(table)
	if s.columns[t] == nil {
		s.columns[t] = make(map[string]struct{})
	}
	s.columns[t][normIdent(column)] = struct{}{}
}

func (s *schemaSnapshot) addIndex(table, index string) {
	t := normIdent(table)
	if s.indexes[t] == nil {
		s.indexes[t] = make(map[string]struct{})
	}
	s.indexes[t][normIdent(index)] = struct{}{}
}

func (s *schemaSnapshot) seedFromSchema(table string, dbNames []string, indexNames []string) {
	s.addTable(table)
	for _, col := range dbNames {
		s.addColumn(table, col)
	}
	for _, idx := range indexNames {
		s.addIndex(table, idx)
	}
}
