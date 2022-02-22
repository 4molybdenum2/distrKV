package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	db "github.com/4molybdenum2/distrKV/db"
	"github.com/4molybdenum2/distrKV/pkg/web"
	"github.com/gorilla/mux"
)

// defined database location
var (
	path     = flag.String("loc", "", "Path to Bolt DB database")
	httpAddr = flag.String("addr", "127.0.0.1:8080", "HTTP host endpoint")
)

func parseFlags() {
	flag.Parse()
	if *path == "" {
		log.Fatal("Must provide database location, ...in this case a file")
	}
}

var d *db.Database
var closeFunc func() error
var err error

func main() {
	fmt.Println("distrKV is a Distributed Key-Value Store")
	parseFlags()

	d, closeFunc, err = db.NewDatabase(*path)
	if err != nil {
		log.Fatalf("New Database (%q) : %v", *path, err)
	}
	defer closeFunc()

	// create new server
	srv := web.NewServer(d)

	// defined http router
	r := mux.NewRouter()
	r.HandleFunc("/get", srv.GetKeyHandler)
	r.HandleFunc("/set", srv.SetKeyHandler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", r))
}
