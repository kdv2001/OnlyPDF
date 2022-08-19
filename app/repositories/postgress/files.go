package postgress

import "github.com/jmoiron/sqlx"

type FilesPostgres struct {
	connection *sqlx.DB
}

func CreateFilesPostgres(db *sqlx.DB) (*FilesPostgres, error) {
	_, err := db.Exec("")
	if err != nil {
		return nil, err
	}
	return &FilesPostgres{connection: db}, nil
}

func (db *FilesPostgres) Add() error {

	return nil
}

func (db *FilesPostgres) Update() error {

	return nil
}

func (db *FilesPostgres) Delete() error {

	return nil
}
