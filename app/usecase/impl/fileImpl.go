package impl

type FileUseCase struct {
}

func CreateFileUseCase() FileUseCase {
	return FileUseCase{}
}

func (f *FileUseCase) Merge([][]byte) []byte {
	return []byte("0")
}
