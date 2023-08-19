package tgbot

import (
	"context"
	"log"
	"time"

	"github.com/V1taly5/APIYDisk/internal/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (telegramBot *TelegramBot) handleCommand(chatID int, update tgbotapi.Update) {
	command := update.Message.CommandWithAt()
	switch command {
	case "start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я бот для  отправки файлов на Я.Диск.")
		_, err := telegramBot.API.Send(msg)
		if err != nil {
			log.Println(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		telegramBot.UseCase.CreateUser(ctx, int(update.Message.Chat.ID))
	// case "register":
	// 	find := telegramBot.findUser(chatID)
	// 	if !find {
	// 		telegramBot.State.SetState("start")
	// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Начинаем регистрацию. Введите ваше имя:")
	// 		_, err := telegramBot.API.Send(msg)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 		err = telegramBot.State.Event(ctx, "register")
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 	} else {
	// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы уже зарегестрированы")
	// 		telegramBot.API.Send(msg)
	// 	}
	case "open":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyMarkup = numericKeyboard
		if _, err := telegramBot.API.Send(msg); err != nil {
			panic(err)
		}
	case "set":

	}

}

func (telegramBot *TelegramBot) handeletMessage(chatID int, update tgbotapi.Update) {
	if telegramBot.State.Current() == "registering" {
		var user entity.User
		user.ChatID = chatID
		user.YandexDiskToken = update.Message.Text
		telegramBot.State.Event(ctx, "cencel")
		user.State = telegramBot.State.Current()
		// err := telegramBot.Repo.Create(&user)
		// if err != nil {
		// 	msg := tgbotapi.NewMessage(int64(chatID), "Произошла ошибка! = Бот работает не правильно")
		// 	telegramBot.API.Send(msg)
		// 	return
		// }

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Спасибо за регистрацию!")
		telegramBot.API.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не могу обработать это сообщение!")
		telegramBot.API.Send(msg)
	}
}
