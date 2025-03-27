package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/st3w4r/santa-cruz/manager"
)

type server struct {
	controlDb manager.DbStorageSystem
	dataDbs   []manager.ManagedDb
}

func NewServer(controlDb manager.DbStorageSystem, dataDbs []manager.ManagedDb) *server {
	return &server{
		controlDb: controlDb,
		dataDbs:   dataDbs,
	}
}

func (s *server) getRootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Database mangement system is running"))
}

type ResponseColumn struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Notnull bool   `json:"not_null"`
}

type ResponseTable struct {
	Name      string           `json:"name"`
	Type      string           `json:"type"`
	TableName string           `json:"table_name"`
	Sql       string           `json:"sql"`
	Columns   []ResponseColumn `json:"columns,omitempty"`
}

type ResponseDatabase struct {
	Id          int             `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Tables      []ResponseTable `json:"tables,omitempty"`
}

type ResponseDatabases struct {
	Databases []ResponseDatabase `json:"databases"`
}

func (s *server) listDatabasesHandler(w http.ResponseWriter, r *http.Request) {
	dbs, err := s.controlDb.ListDbs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respDbs := []ResponseDatabase{}
	for _, db := range dbs {
		respDbs = append(respDbs, ResponseDatabase{
			Id:          int(db.ID),
			Name:        db.Name,
			Description: db.Description.String,
		})
	}

	respJson, err := json.Marshal(ResponseDatabases{
		Databases: respDbs,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func (s *server) getDatabaseHandler(w http.ResponseWriter, r *http.Request) {

	dbIdParam := chi.URLParam(r, "dbId")

	dbId, err := strconv.Atoi(dbIdParam)
	if err != nil {
		err = fmt.Errorf("database id must be an integer")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := s.controlDb.GetDb(int64(dbId), "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dataDb, err := manager.ConnectToManagedDb(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dataDb.Db.Close()

	tables, err := manager.ListTablesDb(dataDb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respTables := []ResponseTable{}
	for _, table := range tables {

		columns, err := manager.ListColumnsTable(dataDb, table.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respColumns := []ResponseColumn{}
		for _, column := range columns {
			respColumns = append(respColumns, ResponseColumn{
				Name:    column.Name,
				Type:    column.Type,
				Notnull: column.Notnull,
			})
		}

		respTables = append(respTables, ResponseTable{
			Name:      table.Name,
			Type:      table.Type,
			TableName: table.TblName,
			Sql:       table.Sql,
			Columns:   respColumns,
		})
	}

	respJson, err := json.Marshal(ResponseDatabase{
		Id:          int(db.ID),
		Name:        db.Name,
		Description: db.Description.String,
		Tables:      respTables,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

type respQuery struct {
	Columns []string         `json:"columns"`
	Data    []map[string]any `json:"data"`
}

func parseColValue(col *sql.ColumnType, val []byte) any {
	colType := col.DatabaseTypeName()
	if colType == "INTEGER" || colType == "INT" {
		val, err := strconv.Atoi(string(val))
		if err != nil {
			return nil
		}
		return val
	} else if colType == "BOOLEAN" || colType == "BOOL" {
		val, err := strconv.ParseBool(string(val))
		if err != nil {
			return nil
		}
		return val
	} else if colType == "REAL" || colType == "FLOAT" {
		val, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return nil
		}
		return val
	} else if colType == "TEXT" {
		return string(val)
	} else if colType == "BLOB" {
		return string(val)
	} else if colType == "JSON" {
		if string(val) == "" {
			return nil
		}
		var jsonVal any
		err := json.Unmarshal(val, &jsonVal)
		if err != nil {
			return nil
		}
		return jsonVal
	} else if colType == "NULL" {
		return nil
	}
	return string(val)
}

func (s *server) getRunQueryHandler(w http.ResponseWriter, r *http.Request) {
	dbIdParam := chi.URLParam(r, "dbId")

	dbId, err := strconv.Atoi(dbIdParam)
	if err != nil {
		err = fmt.Errorf("database id must be an integer")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := s.controlDb.GetDb(int64(dbId), "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		err = fmt.Errorf("?query parameter is required, example: ?query=SELECT * FROM table")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataDb, err := manager.ConnectToManagedDb(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dataDb.Db.Close()

	ctx := r.Context()
	rows, err := dataDb.Db.QueryContext(ctx, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cols, err := rows.ColumnTypes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	row := make([][]byte, len(cols))
	rowPtr := make([]any, len(cols))
	for i := range row {
		rowPtr[i] = &row[i]
	}
	data := []map[string]any{}

	for rows.Next() {
		_ = rows.Scan(rowPtr...)
		rowData := map[string]any{}
		for i, r := range row {
			colName := cols[i].Name()
			rowData[colName] = parseColValue(cols[i], r)
		}
		data = append(data, rowData)
	}

	colNames := []string{}
	for _, cols := range cols {
		colNames = append(colNames, cols.Name())
	}

	resp := respQuery{
		Columns: colNames,
		Data:    data,
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func (s *server) Serve(host, port string) error {
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)

	r.Get("/", s.getRootHandler)
	r.Get("/databases", s.listDatabasesHandler)
	r.Get("/databases/{dbId}", s.getDatabaseHandler)
	r.Get("/databases/{dbId}/run", s.getRunQueryHandler)

	err := http.ListenAndServe(host+":"+port, r)
	if err != nil {
		return err
	}
	return nil
}
