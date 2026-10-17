package main

import (
	"bot/bot"
	"bot/config"
	"bot/plugin"
	"log/slog"
	"os"
)

func main() {
	os.Exit(run())
}

func run() int {
	// Setup structural JSON logging.
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// Load config.
	cfg, err := config.Parse()
	if err != nil {
		slog.Error("failed to parse config",
			"error", err.Error(),
		)
		return 1
	}

	slog.Info("connecting to Discord", "guildID", cfg.GuildID)
	bot := bot.Bot{}
	if err := bot.Connect(cfg.Token, cfg.GuildID); err != nil {
		slog.Error("failed to connect to Discord",
			"error", err.Error(),
		)
		return 1
	}
	defer bot.Close()

	bot.Info()

	plugins, err := plugin.LoadDir(cfg.PluginDir)
	if err != nil {
		slog.Error("failed to load plugins",
			"error", err.Error(),
		)
		return 1
	}

	if err := bot.RegisterAll(plugins); err != nil {
		slog.Error("failed while loading plugins",
			"error", err.Error(),
		)
		return 1
	}

	bot.Run()

	return 0
}
