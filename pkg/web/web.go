package web

import (
	"fmt"
	"hash/fnv"
	"io"
	"net/http"

	"github.com/4molybdenum2/distrKV/db"
)

type Server struct {
	db         *db.Database
	shardCount int
	shardIdx   int
	addr       map[int]string
}

func NewServer(db *db.Database, shardCount int, shardIdx int, addr map[int]string) *Server {
	return &Server{
		db:         db,
		shardCount: shardCount,
		shardIdx:   shardIdx,
		addr:       addr,
	}
}

func (s *Server) redirect(shardId int, w http.ResponseWriter, r *http.Request) {
	url := "http://" + s.addr[shardId] + r.RequestURI
	fmt.Fprintf(w, "redirecting from shard %d to shard %d (%q)\n", s.shardIdx, shardId, url)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error redirecting the request: %v", err)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

func (s *Server) getShard(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	shardIdx := int(h.Sum64() % uint64(s.shardCount))
	return shardIdx
}

func (s *Server) GetKeyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	// get shard
	shardId := s.getShard(key)
	value, err := s.db.GetKey(key)

	// if not current shard redirect request to necessary shard
	if shardId != s.shardIdx {
		s.redirect(shardId, w, r)
		return
	}
	fmt.Fprintf(w, "Shard = %d, current shard = %d, addr = %q, Value = %q, error = %v", shardId, s.shardIdx, s.addr[shardId], value, err)
}

func (s *Server) SetKeyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	shardId := s.getShard(key)
	// if calculated shard in which we want to set key is not the current shard then we need to
	// make a set request on another address of the calculated shard
	if shardId != s.shardIdx {
		// redirect to necessary shard
		s.redirect(shardId, w, r)
		return
	}

	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "Error = %v, shard= %d", err, shardId)
}
