package telegram

import (
	"awesomeProject/internal/myLogger"
	"flag"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
)

var webhook *string

var bot *tgbotapi.BotAPI

const GetYourCurrencyRate string = "5655548481:AAFlSdYiyhgf7VX1jX5k5h2WsA6Fu-RIUZI"

func init() {
	webhook = flag.String("webhook", "", "Set webhook link if necessary")
	flag.Parse()

	var err error
	bot, err = tgbotapi.NewBotAPI(GetYourCurrencyRate)
	if err != nil {
		zap.L().Error("tgbotapi.NewBotAPI", zap.Error(err))
		return
	}
	bot.Debug = false
	zap.L().Info("Authorized on account", zap.String("bot.Self.UserName", bot.Self.UserName))
}

func newBot() tgbotapi.UpdatesChannel {
	if *webhook != "" {
		wh, err := tgbotapi.NewWebhook(*webhook)
		if err != nil {
			zap.L().Error("", zap.Error(err))
			return nil
		}

		if _, err = bot.Request(wh); err != nil {
			zap.L().Error("", zap.Error(err))
			return nil
		}
		info, err := bot.GetWebhookInfo()
		if err != nil {
			zap.L().Error("", zap.Error(err))
			return nil
		}

		if info.LastErrorDate != 0 {
			myLogger.Logger.Printf("Telegram callback failed: %s", info.LastErrorMessage)
		}

		updates := bot.ListenForWebhook("/")
		go func() {
			if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
				myLogger.Logger.Fatal(err)
			}
		}()
		return updates
	} else {
		wh := tgbotapi.DeleteWebhookConfig{DropPendingUpdates: true}

		if _, err := bot.Request(wh); err != nil {
			zap.L().Error("", zap.Error(err))
			return nil
		}
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates := bot.GetUpdatesChan(u)
		return updates
	}

}

func getUpdates() {
	for update := range newBot() {
		update.SentFrom()
		var u Update
		u = Update(update)
		u.defineMessagetype()
	}
}

type Update tgbotapi.Update

func (u *Update) defineMessagetype() {
	if u.Message != nil {
		switch {
		case u.Message.IsCommand():
			u.commandProcess()
		case u.isText():

		case u.CallbackQuery != nil:

		default:

		}
	}

}

func (u *Update) isText() bool {
	if u.Message.Text != "" {
		return true
	}
	return false
}

func (u *Update) commandProcess() {
	switch u.Message.Command() {
	case "start":
		//tgbotapi.NewMessage(,"Hello")
	case "auth":

	default:

	}
}

type Updater interface {
	SentFrom()
}
