package impl

import (
	"OnlyPDF/app/repositories"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"gopkg.in/telebot.v3"
	"os"
	"strconv"
	"strings"
)

type FileDownLoader interface {
	DownloadFile(fileId, localFileName string) error
}

type FileUseCaseImpl struct {
	filesRepo repositories.FilesRepositories
	loader    FileDownLoader
}

func CreateFileUseCase(repo repositories.FilesRepositories, loader FileDownLoader) *FileUseCaseImpl {
	return &FileUseCaseImpl{filesRepo: repo, loader: loader}
}

func (impl *FileUseCaseImpl) AddFile(user string, file telebot.Document) error {
	err := impl.filesRepo.Add(user, file)
	if err != nil {
		return err
	}
	return nil
}

func (impl *FileUseCaseImpl) MergeFiles(userId, resultFileName string) (string, error) {
	if resultFileName == "" {
		resultFileName = fmt.Sprint(userId, "_result.pdf")
	}

	if _, err := os.Stat(userId); !os.IsNotExist(err) {
		os.RemoveAll("./" + userId)
	} else {
		if err = os.Mkdir(userId, 0777); err != nil {
			return "", err
		}
	}

	resultFileName = fmt.Sprint(userId, "/", resultFileName)

	files, err := impl.filesRepo.Get(userId)
	if err != nil {
		return "", err
	}

	downloadFilesNames := make([]string, 0, len(files))

	for _, val := range files {
		downloadFileName := fmt.Sprint(userId, "/", val.FileName)
		downloadFilesNames = append(downloadFilesNames, downloadFileName)
		if err = impl.loader.DownloadFile(val.FileID, downloadFileName); err != nil {
			return "", err
		}
	}

	if err = api.MergeCreateFile(downloadFilesNames, resultFileName, nil); err != nil {
		return "", fmt.Errorf("error merge file")
	}

	return resultFileName, nil
}
func (impl *FileUseCaseImpl) ClearFiles(user string) error {
	if err := impl.filesRepo.Delete(user); err != nil {
		return err
	}

	os.RemoveAll("./" + user)

	return nil
}

func (impl *FileUseCaseImpl) GetFilesNames(user string) (string, error) {
	documents, err := impl.filesRepo.Get(user)
	if err != nil {
		return "", err
	}
	var filesNames []string
	for idx, val := range documents {
		newName := strconv.Itoa(idx+1) + ") " + val.FileName
		filesNames = append(filesNames, newName)
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
