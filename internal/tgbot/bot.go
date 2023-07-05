package tgbot

import (
	"context"
	"fmt"
	"log"

	"github.com/V1taly5/APIYDisk/internal/entity"
	"github.com/V1taly5/APIYDisk/internal/infrastructure/repository"
	"github.com/V1taly5/APIYDisk/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/looplab/fsm"
)

var ctx = context.Background()

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	),
)

type TelegramBot struct {
	API     *tgbotapi.BotAPI
	Updates tgbotapi.UpdatesChannel
	Repo    entity.UserRepository
	State   *fsm.FSM
}

func (telegramBot *TelegramBot) Init(repo entity.UserRepository) {
	botAPI, err := tgbotapi.NewBotAPI(TG_Bot_Token)
	if err != nil {
		log.Fatal(err)
	}
	telegramBot.API = botAPI
	telegramBot.API.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	botUpdates := telegramBot.API.GetUpdatesChan(updateConfig)

	telegramBot.Updates = botUpdates

	telegramBot.Repo = repo

}

func (telegramBot *TelegramBot) ReceiveUpdates(ctx context.Context) {

	telegramBot.State = fsm.NewFSM(
		"start",
		fsm.Events{
			{Name: "register", Src: []string{"start"}, Dst: "registering"},
			{Name: "cancel", Src: []string{"registering"}, Dst: "start"},
		},
		fsm.Callbacks{},
	)
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-telegramBot.Updates:
			telegramBot.analyzeUpdate(update, telegramBot.State)
			break
		}
	}
}

func (telegramBot *TelegramBot) analyzeUpdate(update tgbotapi.Update, newFSM *fsm.FSM) {
	if update.CallbackQuery != nil {
		// Respond to the callback query, telling Telegram to show the user
		// a message with the data received.
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		if _, err := telegramBot.API.Request(callback); err != nil {
			panic(err)
		}

		// And finally, send a message containing the data received.
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
		if _, err := telegramBot.API.Send(msg); err != nil {
			panic(err)
		}
		return
	}
	// TODO
	chatID := update.Message.Chat.ID

	if telegramBot.findUser(int(chatID)) {
		telegramBot.analyzeUser(update)
		// user exists
	} else {
		// message from new user
		msg := tgbotapi.NewMessage(chatID, "Пожалуйста зарегестрируйтесь!")
		telegramBot.API.Send(msg)
	}
	if update.Message.IsCommand() {
		telegramBot.handleCommand(int(chatID), update)
	} else if update.Message.Document != nil {
		msg := tgbotapi.NewMessage(chatID, "This is Document")
		telegramBot.API.Send(msg)
	} else if update.Message.Text != "" {
		telegramBot.handeletMessage(int(chatID), update)
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "бот работает только с файласм")
		msg.ReplyToMessageID = update.Message.MessageID
		telegramBot.API.Send(msg)

	}

	// check for message type

	// if update.Message.Text != "" {
	// 	if telegramBot.State.Current() == "Start" {
	// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Начинаем регистрацию. Введите ваше имя:")
	// 		_, err := telegramBot.API.Send(msg)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 		err = telegramBot.State.Event(ctx, "register")
	// 		telegramBot.ActiveContactRequests[chatID] = "registering"
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 	}
	// 	if telegramBot.State.Current() == "registering" {
	// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы уже зарегистрированы.")
	// 		_, err := telegramBot.API.Send(msg)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 	}
	// }
}

func (telegramBot *TelegramBot) analyzeUser(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	user, err := telegramBot.Repo.GetUser(chatID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Произошла ошибка! Бот работает не правильно")
		telegramBot.API.Send(msg)
		return
	}
	telegramBot.State.SetState(user.State)
	telegramBot.State.Event(ctx, "cancel")
	fmt.Println(telegramBot.State.Current())
	// TODO analyxeDate
}

func (telegramBot *TelegramBot) findUser(chatID int) bool {
	find, err := telegramBot.Repo.Find(chatID)
	if err != nil {
		msg := tgbotapi.NewMessage(int64(chatID), "Произошла ошибка! Бот может работать не правильнго!")
		telegramBot.API.Send(msg)
	}
	return find
}

// func StartTelegramBot(token string, debugMod bool, disk *repository.YandexDiskAPI) error {
// 	bot, err := tgbotapi.NewBotAPI(token)
// 	if err != nil {
// 		return err
// 	}

// 	bot.Debug = debugMod

// 	updateConfig := tgbotapi.NewUpdate(0)
// 	updateConfig.Timeout = 60
// 	updates := bot.GetUpdatesChan(updateConfig)

// 	docUseCase := usecase.NewDocumentUseCase(*disk)

// 	for update := range updates {
// 		if update.Message == nil {
// 			continue
// 		}
// 		fmt.Println(update)
// 		if update.Message.Document != nil {
// 			fmt.Println("!!!!!!!!!!!! ", update.Message.Photo, " !!!!!!!!!")
// 			// handelerFile(update, bot, disk)
// 			handelerDoc(update, bot, docUseCase)
// 		} else if update.Message.Text != "" {
// 			msg := tgbuser, err := Connection.GetUser(chatID)  otapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
// 			switch update.Message.Text {
// 			case "open":
// 				msg.ReplyMarkup = numericKeyboard_2
// 				bot.Send(msg)
// 			case "setCurentPath":
// 			}
// 		} else {
// 			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "бот работает только с файласм")
// 			msg.ReplyToMessageID = update.Message.MessageID
// 			bot.Send(msg)
// 		}

// 	}
// 	return nil
// }

func handelerDoc(update tgbotapi.Update, bot *tgbotapi.BotAPI, docUseCase *usecase.DocumentUseCase) {
	var fileConfig tgbotapi.FileConfig

	fileConfig.FileID = update.Message.Document.FileID
	file, err := bot.GetFile(fileConfig)
	if err != nil {
		return
	}

	fmt.Println("@@@@@@@@@@@@@")
	url := file.Link(bot.Token)
	fmt.Println(url)
	fmt.Println("@@@@@@@@@@@@@")

	document, err := docUseCase.UploadDocument(url, "api4")
	if err != nil {
		panic(err)
	}
	fmt.Println("Response to yandex disk: ", document)
}

func handelerFile(update tgbotapi.Update, bot *tgbotapi.BotAPI, disk *repository.YandexDiskAPI) {
	var fileConfig tgbotapi.FileConfig

	fileConfig.FileID = update.Message.Document.FileID
	file, err := bot.GetFile(fileConfig)
	if err != nil {
		return
	}

	fmt.Println("@@@@@@@@@@@@@")
	url := file.Link(bot.Token)
	fmt.Println(url)
	fmt.Println("@@@@@@@@@@@@@")

	_, err = disk.UploadFileLink(url, "api3")
	if err != nil {
		panic(err)
	}
}
