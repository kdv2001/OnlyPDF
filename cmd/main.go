package main

import "log"

func main() {
	bot, err := CreateOnlyPDFBot()
	if err != nil {
		log.Fatal(err)
	}

	err = bot.StartListenAndServ()
	if err != nil {
		log.Fatal(err)
	}
}
