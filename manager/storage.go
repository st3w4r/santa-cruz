package manager

import (
	"context"
	"database/sql"
	"time"

	"github.com/st3w4r/santa-cruz/dbsqlc"
	"github.com/st3w4r/santa-cruz/dbsqlc/dbmanager"
)

type StorageSystem interface {
	ListDbs() error
	CreateDb() error
	GetDb() error
}

type dbStorageSystem struct {
	db      *sql.DB
	queries *dbmanager.Queries
	dbPath  string
}

func InitDBManger() (dbStorageSystem, error) {
	ctx := context.Background()
	db, err := sql.Open("sqlite", ":memory:")

	if err != nil {
		return dbStorageSystem{}, err
	}

	if _, err := db.ExecContext(ctx, dbsqlc.CreateTableSchema); err != nil {
		return dbStorageSystem{}, err
	}

	queries := dbmanager.New(db)

	return dbStorageSystem{
		db:      db,
		queries: queries,
	}, nil
}

func (ds *dbStorageSystem) ListDbs() ([]dbmanager.Database, error) {
	ctx := context.Background()
	dbs, err := ds.queries.ListDatabases(ctx)
	if err != nil {
		return nil, err
	}
	return dbs, nil
}

func (ds *dbStorageSystem) AddDb(name, path string) (dbmanager.Database, error) {
	ctx := context.Background()
	timeNow := time.Now().UTC().Format(time.RFC3339)
	createdDb, err := ds.queries.CreateDatabase(ctx, dbmanager.CreateDatabaseParams{
		Name:      name,
		Path:      path,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	})
	if err != nil {
		return dbmanager.Database{}, err
	}
	return createdDb, nil
}
