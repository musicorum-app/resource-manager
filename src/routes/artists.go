package routes

import (
	"encoding/json"
	"fmt"
	"github.com/chebyrash/promise"
	"github.com/musicorum-app/resource-manager/api"
	"github.com/musicorum-app/resource-manager/database"
	"github.com/musicorum-app/resource-manager/queue"
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
	return promise.New(func(resolveFinal func(interface{}), reject func(error)) {
		fmt.Println("Searching for " + artist)
		redisSearch := redis.FindArtist(utils.Hash(artist))
		if redisSearch != nil {
			resolveFinal(redisSearch)
			return
		}
		var result *structs.ArtistResponse
		result = new(structs.ArtistResponse)
		search := database.FindResource("artists", artist)
		if search != nil {
			result.Name = search.Key
			result.Hash = utils.Hash(artist)
			result.Url = search.Url
			result.Spotify = search.Spotify
			fmt.Println(search.Key + "  /   " + string(search.CachedAt))
			resolveFinal(result)
			redis.SetArtist(result)
		} else {
			p := promise.New(func(resolve func(interface{}), reject func(error)) {
				fmt.Errorf("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
				response := api.SearchArtist(artist)
				if response == nil {
					resolve(nil)
					resolveFinal(nil)
					return
				}
				result.Name = response.Name
				result.Hash = utils.Hash(artist)
				result.Url = response.Image
				result.Spotify = response.ID
				go func() {
					database.InsertArtist(response.Name, response.ID, response.Image)
				}()
				fmt.Println("RESULT THIG HERE")
				resolve(response)
				resolveFinal(result)
				redis.SetArtist(result)
			})
			action := queue.AddItem("spotify", p)
			response, err := action.Await()
			if err != nil {
				resolveFinal(nil)
				return
			}
			if response == nil {
				resolveFinal(nil)
				return
			}
		}
	})
}
