package store

import (
	"database/sql"
	"fmt"
	"strings"

	"gooseboard/internal/schema"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Migrate(entities map[string]schema.Entity) error {
	for _, e := range entities {
		cols := []string{"id TEXT PRIMARY KEY"}
		for _, f := range e.Fields {
			cols = append(cols, fmt.Sprintf("%s %s", f.Name, sqlType(f.Type)))
		}
		stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", e.Name, strings.Join(cols, ", "))
		if _, err := s.db.Exec(stmt); err != nil {
			return fmt.Errorf("migrating %s: %w", e.Name, err)
		}
	}
	return nil
}

func sqlType(t schema.FieldType) string {
	switch t {
	case schema.FieldInt:
		return "INTEGER"
	case schema.FieldBool:
		return "INTEGER"
	default:
		return "TEXT"
	}
}

func (s *SQLiteStore) Create(entity string, data map[string]any) (string, error) {
	id := uuid.New().String()

	cols := make([]string, 0, len(data)+1)
	placeholders := make([]string, 0, len(data)+1)
	args := make([]any, 0, len(data)+1)

	cols = append(cols, "id")
	placeholders = append(placeholders, "?")
	args = append(args, id)

	for col, val := range data {
		cols = append(cols, col)
		placeholders = append(placeholders, "?")
		args = append(args, val)
	}

	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", entity, strings.Join(cols, ", "), strings.Join(placeholders, ", "))
	if _, err := s.db.Exec(stmt, args...); err != nil {
		return "", fmt.Errorf("creating %s: %w", entity, err)
	}
	return id, nil
}

func (s *SQLiteStore) Get(entity, id string) (map[string]any, error) {
	stmt := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", entity)
	rows, err := s.db.Query(stmt, id)
	if err != nil {
		return nil, fmt.Errorf("getting %s: %w", entity, err)
	}
	defer rows.Close()

	results, err := scanRows(rows)
	if err != nil {
		return nil, fmt.Errorf("getting %s: %w", entity, err)
	}
	if len(results) == 0 {
		return nil, sql.ErrNoRows
	}
	return results[0], nil
}

func (s *SQLiteStore) List(entity string, filter map[string]any) ([]map[string]any, error) {
	stmt := fmt.Sprintf("SELECT * FROM %s", entity)

	var args []any
	if len(filter) > 0 {
		conds := make([]string, 0, len(filter))
		for col, val := range filter {
			conds = append(conds, fmt.Sprintf("%s = ?", col))
			args = append(args, val)
		}
		stmt += " WHERE " + strings.Join(conds, " AND ")
	}

	rows, err := s.db.Query(stmt, args...)
	if err != nil {
		return nil, fmt.Errorf("listing %s: %w", entity, err)
	}
	defer rows.Close()

	results, err := scanRows(rows)
	if err != nil {
		return nil, fmt.Errorf("listing %s: %w", entity, err)
	}
	return results, nil
}

func (s *SQLiteStore) Update(entity, id string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}

	sets := make([]string, 0, len(data))
	args := make([]any, 0, len(data)+1)
	for col, val := range data {
		sets = append(sets, fmt.Sprintf("%s = ?", col))
		args = append(args, val)
	}
	args = append(args, id)

	stmt := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", entity, strings.Join(sets, ", "))
	if _, err := s.db.Exec(stmt, args...); err != nil {
		return fmt.Errorf("updating %s: %w", entity, err)
	}
	return nil
}

func (s *SQLiteStore) Delete(entity, id string) error {
	stmt := fmt.Sprintf("DELETE FROM %s WHERE id = ?", entity)
	if _, err := s.db.Exec(stmt, id); err != nil {
		return fmt.Errorf("deleting %s: %w", entity, err)
	}
	return nil
}

func scanRows(rows *sql.Rows) ([]map[string]any, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := make([]map[string]any, 0)
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}

		row := make(map[string]any, len(cols))
		for i, col := range cols {
			row[col] = vals[i]
		}
		results = append(results, row)
	}
	return results, rows.Err()
}
