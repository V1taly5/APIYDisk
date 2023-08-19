package main

import (
	"bufio"
	"context"
	"log"
	"os"

	mongo "github.com/V1taly5/APIYDisk/internal/infrastructure/repository/mongo"
	"github.com/V1taly5/APIYDisk/internal/tgbot"
	"github.com/V1taly5/APIYDisk/internal/usecase"
)

func main() {
	// cfg := config.MustLoad("./config.env")

	client, err := mongo.InitDataLayer()
	if err != nil {
		log.Fatal("failed to init storege", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)

	}

	userRepository := mongo.NewUserRepository(client)

	userUseCase := usecase.NewUserUseCase(userRepository)

	// config, err := helper.LoadConfig(".")
	// if err != nil {
	// 	log.Fatal("cannot load config:", err)
	// }
	var telegramBot tgbot.TelegramBot
	telegramBot.Init(userUseCase)

	ctx, cancel := context.WithCancel(context.Background())
	go telegramBot.ReceiveUpdates(ctx)

	log.Println("Start listening for updates. Press enter to stop")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()

	// // bot, err := tgbotapi.NewBotAPI(config.TG_Bot_Token)
	// // if err != nil {
	// // 	panic(err)
	// // }
	// // bot.Debug = true

	// updateConfig := tgbotapi.NewUpdate(0)

	// updateConfig.Timeout = 60

	// updates := bot.GetUpdatesChan(updateConfig)

	// for update := range updates {

	// 	if update.Message == nil {
	// 		continue
	// 	}
	// 	fmt.Println(update)

	// 	var fileConfig tgbotapi.FileConfig

	// 	fileConfig.FileID = update.Message.Document.FileID
	// 	file, err := bot.GetFile(fileConfig)
	// 	if err != nil {
	// 		return
	// 	}
	// 	fmt.Println("@@@@@@@@@@@@@")
	// 	url := file.Link(config.TG_Bot_Token)
	// 	fmt.Println(url)
	// 	fmt.Println("@@@@@@@@@@@@@")

	// 	_, err = YandexDisk.UploadFileLink(url, "api3")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
}
