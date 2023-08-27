package app

import (
	"gopkg.in/telebot.v3"
)

var (
	menu = &telebot.ReplyMarkup{ResizeKeyboard: true}

	BtnPrint   = menu.Text("🖨 Отобразить")
	BtnMerge   = menu.Text("💾 Объединить")
	BtnClear   = menu.Text("🧫 Очистить")
	BtnHelp    = menu.Text("🆘 Помощь")
	BtnConvert = menu.Text("⚙ Конвертировать")
)

func ReturnMainMenu() *telebot.ReplyMarkup {
	menu.Reply(
		menu.Row(BtnPrint, BtnMerge),
		menu.Row(BtnClear, BtnConvert),
		menu.Row(BtnHelp),
	)
	return menu
}

var DefaultMsg = "Привет, это OnlyPDFBot, предназначенный " +
	"для соединения pdf файлов. Для более подробной информации о " +
	"доступных командах воспользуйся командой /help или нажми на" +
	" кнопку помощь."

var HelpMsg = "Как оно работает? \n" +
	"Каждый отправленный пдф файл в чат с ботом добавляется в очередь, " +
	"а затем вся очередь объединяется cверху вниз.\n\n" +
	"Доступные команды: \n" +
	"Объединить (/merge) - объединение файлов \n" +
	"отобразить (/print) - отображение загруженных файлов \n" +
	"очистить (/clear) - очистка загруженных файлов в очередью \n" +
	"помощь (/help) - помощь. \n" +
	"Чтобы изменить имя конечного файла, используйте /merge filename, где filename - имя конечного файла. \n" +
	"На данный момент бот может отправить итоговый файла максимум размером 50 Мб, \n" +
	"а получить файл максимум 20 Мб."
