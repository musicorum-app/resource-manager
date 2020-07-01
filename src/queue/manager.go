package queue

import (
	"fmt"
	"github.com/musicorum-app/resource-manager/database"
	"github.com/musicorum-app/resource-manager/utils"
	"time"
)

type ImageResource struct {
	URL      string `json:"url"`
	Hash     string `json:"hash"`
	FilePath string `json:"file_path"`
}

type ItemResult struct {
	ImageResource ImageResource
}

type Item struct {
	Resource string
	Data     string
	Callback func(result ItemResult)
}

var queueItems map[string]utils.Source
var queueTickers map[string]int8
var queueCallbacks map[string][]Item

func Initialize() {
	sources := utils.GetConfig().Sources
	queueItems = make(map[string]utils.Source)
	queueTickers = make(map[string]int8)
	queueCallbacks = make(map[string][]Item)

	for _, source := range sources {
		queueItems[source.Name] = source
		queueTickers[source.Name] = 0
	}

	for {
		go tick()
		time.Sleep(time.Second)
	}
}

func AddItem(source string, resource string, data string, done chan ItemResult) {
	println("NEW ITEM")
	go func() {
		callback := func(result ItemResult) {
			println("CHANNEL")
			done <- result
		}
		queueCallbacks[source] = append(queueCallbacks[source], Item{
			Resource: resource,
			Data:     data,
			Callback: callback,
		})
	}()
}

func tick() {
	for _, source := range queueItems {
		ticker := queueTickers[source.Name]
		print("Source " + source.Name + "; Ticker: ")
		print(ticker)
		println()
		println(len(queueCallbacks[source.Name]))
		for _, callbackItem := range queueCallbacks[source.Name] {
			fmt.Println(callbackItem.Data)
		}
		for _, callbackItem := range queueCallbacks[source.Name] {
			fmt.Println("CURRENT CALLBACK " + callbackItem.Data)
			queueCallback := queueCallbacks[source.Name]
			callbackItem := callbackItem
			go func() {
				runItem(callbackItem)
			}()
			queueCallbacks[source.Name] = queueCallback[1:]
			println("APPEND " + string(len(queueCallbacks[source.Name])))
			println(len(queueCallback))
		}
		queueTickers[source.Name] = ticker + 1
	}
}

func runItem(item Item) {
	println("RUNNING ITEM")
	//print(item.Resource)
	//println()

	switch item.Resource {
	case "ARTIST":
		runSpotifyArtist(item)
	}
}

func runSpotifyArtist(item Item) {
	result := database.FindResource("artists", item.Data)
	println(result.Spotify)
	item.Callback(ItemResult{
		ImageResource{
			URL:      "alo",
			Hash:     "asdasddasd",
			FilePath: "asdasdkkk",
		},
	})
	println("RUN FINISHED")
}
