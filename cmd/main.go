package main

import (
	"math/rand"
	"time"
)

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func main() {
	bot := CreateOnlyPDFBot()
	bot.StartListenAndServ()
}
