package usecase

import (
	"OnlyPDF/app/models"
	"context"
)

type FilesUseCases interface {
	AddFile(user string, file models.File) error
	ConvertFiles(ctx context.Context, userId, resultFileName string, needMerge bool) (string, error)
	ClearFiles(user string) error
	GetFilesNames(user string) (string, error)
	GetFilesIds(user string) ([]string, error)
}
