package repositories

import (
	"OnlyPDF/app/models"
)

type FilesRepositories interface {
	Add(userName string, document models.File) error
	Update() error
	Get(userId string) ([]models.File, error)
	Delete(userName string) error
}
