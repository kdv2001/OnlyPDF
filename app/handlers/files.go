package handlers

import (
	"OnlyPDF/app/usecase"
	"fmt"
	"github.com/unidoc/unipdf/v3/model"
	"gopkg.in/telebot.v3"
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
	fmt.Println(document.MIME)
	if !strings.Contains(document.MIME, "pdf") {
		return telebot.ErrWrongFileID
	}
	err := h.useCase.AddFile(document.File)
	if err != nil {
		return telebot.ErrWrongFileID
	}
	h.mock[user] = append(h.mock[user], document.FileID)
	return nil
}

func (h *Handlers) Merge(ctx telebot.Context, bot *telebot.Bot) error {
	user := ctx.Message().Sender.Username
	writer := model.NewPdfWriter()
	var fileNames []string
	for idx, val := range h.mock[user] {
		fileName := user + "_" + strconv.Itoa(idx)
		fileNames = append(fileNames, fileName)
		document := telebot.File{FileID: val}

		err := bot.Download(&document, fileName)
		if err != nil {
			return telebot.ErrWrongFileID
		}
		reader, _, err := model.NewPdfReaderFromFile(fileName, nil)
		if err != nil {
			fmt.Println(err)
		}
		for _, page := range reader.PageList {
			writer.AddPage(page)
		}
	}
	resultName := user + "_result.pdf"
	err := writer.WriteToFile(resultName)
	if err != nil {
		fmt.Println(err)
	}
	file := &telebot.Document{FileName: resultName, File: telebot.FromDisk(resultName), MIME: "pdf"}
	_, err = bot.Send(ctx.Recipient(), file)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (h *Handlers) ShowFiles(ctx telebot.Context) error {
	user := ctx.Message().Sender.Username
	files, findFlag := h.mock[user]
	if !findFlag {
		ctx.Send("files not found")
		return nil
	}
	ctx.Send(strings.Join(files, "/n"))
	return nil
}
