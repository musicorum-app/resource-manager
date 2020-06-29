package queue

import (
	"time"
)

var tickers map[string]int8

func Initialize() {
	tickers = make(map[string]int8)
	tickers["spotify"] = 0

	for {
		tick()
		time.Sleep(time.Second)
	}
}

func tick() {
	tickers["spotify"] = tickers["spotify"] + 1
	println("Ticking...")
	println(tickers["spotify"])
}
