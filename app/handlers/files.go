package handlers

import (
	"OnlyPDF/app"
	"OnlyPDF/app/usecase"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/telebot.v3"
)

type Handlers struct {
	useCase usecase.FilesUseCases
}

func CreateHandlers(useCase usecase.FilesUseCases) *Handlers {
	return &Handlers{useCase: useCase}
}

func (h *Handlers) AddFiles(ctx telebot.Context) error {
	menu := app.ReturnMainMenu()
	document := ctx.Message().Document

	userId := strconv.FormatInt(ctx.Message().Sender.ID, 10)
	if !strings.Contains(document.MIME, "pdf") {
		ctx.Send("Не поддерживаемый формат файла.", menu)
		return telebot.NewError(http.StatusNotFound, "Bad request: not supported format")
	}

	if err := h.useCase.AddFile(userId, *document); err != nil {
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

	resultNameOnDisk, err := h.useCase.MergeFiles(userId, "")
	if err != nil {
		ctx.Send("Не могу объединить файлы.", menu)
		return telebot.NewError(http.StatusInternalServerError, "Can't send file. Err: "+err.Error())
	}

	if len(ctx.Message().Payload) > 0 {
		resultNameOnDisk = ctx.Message().Payload + ".pdf"
	}

	file := &telebot.Document{FileName: resultNameOnDisk, File: telebot.FromDisk(resultNameOnDisk), MIME: "pdf"}
	if _, err = bot.Send(ctx.Recipient(), file); err != nil {
		ctx.Send("Не могу объединить файлы.", menu)
		return telebot.NewError(http.StatusInternalServerError, "Can't send file. Err: "+err.Error())
	}

	if err = h.useCase.ClearFiles(userId); err != nil {
		ctx.Send("Не могу объединить файлы.", menu)
		return telebot.NewError(http.StatusInternalServerError, "can't remove folder. Err: "+err.Error())
	}

	return nil
}

func (h *Handlers) MergeCommand(ctx telebot.Context) error {
	menu := app.ReturnMainMenu()
	bot := ctx.Bot()
	userId := strconv.FormatInt(ctx.Sender().ID, 10)

	resultFileName := ""
	if len(ctx.Args()) > 0 {
		resultFileName = ctx.Args()[0]
	}

	resultNameOnDisk, err := h.useCase.MergeFiles(userId, resultFileName)
	if err != nil {
		ctx.Send("Не могу объединить файлы.", menu)

		return telebot.NewError(http.StatusInternalServerError, "Can't send file. Err: "+err.Error())
	}

	file := &telebot.Document{FileName: resultFileName, File: telebot.FromDisk(resultNameOnDisk), MIME: "pdf"}

	_, err = bot.Send(ctx.Recipient(), file)
	if err != nil {
		ctx.Send("Не могу объединить файлы.", menu)

		return telebot.NewError(http.StatusInternalServerError, "Can't send file. Err: "+err.Error())
	}

	if err = h.useCase.ClearFiles(userId); err != nil {
		ctx.Send("Не могу объединить файлы.", menu)
		return telebot.NewError(http.StatusInternalServerError, "can't remove folder. Err: "+err.Error())
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
	if err := h.useCase.ClearFiles(userId); err != nil {
		ctx.Send("Не могу очистить очередь файлов.", menu)
		return err
	}

	ctx.Send("Текущая очередь очищена.")

	return nil
}
