package pgxdb

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// bridgeConnect opens a PostgreSQL connection.
// Called from db.glsp as (bridge-connect url) → bridgeConnect(url).
func bridgeConnect(url string) (any, error) {
	return pgx.Connect(context.Background(), url)
}

// bridgeClose closes the connection.
func bridgeClose(conn any) error {
	return conn.(*pgx.Conn).Close(context.Background())
}

// bridgeQuery runs a SELECT and returns rows as []any where each element is map[string]any.
// Returns []any (not []map[string]any) so glisp's _glispToSlice can iterate over it.
// args is []any or nil.
func bridgeQuery(conn any, sql string, args any) ([]any, error) {
	var pgxArgs []any
	if args != nil {
		pgxArgs = args.([]any)
	}
	rows, err := conn.(*pgx.Conn).Query(context.Background(), sql, pgxArgs...)
	if err != nil {
		return nil, err
	}
	typed, err := pgx.CollectRows(rows, pgx.RowToMap)
	if err != nil {
		return nil, err
	}
	result := make([]any, len(typed))
	for i, r := range typed {
		result[i] = r
	}
	return result, nil
}

// bridgeExec runs an INSERT/UPDATE/DELETE and returns rows affected.
// args is []any or nil.
func bridgeExec(conn any, sql string, args any) (int64, error) {
	var pgxArgs []any
	if args != nil {
		pgxArgs = args.([]any)
	}
	tag, err := conn.(*pgx.Conn).Exec(context.Background(), sql, pgxArgs...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}
