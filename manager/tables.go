package manager

import (
	"context"
	"database/sql"

	"github.com/st3w4r/santa-cruz/dbsqlc/dbmanager"
	"github.com/st3w4r/santa-cruz/dbsqlcmanaged/dbmanaged"
)

type Table struct {
	Name    string
	Type    string
	TblName string
	Sql     string
}

type ManagedDb struct {
	Db      *sql.DB
	Queries *dbmanaged.Queries
	DbPath  string
}

func ConnectToManagedDb(mdb dbmanager.Database) (ManagedDb, error) {
	db, err := sql.Open("sqlite", mdb.Path)
	if err != nil {
		return ManagedDb{}, err
	}

	queries := dbmanaged.New(db)

	return ManagedDb{
		Db:      db,
		Queries: queries,
		DbPath:  mdb.Path,
	}, nil
}

func ListTablesDb(mdb ManagedDb) ([]Table, error) {

	rows, err := mdb.Queries.ListTables(context.Background())
	if err != nil {
		return nil, err
	}

	var tables []Table
	for _, row := range rows {
		tables = append(tables, Table{
			Name:    row.Name.String,
			Type:    row.Type.String,
			TblName: row.TblName.String,
			Sql:     row.Sql.String,
		})
	}

	return tables, nil
}
