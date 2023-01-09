package main

import (
	"github.com/Bukhashov/filechain/internal/config"
	"github.com/Bukhashov/filechain/pkg/logging"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	logger := logging.GetLogger()
	cfg := config.GetConfig()
	
	bot, err := tgbotapi.NewBotAPI(cfg.Token.Telegram); if err != nil {
		logger.Info(err)
		return
	}
	// bot.Debug = true
	logger.Info("Telegram bot api start")

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		// ignore any non-Message updates
		if update.Message == nil {
			continue
		}
		// ignore any non-command Messages
		if !update.Message.IsCommand() {
            continue
        }

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
			case "start" :
				msg.Text = `
					Сәлем қолданушы FileChain жобамызға қош келіпсіз.
					/singin командасы арқылы жеке Папакаларынызға кізініз.
					Егер сіз біздін жобады жаңадан қоладушы болсаныз /singup командасы арқылы тіркеле аласыз.
					Егерде түсінбей жатқа заттарыныз болса /help коммандасы арқылы ақпарат алсаныз болады
				`

				if _, err := bot.Send(msg); err != nil {
					logger.Info(err)
				}

			case "singin" :
				msg.Text = `
					әсімізді [username] және почта [email] ды жіберіңіз
					Ескерту: есімініз ағылшын әріптеріме [a-z A-Z] жазыныз
				`
				if _, err := bot.Send(msg); err != nil {
					logger.Info(err)
				}

				

			case "help" :
				msg.Text = `Text /help`
		}

		

	}

	
}