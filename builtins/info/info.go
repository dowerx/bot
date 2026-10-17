// Prints some info about the bot.

package main

import (
	"bot/plugin"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
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

// Register implements [plugin.Plugin].
func (p *InfoPlugin) Register(helpers plugin.PluginHelpers) []*discordgo.ApplicationCommand {
	p.helpers = helpers
	return []*discordgo.ApplicationCommand{
		{
			Name:        "info",
			Description: "Prints some info about the bot.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "area",
					Description: "Area you are interested in.",
					Required:    true,
				},
			},
		},
	}
}

func (p *InfoPlugin) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	if data.Name != "info" {
		return
	}

	var message string
	area := data.Options[0].StringValue()
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

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	}); err != nil {
		slog.Error("failed to send echo reply", "error", err.Error())
	}
}

func main() {}
