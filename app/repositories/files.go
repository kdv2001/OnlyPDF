package repositories

type FilesRepositories interface {
	Add() error
	Update() error
	Delete() error
}
