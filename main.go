package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type ApiInfo struct {
	Uptime  string `json:"Uptime"`
	Info    string `json:"Info"`
	Version string `json:"Version"`
}

type TrackInfo struct {
	HeaderDate  string `json:"Header date"`
	Pilot       string `json:"Pilot"`
	Glider      string `json:"Glider"`
	GliderId    string `json:"Glider id"`
	TrackLength int    `json:"Track length"`
}

type TracksDB struct {
	tracks map[string]TrackInfo
}

var db TracksDB = TracksDB{}

var startTime time.Time

func init() {
	startTime = time.Now()
}

func uptime() string {
	now := time.Now()
	now.Format(time.RFC3339)
	startTime.Format(time.RFC3339)

	return now.Sub(startTime).String()
}

func (db *TracksDB) init() {
	db.tracks = make(map[string]TrackInfo)
}

func (db *TracksDB) Get(keyId string) (TrackInfo, bool) {
	s, ok := db.tracks[keyId]
	return s, ok
}

func replyWithAllTracksId(w *http.ResponseWriter, db *TracksDB) {

	if db.tracks == nil {
		json.NewEncoder(*w).Encode([]TrackInfo{})
	} else {
		json.NewEncoder(*w).Encode(db.tracks)
	}
}

func replyWithTracksId(w *http.ResponseWriter, db *TracksDB, id string) {
	//make sure i is valid
	s, ok := db.Get(id)
	if !ok {
		http.Error(*w, http.StatusText(404), http.StatusNotFound)
	}
	//handle TrackID
	json.NewEncoder(*w).Encode(s)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(404), 404)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	if parts[2] == "api" && parts[3] == "" {

		api := ApiInfo{uptime(), "Service for IGC tracks.", "v1"}

		json.NewEncoder(w).Encode(api)

	} else {
		http.Error(w, http.StatusText(404), 404)
	}
}

func igcHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		if r.Body == nil {
			http.Error(w, "IGC post requests must have an URL", http.StatusBadRequest)
			return
		}
	case "GET":
		http.Header.Add(w.Header(), "content-type", "application/json")
		parts := strings.Split(r.URL.Path, "/")

		if parts[2] == "api" && parts[3] == "igc" && len(parts) == 5 {
			if parts[4] == "" {
				replyWithAllTracksId(&w, &db)
			} else {
				http.Error(w, http.StatusText(404), 404)
			}
		} else if parts[2] == "api" && parts[3] == "igc" && parts[5] == "" && len(parts) == 6 {
			replyWithTracksId(&w, &db, parts[4])
		}
	default:
		http.Error(w, "Not implemented yet", http.StatusNotImplemented)
		return
	}
}

func main() {

	db.init()

	http.HandleFunc("/igcinfo/api/igc/", igcHandler)
	http.HandleFunc("/igcinfo/api/", apiHandler)
	http.HandleFunc("/igcinfo/", errorHandler)
	http.HandleFunc("/", errorHandler)

	http.ListenAndServe("127.0.0.1:8080", nil)

}

/*
	What: returns the meta information about a given track with the provided <id>, or NOT FOUND response code with an empty body.
	Response type: application/json
	Response code: 200 if everything is OK, appropriate error code otherwise.
*/
