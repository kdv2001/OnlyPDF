package repositories

import "gopkg.in/telebot.v3"

type FilesRepositories interface {
	Add(userName string, document telebot.Document) error
	Update() error
	Get(userName string) ([]telebot.Document, error)
	Delete(userName string) error
}
