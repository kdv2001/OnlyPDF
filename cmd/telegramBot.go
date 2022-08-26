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

var defaultMsg = "Привет, это OnlyPDFBot, предназначенный " +
	"для соединения pdf файлов. Для более подробной информации о " +
	"доступных командах воспользуйся командой /help или нажми на" +
	" кнопку помощь."

var helpMsg = "Как оно работает? \n" +
	"Каждый отправленный пдф файл в чат с ботом добавляется в очередь, " +
	"а затем вся очередь объединяется cверху вниз.\n\n" +
	"Доступные команды: \n" +
	"Объединить (/merge) - объединение файлов \n" +
	"отобразить (/print) - отображение загруженных файлов \n" +
	"очистить (/clear) - очистка загруженных файлов в очередью \n" +
	"помощь (/help) - помощь. \n" +
	"На данный момент бот отправить итоговый файла максимум размером 50 Мб, \n" +
	"а получить файл максимум 20 Мб."

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
		return ctx.Send(defaultMsg, menu)
	})
	bot.Handle(&btnHelp, func(ctx telebot.Context) error {
		return ctx.Send(helpMsg, menu)
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
	bot.Handle("/help", func(ctx telebot.Context) error {
		return ctx.Send(helpMsg, menu)
	})
	bot.Handle("/print", b.handler.ShowFiles)
	bot.Handle("/clear", b.handler.ClearFiles)
	bot.Handle("/merge", func(ctx telebot.Context) error {
		err := b.handler.Merge(ctx, bot)
		if err != nil {
			return err
		}
		return nil
	})
	bot.Start()
}
