package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/musicorum-app/resource-manager/api"
	"github.com/musicorum-app/resource-manager/database"
	"github.com/musicorum-app/resource-manager/redis"
	"github.com/musicorum-app/resource-manager/routes"
	"github.com/musicorum-app/resource-manager/utils"
	"go/types"
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
		// queue.Initialize()
	}()

	wg.Wait()
}

func server() <-chan types.Nil {
	log.Println("Starting web server...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/status", status)
	router.HandleFunc("/fetch/artists", routes.ArtistsHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":"+utils.GetEnvVar("PORT"), router))

	return nil
}

func index(w http.ResponseWriter, _ *http.Request) {
	mapIndex := map[string]string{"working": "ok"}
	marshal, _ := json.Marshal(mapIndex)
	fmt.Fprintln(w, string(marshal))
}

func status(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Todo Index!")
}
