package handlers

import (
	"context"
	"gopkg.in/telebot.v3"
)

type cache interface {
	Set(ctx context.Context)
	Get(ctx context.Context)
}

type Middlewares struct {
	cache cache
}

func StateMiddleware(handlerFunc telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {

		return handlerFunc(ctx)
	}
}

//func stateMiddleware(handlerFunc telebot.HandlerFunc) telebot.HandlerFunc {
//
//	return nil
//}
