package manager

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"

	"github.com/st3w4r/santa-cruz/config"
	"github.com/st3w4r/santa-cruz/dbsqlc"
	"github.com/st3w4r/santa-cruz/dbsqlc/dbmanager"
)

type dbStorageSystem struct {
	db      *sql.DB
	queries *dbmanager.Queries
	dbPath  string
}

func InitDBManger() (dbStorageSystem, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return dbStorageSystem{}, err
	}

	dbPath := cfg.DB_MANAGER_PATH
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(dbPath), os.ModePerm)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return dbStorageSystem{}, err
	}

	ctx := context.Background()
	if _, err := db.ExecContext(ctx, dbsqlc.CreateTableSchema); err != nil {
		return dbStorageSystem{}, err
	}

	queries := dbmanager.New(db)

	return dbStorageSystem{
		db:      db,
		queries: queries,
		dbPath:  dbPath,
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

func (ds *dbStorageSystem) AddDb(name, path, desc string) (dbmanager.Database, error) {
	ctx := context.Background()
	timeNow := time.Now().UTC().Format(time.RFC3339)
	createdDb, err := ds.queries.CreateDatabase(ctx, dbmanager.CreateDatabaseParams{
		Name:        name,
		Path:        path,
		Description: sql.NullString{String: desc, Valid: true},
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return dbmanager.Database{}, fmt.Errorf("database with name '%s' already exists", name)
		}
		return dbmanager.Database{}, err
	}
	return createdDb, nil
}

func (ds *dbStorageSystem) GetDb(id int64) (dbmanager.Database, error) {
	ctx := context.Background()
	db, err := ds.queries.GetDatabase(ctx, id)
	if err != nil {
		return dbmanager.Database{}, err
	}
	return db, nil
}

func (ds *dbStorageSystem) RemoveDb(id int64) error {
	ctx := context.Background()
	err := ds.queries.DeleteDatabase(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (ds *dbStorageSystem) UpdateDb(id int64, name, path, desc string) (dbmanager.Database, error) {
	ctx := context.Background()
	timeNow := time.Now().UTC().Format(time.RFC3339)

	db, err := ds.queries.UpdateDatabase(ctx, dbmanager.UpdateDatabaseParams{
		ID:          id,
		Name:        sql.NullString{String: name, Valid: name != ""},
		Path:        sql.NullString{String: path, Valid: path != ""},
		Description: sql.NullString{String: desc, Valid: desc != ""},
		UpdatedAt:   timeNow,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return dbmanager.Database{}, fmt.Errorf("database with name '%s' already exists", name)
		}
		return dbmanager.Database{}, err
	}
	return db, nil
}
