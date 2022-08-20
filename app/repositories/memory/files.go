package memory

import "gopkg.in/telebot.v3"

type FilesMemory struct {
	dataBase map[string][]telebot.Document
}

func CreateFilesPostgresInMemory() (*FilesMemory, error) {
	db := make(map[string][]telebot.Document)
	return &FilesMemory{dataBase: db}, nil
}

func (db *FilesMemory) Add(userName string, document telebot.Document) error {
	sumFileSize := 0
	for _, val := range db.dataBase[userName] {
		sumFileSize += val.FileSize
	}
	sumFileSize += document.FileSize
	if sumFileSize >= 50000000 {
		return telebot.ErrCantUploadFile
	}
	db.dataBase[userName] = append(db.dataBase[userName], document)
	return nil
}

func (db *FilesMemory) Get(userName string) ([]telebot.Document, error) {
	files, ok := db.dataBase[userName]
	if !ok {
		return []telebot.Document{}, telebot.ErrNotFound
	}
	return files, nil
}

func (db *FilesMemory) Update() error {
	return nil
}

func (db *FilesMemory) Delete(userName string) error {
	delete(db.dataBase, userName)
	return nil
}
