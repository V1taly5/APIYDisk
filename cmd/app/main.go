package main

import (
	"bufio"
	"context"
	"log"
	"os"

	"github.com/V1taly5/APIYDisk/internal/entity"
	"github.com/V1taly5/APIYDisk/internal/infrastructure/repository"
	"github.com/V1taly5/APIYDisk/internal/tgbot"
)

func main() {
	// cfg := config.MustLoad("./config.env")

	client, err := repository.InitDataLayer()
	if err != nil {
		log.Fatal("failed to init storege", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)

	}

	userRepository := repository.NewUserRepository(client)

	// config, err := helper.LoadConfig(".")
	// if err != nil {
	// 	log.Fatal("cannot load config:", err)
	// }
	var telegramBot tgbot.TelegramBot
	telegramBot.Init(userRepository)
	ctx, cancel := context.WithCancel(context.Background())
	go telegramBot.ReceiveUpdates(ctx)

	// var User entity.User
	var User3 entity.User
	// User.ChatID = 34242342
	// User.YandexDiskToken = "y0_AgAAAAAWC_PqAADLWwAAAADlhx9_puyF720NQDaD8to781KL6yfqMhA"
	// telegramBot.Repo.Create(&User)
	User3, err = telegramBot.Repo.FindByChatID(34242342)
	User3.State = ""
	// fmt.Println(User3)
	telegramBot.Repo.Update(&User3)

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
