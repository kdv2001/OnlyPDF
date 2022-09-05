package main

import (
	"fmt"
	"os"
)

func main() {
	bot, err := CreateOnlyPDFBot()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = bot.StartListenAndServ()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
