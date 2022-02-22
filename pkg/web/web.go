package web

import (
	"fmt"
	"net/http"

	"github.com/4molybdenum2/distrKV/db"
)

type Server struct {
	db *db.Database
}

func NewServer(db *db.Database) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) GetKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get Endpoint hit!\n"))
	r.ParseForm()
	key := r.Form.Get("key")
	value, err := s.db.GetKey(key)
	fmt.Fprintf(w, "Value = %q, error = %v", value, err)
}

func (s *Server) SetKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Set Endpoint hit!\n"))
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "Error = %v", err)
}
