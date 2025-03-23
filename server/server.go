package server

import (
	"encoding/json"
	"net/http"

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

type ResponseListDatabases struct {
	Databases []string `json:"databases"`
}

func (s *server) listDatabases(w http.ResponseWriter, r *http.Request) {
	dbs, err := s.controlDb.ListDbs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dbNames := []string{}
	for _, db := range dbs {
		dbNames = append(dbNames, db.Name)
	}

	respJson, err := json.Marshal(ResponseListDatabases{
		Databases: dbNames,
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
	r.Get("/databases", s.listDatabases)

	err := http.ListenAndServe(host+":"+port, r)
	if err != nil {
		return err
	}
	return nil
}
