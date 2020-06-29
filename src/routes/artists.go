package routes

import (
	"encoding/json"
	"net/http"
)

type ArtistsRequest struct {
	Artists []string `json:"artists"`
}

func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data ArtistsRequest
	_ = json.NewDecoder(r.Body).Decode(&data)
	for _, artist := range data.Artists {
		println(artist)
	}
	json.NewEncoder(w).Encode(&data)
}
