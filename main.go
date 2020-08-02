package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shreddedbacon/fake-powerwall/api"
)

func main() {
	var inverter string
	flag.StringVar(&inverter, "inverter", "http://192.168.1.50", "URL or IP for the host")
	flag.Parse()
	a := api.FakePowerwall{
		Inverter: inverter,
	}
	r := mux.NewRouter()
	r.HandleFunc("/favicon.ico", faviconHandler)
	// API
	// Handle setting the maximum available amperate
	r.HandleFunc("/api/meters/aggregates", a.GetMetersAggregates).Methods("GET")
	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}

}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	favicon := `data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQEAYAAABPYyMiAAAABmJLR0T///////8JWPfcAAAACXBIWXMAAABIAAAASABGyWs+AAAAF0lEQVRIx2NgGAWjYBSMglEwCkbBSAcACBAAAeaR9cIAAAAASUVORK5CYII=\n`
	fmt.Fprintln(w, favicon)
}
