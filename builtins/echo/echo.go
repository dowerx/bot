// Returns the same message as passed in.

package main

import (
	"bot/plugin"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type EchoPlugin struct {
}

func (p *EchoPlugin) Metadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name:    "echo",
		Version: "0.1.0",
		Author:  "dowerx",
	}
}

func Entry() (plugin.Plugin, error) {
	return &EchoPlugin{}, nil
}

func (p *EchoPlugin) Register(plugin.PluginHelpers) []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "echo",
			Description: "Repeates a message.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Message to repeat.",
					Required:    true,
				},
			},
		},
	}
}

func (p *EchoPlugin) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	if data.Name != "echo" {
		return
	}

	message := data.Options[0].StringValue()

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
