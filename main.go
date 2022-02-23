package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	db "github.com/4molybdenum2/distrKV/db"
	"github.com/4molybdenum2/distrKV/pkg/config"
	"github.com/4molybdenum2/distrKV/pkg/web"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
)

// defined database location
var (
	path     = flag.String("loc", "", "Path to Bolt DB database")
	httpAddr = flag.String("addr", "127.0.0.1:8080", "HTTP host endpoint")
	cnf      = flag.String("config", "sharding.toml", "Config file for sharding data")
	shard    = flag.String("shard", "", "Name of the shard for data")
)

func parseFlags() {
	flag.Parse()
	if *path == "" {
		log.Fatal("Must provide database location, ...in this case a file")
	}
	if *shard == "" {
		log.Fatal("Must provide name of the shard")
	}
}

func main() {
	fmt.Println("distrKV is a Distributed Key-Value Store")
	parseFlags()

	// get toml data from file
	var config config.Config
	if _, err := toml.DecodeFile(*cnf, &config); err != nil {
		log.Fatalf("DecodeFile(%q): %v", *cnf, err)
	}

	shardCount := len(config.Shards)
	var shardIdx int = -1
	var shardAddr = make(map[int]string)

	for _, s := range config.Shards {
		shardAddr[s.Idx] = s.Address
		if s.Name == *shard {
			shardIdx = s.Idx
		}
	}

	if shardIdx < 0 {
		log.Fatalf("Shard not found with name %q", *shard)
	}
	log.Printf("Shard count is %d and current shard is %d", shardCount, shardIdx)

	d, closeFunc, err := db.NewDatabase(*path)
	if err != nil {
		log.Fatalf("New Database (%q) : %v", *path, err)
	}
	defer closeFunc()

	// create new server
	srv := web.NewServer(d, shardCount, shardIdx, shardAddr)

	// defined http router
	r := mux.NewRouter()
	r.HandleFunc("/get", srv.GetKeyHandler)
	r.HandleFunc("/set", srv.SetKeyHandler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(*httpAddr, r))
}
