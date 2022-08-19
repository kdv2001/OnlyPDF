package usecase

import "gopkg.in/telebot.v3"

type FilesUseCases interface {
	AddFile(file telebot.File) error
}
