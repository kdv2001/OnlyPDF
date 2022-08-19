package main

import (
	"OnlyPDF/app/handlers"
	"OnlyPDF/app/usecase/impl"
	"fmt"
	"gopkg.in/telebot.v3"
	"os"
	"time"
)

type OnlyPDFBot struct {
	handler *handlers.Handlers
}

func CreateOnlyPDFBot() *OnlyPDFBot {
	fileUseCase := impl.CreateFileUseCase()
	handler := handlers.CreateHandlers(&fileUseCase)
	return &OnlyPDFBot{&handler}
}

func (b *OnlyPDFBot) StartListenAndServ() {
	bot, err := telebot.NewBot(telebot.Settings{Token: os.Getenv("OnlyPDFBotToken"), Poller: &telebot.LongPoller{Timeout: 10 * time.Second}})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//bot.Use(middleware.Logger())
	bot.Handle(telebot.OnDocument, func(ctx telebot.Context) error {
		fmt.Println(ctx.Message().Document.FileID)
		fmt.Println(ctx.Message().ID)
		return nil
	})

	bot.Handle("/m", b.handler.AddFiles)
	bot.Handle("/t", func(m telebot.Context) error {
		//document := &telebot.Document{File: telebot.FromDisk("api2.png"), FileName: "Фото.png"}
		//msg := telebot.Message{ID: m.Message().ID - 1}
		file, err := bot.FileByID(bot.)
		if err != nil {
			fmt.Println(err)
		}
		bot.Download(&file, m.Message().Sender.Username)
		fmt.Println(file.FileID)
		return nil
	})

	bot.Start()
}
