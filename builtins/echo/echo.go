// Returns the same message as passed in.

package main

import (
	"bot/plugin"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
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

func (p *EchoPlugin) Commands(_ plugin.PluginHelpers) []plugin.SlashCommand {
	return []plugin.SlashCommand{
		{
			Definition: discord.SlashCommandCreate{
				Name:        "echo",
				Description: "Repeates a message.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:        "message",
						Description: "Message to repeat.",
						Required:    true,
					},
				},
			},
			Handler: func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
				message := data.String("message")

				if err := e.CreateMessage(discord.MessageCreate{
					Content: message,
				}); err != nil {
					return fmt.Errorf("failed to send echo reply: %w", err)
				}
				return nil
			},
		},
	}
}

func main() {}
