package usecase

type FilesUseCases interface {
	Merge([][]byte) []byte
}
