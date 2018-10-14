package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"../../../src/github.com/marni/goigc"
)

func checkId(id string) bool {
	idExists := false
	for i := 0; i < len(IDs); i++ {
		if IDs[i] == strings.ToUpper(id) {
			idExists = true
			break
		}
	}
	if idExists == true {
		return true
	} else {
		return false
	}
}

func checkURL(u string) bool {
	check, _ := regexp.MatchString("^(http://skypolaris.org/wp-content/uploads/IGS%20Files/)(.*?)(%20)(.*?)(.igc)$", u)
	if check == true {
		return true
	}
	return false
}

func replyWithAllTracksId(w http.ResponseWriter, db TrackDB) {
	if len(IDs) == 0 {
		IDs = make([]string, 0)
	}
	json.NewEncoder(w).Encode(IDs)
	return
}

func replyWithTracksId(w http.ResponseWriter, db TrackDB, id string) {
	t, _ := db.Get(strings.ToUpper(id))
	http.Header.Set(w.Header(), "content.type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func replyWithField(w http.ResponseWriter, db TrackDB, id string, field string) {
	t, _ := db.Get(strings.ToUpper(id))

	switch strings.ToUpper(field) {
	case "PILOT":
		fmt.Fprint(w, t.Pilot)
	case "GLIDER":
		fmt.Fprint(w, t.Glider)
	case "GLIDER_ID":
		fmt.Fprint(w, t.GliderId)
	case "TRACK_LENGTH":
		fmt.Fprint(w, t.TrackLength)
	case "H_DATE":
		fmt.Fprint(w, t.HeaderDate)
	default:
		http.Error(w, "Not a valid option", http.StatusNotFound)
		return
	}
}

func handlerIgc(w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case "POST":
		if r.Body == nil {
			http.Error(w, "Missing body", http.StatusBadRequest)
			return
		}
		var u URL
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if checkURL(u.URL) == false {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		track, err := igc.ParseLocation(u.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		totalDistance := CalculatedDistance(track)
		var i ID
		i.ID = ("ID" + strconv.Itoa(lastUsed))
		t := Track{track.Header.Date,
			track.Pilot,
			track.GliderType,
			track.GliderID,
			totalDistance}
		lastUsed++
		if db.tracks == nil {
			db.Init()
		}
		db.Add(t, i)
		return
	case "GET":
		/*if len(IDs) == 0 {
			IDs = make([]string, 0)
		}
		json.NewEncoder(w).Encode(IDs)
		return*/

		parts := strings.Split(r.URL.Path, "/")




		if len(parts) == 5 {
			if parts[4] == ""{
				replyWithAllTracksId(w, db)
			} else {
				http.Error(w, http.StatusText(404), 404)
			}
		} else if (parts[5] == "" && len(parts) == 6) || len(parts) == 5 && checkId(parts[4]) {
			/*idExists := false
			for i := 0; i < len(IDs); i++ {
				if IDs[i] == strings.ToUpper(parts[4]) {
					idExists = true
					break
				}
			}*/
			checkId(parts[4])
			if !checkId(parts[4])/*!idExists*/ {
				http.Error(w, "ID out of range.", http.StatusNotFound)
				return
			} else {
				replyWithTracksId(w, db, parts[4])
			}
		} else if (parts[6] == "" && len(parts) == 7) || len(parts) == 6 && checkId(parts[4]) {
			replyWithField(w, db, parts[4], parts[5])
		} else {
			http.Error(w, "Not a valid request", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handlerApi(w http.ResponseWriter, r *http.Request){
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	if parts[2] == "api" && parts[3] == "" {
		api := Information{uptime(), Description, Version}
		json.NewEncoder(w).Encode(api)
	} else {
		http.Error(w, http.StatusText(404), 404)
	}
}

/*func handlerIdAndField(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	idExists := false
	for i := 0; i < len(IDs); i++ {
		if IDs[i] == strings.ToUpper(parts[4]) {
			idExists = true
			break
		}
	}
	if !idExists {
		http.Error(w, "ID out of range.", http.StatusNotFound)
		return
	}
	t, _ := db.Get(strings.ToUpper(parts[4]))
	if len(parts) == 5 {
		http.Header.Set(w.Header(), "content.type", "application/json")
		json.NewEncoder(w).Encode(t)
	}
	if len(parts) == 6 {
		switch strings.ToUpper(parts[5]) {
		case "PILOT":
			fmt.Fprint(w, t.Pilot)
		case "GLIDER":
			fmt.Fprint(w, t.Glider)
		case "GLIDER_ID":
			fmt.Fprint(w, t.GliderId)
		case "TRACK_LENGTH":
			fmt.Fprint(w, t.TrackLength)
		case "H_DATE":
			fmt.Fprint(w, t.HeaderDate)
		default:
			http.Error(w, "Not a valid option", http.StatusNotFound)
			return
		}

	}
	if len(parts) > 6 {
		http.Error(w, "Too many /.", http.StatusNotFound)
	}
}*/