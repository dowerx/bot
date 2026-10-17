package config

import (
	"os"

	"github.com/go-playground/validator/v10"
)

const (
	ENV_TOKEN          = "BOT_TOKEN"
	ENV_GUILD_ID       = "BOT_GUILD_ID"
	ENV_PLUGIN_DIR     = "BOT_PLUGIN_DIR"
	DEFAULT_PLUGIN_DIR = "./plugins"
)

func envOrDefault(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Holds application configuration. Fill it with Parse().
type Config struct {
	Token     string `validate:"required,min=1,printascii"`
	GuildID   string `validate:"max=20,printascii"`
	PluginDir string `validate:"dir"`
}

// Sets the default values and parses environment variables
// into a Config struct. Throws an error if validation fails.
func Parse() (Config, error) {
	cfg := Config{
		Token:     envOrDefault(ENV_TOKEN, ""),    // must be prefixed with 'Bot ' or 'Bearer '
		GuildID:   envOrDefault(ENV_GUILD_ID, ""), // empty guild id means all guilds
		PluginDir: envOrDefault(ENV_PLUGIN_DIR, DEFAULT_PLUGIN_DIR),
	}

	v := validator.New()
	if err := v.Struct(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
