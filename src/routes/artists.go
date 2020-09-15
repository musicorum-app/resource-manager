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

func RewindArtistsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data ArtistsRequest
	_ = json.NewDecoder(r.Body).Decode(&data)
	var promises []*promise.Promise

	for _, artist := range data.Artists {
		promises = append(promises, handleArtist(artist))
	}

	artists, err := promise.All(promises...).Await()

	if err != nil {
		utils.FailOnError(err)
	}

	var results []*structs.ArtistRewindResponse

	var toFetch []string

	for _, a := range artists.([]interface{}) {
		artist := a.(*structs.ArtistResponse)
		if artist != nil {
			var obj *structs.ArtistRewindResponse
			obj = new(structs.ArtistRewindResponse)

			if artist.Popularity != nil {
				var pop uint8
				pop = 0
				pop = *artist.Popularity

				obj.Name = artist.Name
				obj.Spotify = artist.Spotify
				obj.Image = artist.Url
				obj.Popularity = pop
			} else {
				toFetch = append(toFetch, artist.Spotify)
				obj.Name = artist.Name
				obj.Spotify = artist.Spotify
				obj.Image = artist.Url
				obj.Popularity = 0
			}

			results = append(results, obj)
		} else {
			results = append(results, nil)
		}
	}

	// TODO: Make this as promise
	fetched := api.GetArtists(toFetch)

	fmt.Println(fetched)

	//var finalResult []*structs.ArtistRewindResponse

	for _, result := range results {
		if result != nil {
			for _, fetch := range fetched {
				if fetch.ID == result.Spotify {
					result.Popularity = fetch.Popularity
					redis.SaveArtistRewindResponse(fetch)
				}
			}
		}
	}

	fmt.Println("To fetch:")
	fmt.Println(toFetch)

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
		search := database.FindArtist("artists", artist)
		if search != nil {
			result.Name = search.Key
			result.Hash = utils.Hash(artist)
			result.Url = search.Url
			result.Spotify = search.Spotify
			resolveFinal(result)
			redis.SetArtist(result)
		} else {
			p := promise.New(func(resolve func(interface{}), reject func(error)) {
				response := api.SearchArtist(artist)
				if response == nil {
					resolve(nil)
					resolveFinal(nil)
					return
				}
				var foundArtist api.ArtistObject
				foundArtist = *response
				redis.SaveArtistRewindResponse(foundArtist)

				pop := new(uint8)
				*pop = response.Popularity

				result.Name = response.Name
				result.Hash = utils.Hash(artist)
				result.Url = response.Image
				result.Spotify = response.ID
				result.Popularity = pop
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
