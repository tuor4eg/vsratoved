package main

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/tuor4eg/vsratoved/internal/bot"
	"github.com/tuor4eg/vsratoved/internal/config"
	"github.com/tuor4eg/vsratoved/internal/llm"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	token := config.C.TelegramBotToken
	if token == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN is not set")
	}

	telegramBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := telegramBot.GetUpdatesChan(u)

	ctx := context.Background()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		command := update.Message.Command()
		chatID := update.Message.Chat.ID
		requestTime := time.Now().Format("2006-01-02 15:04:05")

		log.Printf("chat_id=%d, text=/%s, time=%s", chatID, command, requestTime)

		var messageText string
		var err error
		var advice *llm.AdviceResponse
		switch command {
		case "start":
			messageText = bot.StartMessage()

		case "vsrata":
			advice, err = llm.GetWeirdAdvice(ctx, "clean")
			if err != nil {
				messageText = bot.ErrorMessage()
			} else {
				messageText = bot.AdviceMessage(advice)
			}

		case "vsrata_spicy":
			advice, err = llm.GetWeirdAdvice(ctx, "spicy")
			if err != nil {
				messageText = bot.ErrorMessage()
			} else {
				messageText = bot.AdviceMessage(advice)
			}

		default:
			messageText = bot.UnknownCommandMessage()
		}

		if err != nil {
			log.Printf("error: chat_id=%d, command=/%s, err=%v", chatID, command, err)
		}

		if messageText != "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
			telegramBot.Send(msg)
		}
	}
}
