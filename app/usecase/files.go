package usecase

import "gopkg.in/telebot.v3"

type FilesUseCases interface {
	AddFile(user string, file telebot.Document) error
	MergeFiles(userId, resultFileName string) (string, error)
	ClearFiles(user string) error
	GetFilesNames(user string) (string, error)
	GetFilesIds(user string) ([]string, error)
}
