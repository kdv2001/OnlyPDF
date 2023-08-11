package telegramFiles

import (
	"gopkg.in/telebot.v3"
)

type TelegramFiles struct {
	bot *telebot.Bot
}

func NewTelegramFiles(bot *telebot.Bot) TelegramFiles {
	return TelegramFiles{
		bot: bot,
	}
}

func (t *TelegramFiles) DownloadFile(fileId, localFileName string) error {
	file := telebot.File{FileID: fileId}
	if err := t.bot.Download(&file, localFileName); err != nil {
		return err
	}

	return nil
}
