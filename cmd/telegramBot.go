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
	menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnPrint := menu.Text("Print")
	btnMerge := menu.Text("Merge")
	btnClear := menu.Text("Clear")
	btnHelp := menu.Text("Help")

	menu.Reply(
		menu.Row(btnPrint),
		menu.Row(btnMerge),
		menu.Row(btnClear),
		menu.Row(btnHelp),
	)

	bot.Handle(telebot.OnDocument, b.handler.AddFiles)
	bot.Handle("/start", func(ctx telebot.Context) error {
		return ctx.Send("Hello!", menu)
	})
	bot.Handle(&btnHelp, func(ctx telebot.Context) error {
		msg := "Как оно работает? Каждый отправленный пдф файл в чат с ботом добавляется в очередь," +
			" затем вся очередь объединяется снизу вверх.\n" +
			"merge - объединение файлов \n" +
			"print - отображение загруженных файлов(пока что теграм айди файлов)\n" +
			"clear - очистка загруженных файлов в очередью \n"

		return ctx.Send(msg, menu)
	})
	bot.Handle(&btnPrint, b.handler.ShowFiles)
	bot.Handle(&btnClear, b.handler.ClearFiles)
	bot.Handle(&btnMerge, func(ctx telebot.Context) error {
		err := b.handler.Merge(ctx, bot)
		if err != nil {
			return err
		}
		return nil
	})

	bot.Start()
}
