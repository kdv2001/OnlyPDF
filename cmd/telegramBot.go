package main

import (
	"OnlyPDF/app"
	"OnlyPDF/app/handlers"
	"OnlyPDF/app/repositories/gotemberg"
	"OnlyPDF/app/repositories/memory"
	"OnlyPDF/app/repositories/telegramFiles"
	"OnlyPDF/app/usecase/impl"
	"fmt"
	"net/http"
	"os"
	"time"

	"gopkg.in/telebot.v3"
)

type OnlyPDFBot struct {
	handler *handlers.Handlers
}

const pollTime = 10 * time.Second

func CreateOnlyPDFBot() (OnlyPDFBot, error) {
	repo, err := memory.CreateFilesPostgresInMemory()
	if err != nil {
		return OnlyPDFBot{}, err
	}

	bot, err := telebot.NewBot(telebot.Settings{Token: os.Getenv("OnlyPDFBotToken")})
	if err != nil {
		return OnlyPDFBot{}, err
	}

	g := gotemberg.NewRepo(http.DefaultClient, os.Getenv("GotenbergURL"))
	fileLoader := telegramFiles.NewTelegramFiles(bot)
	fileUseCase := impl.CreateFileUseCase(repo, &fileLoader, g)
	handler := handlers.CreateHandlers(fileUseCase)

	return OnlyPDFBot{handler}, nil
}

func (b *OnlyPDFBot) StartListenAndServ() error {
	bot, err := telebot.NewBot(telebot.Settings{Token: os.Getenv("OnlyPDFBotToken"),
		Poller: &telebot.LongPoller{Timeout: pollTime}})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	bot.Handle(telebot.OnMedia, b.handler.AddPhoto)
	bot.Handle(telebot.OnDocument, b.handler.AddFile)
	bot.Handle("/start", func(ctx telebot.Context) error {
		return ctx.Send(app.DefaultMsg, app.ReturnMainMenu())
	})

	bot.Handle(&app.BtnHelp, func(ctx telebot.Context) error {
		return ctx.Send(app.HelpMsg, app.ReturnMainMenu())
	})
	bot.Handle("/help", func(ctx telebot.Context) error {
		return ctx.Send(app.HelpMsg, app.ReturnMainMenu())
	})

	bot.Handle("/err", func(ctx telebot.Context) error {
		return telebot.NewError(http.StatusInternalServerError, "fdfdfdfd")
	})

	bot.Handle(&app.BtnPrint, b.handler.ShowFiles)
	bot.Handle("/print", b.handler.ShowFiles)

	bot.Handle(&app.BtnClear, b.handler.ClearFiles)
	bot.Handle("/clear", b.handler.ClearFiles, handlers.StateMiddleware)

	bot.Handle(&app.BtnMerge, b.handler.Merge)
	bot.Handle("/merge", b.handler.MergeCommand)

	bot.Handle(&app.BtnConvert, b.handler.ConvertCommand)
	bot.Handle("/convert", b.handler.ConvertCommand)

	bot.Start()
	return nil
}
