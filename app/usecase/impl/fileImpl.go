package impl

import (
	"OnlyPDF/app/repositories"
	"gopkg.in/telebot.v3"
)

type FileUseCaseImpl struct {
	filesRepo repositories.FilesRepositories
}

func CreateFileUseCase(repo repositories.FilesRepositories) *FileUseCaseImpl {
	return &FileUseCaseImpl{filesRepo: repo}
}

func (impl *FileUseCaseImpl) AddFile(file telebot.File) error {
	err := impl.filesRepo.Add()
	if err != nil {
		return err
	}
	return nil
}
