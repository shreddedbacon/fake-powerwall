package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/shreddedbacon/fake-powerwall/api"
)

func main() {
	var inverter string
	var inverterType string
	var inverterCloud bool
	var authToken string
	flag.StringVar(&inverter, "inverter", "http://192.168.1.50", "URL or IP for the host")
	flag.StringVar(&inverterType, "inverter-type", "fronius", "What type of inverter (fronius,solaredge,solaredge)")
	flag.BoolVar(&inverterCloud, "inverter-cloud", false, "Cloud based API?")
	flag.StringVar(&authToken, "auth-token", "", "If inverter-type requires an auth token, define it here")
	flag.Parse()

	// get the inverter host from an envvar
	inverter = getEnv("INVERTER_HOST", inverter)
	inverterType = getEnv("INVERTER_TYPE", inverterType)
	authToken = getEnv("AUTH_TOKEN", authToken)

	a := api.FakePowerwall{
		Inverter:     inverter,
		InverterType: inverterType,
		CloudBased:   inverterCloud,
		AuthToken:    authToken,
	}

	log.Println("Starting Fake Powerwall")

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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
