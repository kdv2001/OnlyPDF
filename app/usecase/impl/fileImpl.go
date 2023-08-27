package impl

import (
	"OnlyPDF/app/models"
	"OnlyPDF/app/repositories"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const dirPermission = 0777

type FileDownLoader interface {
	DownloadFile(fileId, localFileName string) error
}

type pdfTools interface {
	ConvertFiles(ctx context.Context, files []string, resultFileName string, needMerge bool) error
}

type FileUseCaseImpl struct {
	filesRepo repositories.FilesRepositories
	pdfTools  pdfTools
	loader    FileDownLoader
}

func CreateFileUseCase(repo repositories.FilesRepositories, loader FileDownLoader, pdf pdfTools) *FileUseCaseImpl {
	return &FileUseCaseImpl{
		filesRepo: repo,
		loader:    loader,
		pdfTools:  pdf,
	}
}

func (impl *FileUseCaseImpl) AddFile(user string, file models.File) error {
	err := impl.filesRepo.Add(user, file)
	if err != nil {
		return err
	}
	return nil
}

func (impl *FileUseCaseImpl) ConvertFiles(ctx context.Context, userId, resultFileName string, needMerge bool) (string, error) {
	if resultFileName == "" {
		resultFileName = fmt.Sprint(userId, "_result")
	}

	if needMerge {
		resultFileName = fmt.Sprint(resultFileName, ".pdf")
	} else {
		resultFileName = fmt.Sprint(resultFileName, ".zip")
	}

	resultFileName = fmt.Sprint(userId, "/", resultFileName)

	os.RemoveAll(fmt.Sprint("./", userId))

	if _, err := os.Stat(userId); os.IsNotExist(err) {
		if err = os.Mkdir(userId, dirPermission); err != nil {
			return "", err
		}
	}

	files, err := impl.filesRepo.Get(userId)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", fmt.Errorf("очередь пуста")
	}

	downloadFilesNames := make([]string, 0, len(files))

	for _, val := range files {
		downloadFileName := fmt.Sprint(userId, "/", val.Name)
		downloadFilesNames = append(downloadFilesNames, downloadFileName)
		if err = impl.loader.DownloadFile(val.Id, downloadFileName); err != nil {
			return "", err
		}
	}

	if err = impl.pdfTools.ConvertFiles(ctx, downloadFilesNames, resultFileName, needMerge); err != nil {
		return "", fmt.Errorf("error merge file %v", err)
	}

	return resultFileName, nil
}

func (impl *FileUseCaseImpl) ClearFiles(userId string) error {
	if err := impl.filesRepo.Delete(userId); err != nil {
		return err
	}

	os.RemoveAll(fmt.Sprint("./", userId))

	return nil
}

func (impl *FileUseCaseImpl) GetFilesNames(user string) (string, error) {
	documents, err := impl.filesRepo.Get(user)
	if err != nil {
		return "", err
	}

	filesNames := make([]string, 0)
	for idx, val := range documents {
		newName := strconv.Itoa(idx+1) + ") " + val.Name
		filesNames = append(filesNames, newName)
	}
	return strings.Join(filesNames, "\n"), nil
}

func (impl *FileUseCaseImpl) GetFilesIds(user string) ([]string, error) {
	files, err := impl.filesRepo.Get(user)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0)
	for _, val := range files {
		ids = append(ids, val.Id)
	}

	return ids, nil
}
