package queue

import (
	"errors"
	"fmt"
	"github.com/chebyrash/promise"
	"github.com/musicorum-app/resource-manager/utils"
	"math/rand"
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

type Action struct {
	Promise *promise.Promise
	Resolve func(interface{})
	Reject  func(error)
	Key     string
}

var queueItems map[string]utils.Source
var queueActions map[string][]Action
var completeActions map[string]bool

func Initialize() {
	sources := utils.GetConfig().Sources
	queueItems = make(map[string]utils.Source)
	queueActions = make(map[string][]Action)
	completeActions = make(map[string]bool)

	for _, source := range sources {
		queueItems[source.Name] = source
	}

	for {
		go tick()
		time.Sleep(time.Second)
	}
}

func AddItem(source string, action *promise.Promise) *promise.Promise {
	return promise.New(func(resolve func(interface{}), reject func(error)) {
		fmt.Println("ADDING PROMISE ITEM")
		item, found := queueItems[source]
		if !found {
			fmt.Println("Source " + source + " not found.")
			reject(errors.New("Source " + source + " not found."))
		}
		key := utils.Hash(string(rand.Intn(120)))
		queueActions[item.Name] = append(queueActions[item.Name], Action{
			Promise: action,
			Resolve: resolve,
			Reject:  reject,
			Key:     key,
		})
		completeActions[key] = false
		if len(queueActions[item.Name]) < int(item.RateLimit) {
			tick()
		}
	})
}

func tick() {
	fmt.Println("tick")
	for _, source := range queueItems {
		var promises []*promise.Promise
		for n := 0; n < int(source.RateLimit); n++ {
			if len(queueActions[source.Name]) > n {
				action := queueActions[source.Name][n]

				p := promise.New(func(resolve func(interface{}), reject func(error)) {
					result, err := action.Promise.Await()
					fmt.Println("RESOLVING ITEM")
					if err != nil {
						reject(err)
						action.Reject(err)
					} else {
						resolve(result)
						action.Resolve(result)
					}
				}).Then(func(data interface{}) interface{} {
					completeActions[action.Key] = true
					return nil
				}).Catch(func(err error) error {
					completeActions[action.Key] = true
					fmt.Println("AN ERROR HERE TAKE")
					return nil
				})
				promises = append(promises, p)
			}
		}

		if len(promises) > 0 {
			_, err := promise.All(promises...).Await()
			if err != nil {
				fmt.Println("Error on the promise list.")
				fmt.Println(err)
			}
		}
		cleanActions(source)
	}
}

func cleanActions(source utils.Source) {
	actions := queueActions[source.Name]
	var newActions []Action
	for _, ac := range actions {
		if !completeActions[ac.Key] {
			newActions = append(newActions, ac)
		}
	}
	queueActions[source.Name] = newActions
}
