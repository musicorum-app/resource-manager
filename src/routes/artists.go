package routes

import (
	"encoding/json"
	"fmt"
	"github.com/chebyrash/promise"
	"github.com/musicorum-app/resource-manager/api"
	"github.com/musicorum-app/resource-manager/database"
	"github.com/musicorum-app/resource-manager/redis"
	"github.com/musicorum-app/resource-manager/structs"
	"github.com/musicorum-app/resource-manager/utils"
	"io"
	"net/http"
)

type ArtistsRequest struct {
	Artists []string `json:"artists"`
}

func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data ArtistsRequest
	_ = json.NewDecoder(r.Body).Decode(&data)
	var promises []*promise.Promise

	for _, artist := range data.Artists {
		promises = append(promises, handleArtist(artist))
	}

	results, err := promise.All(promises...).Await()

	jsonResult, err := json.Marshal(results)
	utils.FailOnError(err)

	io.WriteString(w, string(jsonResult))
}

func handleArtist(artist string) *promise.Promise {
	return promise.New(func(resolve func(interface{}), reject func(error)) {
		fmt.Println("Searching for " + artist)
		var result *structs.ArtistResponse
		result = new(structs.ArtistResponse)
		search := database.FindResource("artists", artist)
		if search != nil {
			result.Name = search.Key
			result.Hash = utils.Hash(artist)
			result.Url = search.Url
			result.Spotify = search.Spotify
		} else {
			response := api.SearchArtist(artist)
			if response == nil {
				resolve(nil)
			}
			result.Name = response.Name
			result.Hash = utils.Hash(artist)
			result.Url = response.Image
			result.Spotify = response.ID
			go func() {
				database.InsertArtist(response.Name, response.ID, response.Image)
			}()
		}
		redis.SetArtist(result)

		//artists = append(artists, result.ImageResource)
		resolve(result)
	})
}
