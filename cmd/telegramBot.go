package main

import (
	"OnlyPDF/app"
	"OnlyPDF/app/handlers"
	"OnlyPDF/app/repositories/memory"
	"OnlyPDF/app/usecase/impl"
	"fmt"
	"gopkg.in/telebot.v3"
	"os"
	"time"
)

type OnlyPDFBot struct {
	handler *handlers.Handlers
}

func CreateOnlyPDFBot() (OnlyPDFBot, error) {
	repo, err := memory.CreateFilesPostgresInMemory()
	if err != nil {
		return OnlyPDFBot{}, err
	}
	fileUseCase := impl.CreateFileUseCase(repo)
	handler := handlers.CreateHandlers(fileUseCase)
	return OnlyPDFBot{handler}, nil
}

func (b *OnlyPDFBot) StartListenAndServ() error {
	bot, err := telebot.NewBot(telebot.Settings{Token: os.Getenv("OnlyPDFBotToken"), Poller: &telebot.LongPoller{Timeout: 10 * time.Second}})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	bot.Handle(telebot.OnDocument, b.handler.AddFiles)
	bot.Handle("/start", func(ctx telebot.Context) error {
		return ctx.Send(app.DefaultMsg, app.ReturnMainMenu())
	})

	bot.Handle(&app.BtnHelp, func(ctx telebot.Context) error {
		return ctx.Send(app.HelpMsg, app.ReturnMainMenu())
	})
	bot.Handle("/help", func(ctx telebot.Context) error {
		return ctx.Send(app.HelpMsg, app.ReturnMainMenu()) //ctx.Send(helpMsg, menu)
	})

	bot.Handle("/err", func(ctx telebot.Context) error {
		return telebot.NewError(500, "fdfdfdfd") //ctx.Send(helpMsg, menu)
	})

	bot.Handle(&app.BtnPrint, b.handler.ShowFiles)
	bot.Handle("/print", b.handler.ShowFiles)

	bot.Handle(&app.BtnClear, b.handler.ClearFiles)
	bot.Handle("/clear", b.handler.ClearFiles)

	bot.Handle(&app.BtnMerge, b.handler.Merge)
	bot.Handle("/merge", b.handler.Merge)

	bot.Start()
	return nil
}
