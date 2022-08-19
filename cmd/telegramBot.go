package main

import (
	"OnlyPDF/app/handlers"
	"OnlyPDF/app/repositories/postgress"
	"OnlyPDF/app/usecase/impl"
	"fmt"
	//_ "github.com/jackc/pgx/stdlib"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"gopkg.in/telebot.v3"
	"os"
	"time"
)

type OnlyPDFBot struct {
	handler *handlers.Handlers
}

func CreateOnlyPDFBot() (OnlyPDFBot, error) {
	conn, err := sqlx.Open("pgx", "")
	if err != nil {
		return OnlyPDFBot{}, err
	}
	repo, err := postgress.CreateFilesPostgres(conn)
	if err != nil {
		return OnlyPDFBot{}, err
	}
	fileUseCase := impl.CreateFileUseCase(repo)
	handler := handlers.CreateHandlers(fileUseCase)
	return OnlyPDFBot{handler}, nil
}

func (b *OnlyPDFBot) StartListenAndServ() {
	bot, err := telebot.NewBot(telebot.Settings{Token: os.Getenv("OnlyPDFBotToken"), Poller: &telebot.LongPoller{Timeout: 10 * time.Second}})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//bot.Use(middleware.Logger())
	//bot.Handle(telebot.OnDocument, func(ctx telebot.Context) error {
	//	fmt.Println(ctx.Message().Document.FileID)
	//	fmt.Println(ctx.Message().ID)
	//	return nil
	//})

	bot.Handle(telebot.OnDocument, b.handler.AddFiles)
	bot.Handle("/t", b.handler.ShowFiles)
	bot.Handle("/merge", func(ctx telebot.Context) error {
		err := b.handler.Merge(ctx, bot)
		if err != nil {
			return err
		}
		return nil
	})

	bot.Start()
}
