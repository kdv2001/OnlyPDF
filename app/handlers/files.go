package handlers

import (
	"OnlyPDF/app/usecase"
	"fmt"
	"gopkg.in/telebot.v3"
)

type Handlers struct {
	useCase *usecase.FilesUseCases
}

func CreateHandlers(useCase usecase.FilesUseCases) Handlers {
	return Handlers{useCase: &useCase}
}

func (h *Handlers) AddFiles(ctx telebot.Context) error {
	messageId := ctx.Message().ID
	i := messageId

	fmt.Println(i)
	return nil
}
