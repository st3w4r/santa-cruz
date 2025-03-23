package server

import (
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

type ResponseTable struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	TableName string `json:"table_name"`
	Sql       string `json:"sql"`
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
		respTables = append(respTables, ResponseTable{
			Name:      table.Name,
			Type:      table.Type,
			TableName: table.TblName,
			Sql:       table.Sql,
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

	err := http.ListenAndServe(host+":"+port, r)
	if err != nil {
		return err
	}
	return nil
}
