package impl

import (
	"OnlyPDF/app/repositories"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"gopkg.in/telebot.v3"
	"strings"
)

type FileUseCaseImpl struct {
	filesRepo repositories.FilesRepositories
}

func CreateFileUseCase(repo repositories.FilesRepositories) *FileUseCaseImpl {
	return &FileUseCaseImpl{filesRepo: repo}
}

func (impl *FileUseCaseImpl) AddFile(user string, file telebot.Document) error {
	err := impl.filesRepo.Add(user, file)
	if err != nil {
		return err
	}
	return nil
}

func (impl *FileUseCaseImpl) MergeFiles(user string, filesNames []string) (string, error) {
	resultName := user + "/" + user + "_result.pdf"
	err := api.MergeCreateFile(filesNames, resultName, nil)
	if err != nil {
		return "", telebot.ErrInternal
	}
	return resultName, nil
}
func (impl *FileUseCaseImpl) ClearFiles(user string) error {
	err := impl.filesRepo.Delete(user)
	if err != nil {
		return err
	}
	return nil
}

func (impl *FileUseCaseImpl) GetFilesNames(user string) (string, error) {
	documents, err := impl.filesRepo.Get(user)
	if err != nil {
		return "", err
	}
	var filesNames []string
	for _, val := range documents {
		filesNames = append(filesNames, val.FileName)
	}
	return strings.Join(filesNames, "\n"), nil
}

func (impl *FileUseCaseImpl) GetFilesIds(user string) ([]string, error) {
	files, err := impl.filesRepo.Get(user)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, val := range files {
		ids = append(ids, val.FileID)
	}
	return ids, nil
}
