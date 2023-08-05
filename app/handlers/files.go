package handlers

import (
	"OnlyPDF/app"
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
	menu := app.ReturnMainMenu()
	document := ctx.Message().Document
	userId := strconv.FormatInt(ctx.Message().Sender.ID, 10)
	if !strings.Contains(document.MIME, "pdf") {
		ctx.Send("Не поддерживаемый формат файла.", menu)
		return telebot.NewError(404, "Bad request: not supported format")
	}
	err := h.useCase.AddFile(userId, *document)
	if err != nil {
		ctx.Send("Не удалось добавить файл.", menu)
		return err
	}
	ctx.Send("Файл добавлен.", menu)
	return nil
}

func (h *Handlers) Merge(ctx telebot.Context) error {
	menu := app.ReturnMainMenu()
	bot := ctx.Bot()
	userId := strconv.FormatInt(ctx.Message().Sender.ID, 10)
	files, err := h.useCase.GetFilesIds(userId)
	if err == telebot.ErrNotFound {
		fmt.Println(err)
		ctx.Send("Файлы в очереди отсутствуют.", menu)
		return err
	}
	if err != nil {
		fmt.Println(err)
		ctx.Send("Не удалось получить список файлов.", menu)
		return err
	}
	if len(files) <= 1 {
		ctx.Send("Недостаточно файлов.", menu)
		return nil
	}

	ctx.Send("Начинаю объединять файлы.", menu)

	if _, err = os.Stat(userId); !os.IsNotExist(err) {
		os.RemoveAll("./" + userId + "/")
	}

	err = os.Mkdir(userId, 0777)
	if err != nil {
		fmt.Println(err)
		ctx.Send("Не могу объединить файлы.", menu)
		return telebot.NewError(500, "can't create folder")
	}
	defer os.RemoveAll("./" + userId + "/")
	var fileNames []string
	errChan := make(chan error)
	quitChan := make(chan bool)
	for idx, val := range files {
		fileName := userId + "/" + strconv.Itoa(idx)
		fileNames = append(fileNames, fileName)
		document := telebot.File{FileID: val}
		go downloadFile(document, bot, fileName, errChan, quitChan)
	}
	downloadCount := 0
LOOP:
	for {
		select {
		case err = <-errChan:
			if err != nil {
				close(quitChan)
				ctx.Send("Не могу объединить файлы.")
				return err
			} else {
				downloadCount++
				if downloadCount == len(files) {
					break LOOP
				}
			}
		}
	}
	resultNameOnDisk, err := h.useCase.MergeFiles(userId, fileNames)
	nameForFile := resultNameOnDisk
	if len(ctx.Message().Payload) > 0 {
		nameForFile = ctx.Message().Payload + ".pdf"
	}
	file := &telebot.Document{FileName: nameForFile, File: telebot.FromDisk(resultNameOnDisk), MIME: "pdf"}
	_, err = bot.Send(ctx.Recipient(), file)
	if err != nil {
		ctx.Send("Не могу объединить файлы.", menu)
		return telebot.NewError(500, "Can't send file. Err: "+err.Error())
	}
	h.useCase.ClearFiles(userId)
	err = os.Remove(resultNameOnDisk)
	if err != nil {
		fmt.Println(err)
		ctx.Send("Не могу объединить файлы.", menu)
		return telebot.NewError(500, "can't remove folder. Err: "+err.Error())
	}
	return nil
}

func (h *Handlers) ShowFiles(ctx telebot.Context) error {
	menu := app.ReturnMainMenu()
	userId := strconv.FormatInt(ctx.Message().Sender.ID, 10)
	filenames, err := h.useCase.GetFilesNames(userId)
	if err == telebot.ErrNotFound {
		ctx.Send("Файлы в текущей очереди отсутствуют.", menu)
		return nil
	}
	ctx.Send("Текущая очередь:\n" + filenames)
	return nil
}

func (h *Handlers) ClearFiles(ctx telebot.Context) error {
	menu := app.ReturnMainMenu()
	userId := strconv.FormatInt(ctx.Message().Sender.ID, 10)
	err := h.useCase.ClearFiles(userId)
	if err != nil {
		ctx.Send("Не могу очистить очередь файлов.", menu)
		return err
	}
	ctx.Send("Текущая очередь очищена.")
	return nil
}

func downloadFile(doc telebot.File, b *telebot.Bot, file string, errCh chan error, quitChan chan bool) {
	err := b.Download(&doc, file)
	if err != nil {
		fmt.Println(err)
		err = telebot.NewError(500, "can't download file. Err: "+err.Error())
	}
	select {
	case errCh <- err:
		return
	case <-quitChan:
		return
	}
}
