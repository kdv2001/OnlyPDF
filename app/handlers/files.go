package handlers

import (
	"OnlyPDF/app/usecase"
	"fmt"
	"gopkg.in/telebot.v3"
	"os"
	"strconv"
	"strings"
)

type Handlers struct {
	useCase usecase.FilesUseCases
	mock    map[string][]string
}

func CreateHandlers(useCase usecase.FilesUseCases) *Handlers {
	mock := make(map[string][]string)
	return &Handlers{useCase: useCase, mock: mock}
}

func (h *Handlers) AddFiles(ctx telebot.Context) error {
	document := ctx.Message().Document
	user := ctx.Message().Sender.Username
	if !strings.Contains(document.MIME, "pdf") {
		ctx.Send("Не поддерживаемый формат файла.")
		return telebot.NewError(404, "Bad request: not not supported")
	}
	err := h.useCase.AddFile(user, *document)
	if err != nil {
		ctx.Send("Не удалось добавить файл.")
		return err
	}
	ctx.Send("Файл добавлен.")
	return nil
}

func (h *Handlers) Merge(ctx telebot.Context, bot *telebot.Bot) error {
	user := ctx.Message().Sender
	files, err := h.useCase.GetFilesIds(user.Username)
	if err == telebot.ErrNotFound {
		ctx.Send("Файлы в очереди отсутствуют.")
		return err
	}
	if err != nil {
		ctx.Send("Не удалось получить список файлов.")
		return err
	}
	if len(files) <= 1 {
		ctx.Send("Недостаточно файлов.")
		return nil
	}
	err = os.Mkdir(user.Username, 0777)
	if err != nil {
		ctx.Send("Не могу объединить файлы.")
		return telebot.ErrInternal
	}
	defer os.RemoveAll("./" + user.Username + "/")
	fmt.Println("/" + user.Username)
	var fileNames []string
	for idx, val := range files {
		fileName := user.Username + "/" + strconv.Itoa(idx)
		fileNames = append(fileNames, fileName)
		document := telebot.File{FileID: val}
		err = bot.Download(&document, fileName)
		if err != nil {
			ctx.Send("Не могу объединить файлы.")
			return err
		}
	}
	resultName, err := h.useCase.MergeFiles(user.Username, fileNames)
	file := &telebot.Document{FileName: resultName, File: telebot.FromDisk(resultName), MIME: "pdf"}
	_, err = bot.Send(ctx.Recipient(), file)
	if err != nil {
		ctx.Send("Не могу объединить файлы.")
		return err
	}
	h.useCase.ClearFiles(user.Username)
	err = os.Remove(resultName)
	if err != nil {
		ctx.Send("Не могу объединить файлы.")
		return telebot.ErrInternal
	}
	return nil
}

func (h *Handlers) ShowFiles(ctx telebot.Context) error {
	user := ctx.Message().Sender.Username
	filenames, err := h.useCase.GetFilesNames(user)
	if err == telebot.ErrNotFound {
		ctx.Send("Файлы в очереди отсутствуют.")
		return err
	}
	if err != nil {
		ctx.Send("Не могу отобразить файлы.")
		return err
	}
	ctx.Send(filenames)
	return nil
}

func (h *Handlers) ClearFiles(ctx telebot.Context) error {
	user := ctx.Message().Sender.Username
	err := h.useCase.ClearFiles(user)
	if err != nil {
		ctx.Send("Не могу очистить файлы.")
		return err
	}
	ctx.Send("Очередь очищена.")
	return nil
}
