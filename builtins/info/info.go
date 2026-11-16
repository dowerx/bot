// Prints some info about the bot.

package main

import (
	"bot/plugin"
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

type InfoPlugin struct {
	helpers plugin.PluginHelpers
}

func (p *InfoPlugin) Metadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name:    "info",
		Version: "0.1.0",
		Author:  "dowerx",
	}
}

func Entry() (plugin.Plugin, error) {
	return &InfoPlugin{}, nil
}

// Commands implements [plugin.Plugin].
func (p *InfoPlugin) Commands(helpers plugin.PluginHelpers) []plugin.SlashCommand {
	p.helpers = helpers
	return []plugin.SlashCommand{
		{
			Definition: discord.SlashCommandCreate{
				Name:        "info",
				Description: "Prints some info about the bot.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:        "area",
						Description: "Area you are interested in.",
						Required:    true,
					},
				},
			},
			Handler: p.info,
		},
	}
}

func (p *InfoPlugin) info(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	var message string
	area := data.String("area")
	switch area {
	case "plugins":
		plugins := p.helpers.Plugins()
		pluginNames := make([]string, len(plugins))
		i := 0
		for name := range plugins {
			pluginNames[i] = name
			i++
		}
		message = strings.Join(pluginNames, ", ")
	default:
		message = "unknown command"
	}

	if err := e.CreateMessage(discord.MessageCreate{
		Content: message,
	}); err != nil {
		return fmt.Errorf("failed to send info reply: %w", err)
	}
	return nil
}

func main() {}
