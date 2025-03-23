package dbmanaged

import (
	"context"
	"database/sql"
)

type PragmaTableInfo struct {
	Cid       sql.NullInt64
	Name      sql.NullString
	Type      sql.NullString
	Notnull   sql.NullInt64
	DfltValue sql.NullString
	Pk        sql.NullInt64
}

const listColumns = `-- name: ListColumns :many
SELECT
    "cid",
    "name",
    "type",
    "notnull",
    "dflt_value",
    "pk"
FROM pragma_table_info(?)
`

type ListColumnsParams struct {
	TableName string
}

func (q *Queries) ListColumns(ctx context.Context, arg ListColumnsParams) ([]PragmaTableInfo, error) {
	rows, err := q.db.QueryContext(ctx, listColumns, arg.TableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PragmaTableInfo
	for rows.Next() {
		var i PragmaTableInfo
		if err := rows.Scan(
			&i.Cid,
			&i.Name,
			&i.Type,
			&i.Notnull,
			&i.DfltValue,
			&i.Pk,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
