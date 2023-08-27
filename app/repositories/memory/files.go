package memory

import (
	"OnlyPDF/app/models"
	"errors"
	"sync"

	"gopkg.in/telebot.v3"
)

const maxFileSize = 50000000

type FilesMemory struct {
	syncDataBase *sync.Map
}

func CreateFilesPostgresInMemory() (*FilesMemory, error) {
	syncDb := sync.Map{}
	return &FilesMemory{syncDataBase: &syncDb}, nil
}

func (db *FilesMemory) Add(userName string, document models.File) error {
	var sumFileSize int64
	fileSliceAny, ok := db.syncDataBase.Load(userName)

	fileSlice, okConvert := fileSliceAny.([]models.File)
	if ok || !okConvert {
		for _, val := range fileSlice {
			sumFileSize += val.Size
		}
	}

	sumFileSize += document.Size
	if sumFileSize >= maxFileSize {
		return telebot.ErrCantUploadFile
	}

	db.syncDataBase.Store(userName, append(fileSlice, document))

	return nil
}

func (db *FilesMemory) Get(userName string) ([]models.File, error) {
	fileSliceAny, ok := db.syncDataBase.Load(userName)
	if !ok {
		// TODO refactor
		return nil, telebot.ErrNotFound
	}

	fileSlice, ok := fileSliceAny.([]models.File)
	if !ok {
		return nil, errors.New("bad type assertion")
	}

	return fileSlice, nil
}

func (db *FilesMemory) Update() error {
	return nil
}

func (db *FilesMemory) Delete(userName string) error {
	db.syncDataBase.Delete(userName)

	return nil
}
