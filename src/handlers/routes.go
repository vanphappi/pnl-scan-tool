package handlers

import (
	"context"
	"os"
	"os/signal"
	"pnl-scan-tool/package/configs"
	"pnl-scan-tool/src/services"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/gofiber/fiber/v2"
)

func WalletTrackerRoutes(app *fiber.App, taskManager *services.WalletTrackerTaskManager) {

	var env, err = configs.LoadConfig(".")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(env.TELEGRAM_BOT_TOKEN, opts...)
	if err != nil {
		panic(err)
	}

	app.Post("api/wallettracker/add", taskManager.AddWalletTrackerHandler(ctx, b))
	app.Delete("api/wallettracker/delete", taskManager.CancelTaskHandler)
	app.Get("api/wallettracker/list", taskManager.ListTasksHandler)
	app.Post("api/wallettracker/shutdown", taskManager.ShutdownHandler)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
