package main

import (
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

func (b *OnlyPDFBot) StartListenAndServ() {
	bot, err := telebot.NewBot(telebot.Settings{Token: os.Getenv("OnlyPDFBotToken"), Poller: &telebot.LongPoller{Timeout: 10 * time.Second}})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//bot.Use(middleware.Logger())
	menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnPrint := menu.Text("🖨 Отобразить")
	btnMerge := menu.Text("💾 Объединить")
	btnClear := menu.Text("🧫 Очистить")
	btnHelp := menu.Text("🆘 Помощь")

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
		msg := "Как оно работает? \n" +
			"Каждый отправленный пдф файл в чат с ботом добавляется в очередь, а затем вся очередь объединяется cверху вниз.\n\n" +
			"Доступные команды: \n" +
			"Объединить (/merge) - объединение файлов \n" +
			"print - отображение загруженных файлов \n" +
			"clear - очистка загруженных файлов в очередью"

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
