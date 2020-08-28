package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/musicorum-app/resource-manager/api"
	"github.com/musicorum-app/resource-manager/database"
	"github.com/musicorum-app/resource-manager/queue"
	"github.com/musicorum-app/resource-manager/redis"
	"github.com/musicorum-app/resource-manager/routes"
	"github.com/musicorum-app/resource-manager/utils"
	"github.com/rs/cors"
	"log"
	"net/http"
	"sync"
)

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	database.Initialize()
	api.Initialize()
	redis.InitializeRedis()

	go func() {
		server()
	}()

	go func() {
		queue.Initialize()
	}()

	wg.Wait()
}

func server() {
	log.Println("Starting web server...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/fetch/artists", routes.ArtistsHandler)
	router.HandleFunc("/fetch/albums", routes.AlbumsHandler)
	router.HandleFunc("/fetch/tracks", routes.TracksHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Fatal(http.ListenAndServe(":"+utils.GetEnvVar("PORT"), handler))
}

func index(w http.ResponseWriter, _ *http.Request) {
	mapIndex := map[string]string{"working": "ok"}
	marshal, _ := json.Marshal(mapIndex)
	fmt.Fprintln(w, string(marshal))
}
