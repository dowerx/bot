package bot

import (
	"bot/plugin"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

type botPlugin struct {
	plugin plugin.Plugin
}

type pendingPlugin struct {
	meta     plugin.PluginMetadata
	plugin   plugin.Plugin
	router   *handler.Mux
	commands []plugin.SlashCommand
}

// The bots internal data struct.
type Bot struct {
	client        *bot.Client
	guildIDs      []snowflake.ID
	plugins       map[string]botPlugin
	router        *handler.Mux
	commands      []discord.ApplicationCommandCreate
	commandOwners map[string]string
}

// Init creates the Discord client and interaction router.
// Call Start to open the gateway after plugins have been registered.
func (b *Bot) Init(token string, guildID string) error {
	var guildIDs []snowflake.ID
	if guildID != "" {
		id, err := snowflake.Parse(guildID)
		if err != nil {
			return fmt.Errorf("failed to parse guild ID: %w", err)
		}
		guildIDs = []snowflake.ID{id}
	}

	b.router = handler.New()
	b.guildIDs = guildIDs
	b.plugins = make(map[string]botPlugin)
	b.commands = nil
	b.commandOwners = make(map[string]string)

	client, err := disgo.New(token,
		bot.WithEventListeners(b.router),
		bot.WithGatewayConfigOpts(
			// Interaction events do not require any gateway intents.
			gateway.WithIntents(gateway.IntentsNone),
		),
	)
	if err != nil {
		return err
	}
	b.client = client

	return nil
}

// Start opens the gateway connection. Call it after plugin routes and command
// definitions have been registered.
func (b *Bot) Start() error {
	if b.client == nil {
		return fmt.Errorf("bot is not initialized")
	}
	return b.client.OpenGateway(context.Background())
}

func (b *Bot) Close() {
	if b.client != nil {
		b.client.Close(context.Background())
	}
}

// Blocks until an interrupt signal is received, then returns.
func (b *Bot) Run() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	sig := <-sigs
	slog.Info("received signal", "signal", sig.String())
	signal.Stop(sigs)
	close(sigs)
}

// Prints some info about the connected session/user.
func (b *Bot) Info() {
	user, ok := b.client.Caches.SelfUser()
	if !ok {
		slog.Warn("session user not available")
		return
	}

	displayName := user.Username
	if user.GlobalName != nil && *user.GlobalName != "" {
		displayName = *user.GlobalName
	}

	slog.Info("session info",
		"username", user.Username,
		"displayName", displayName,
	)
}

// RegisterAll validates all plugins, synchronizes their complete command set
// with Discord once, and only then mounts their routes and records them.
// If any plugin is invalid or synchronization fails, the bot's state is left
// unchanged.
func (b *Bot) RegisterAll(plugins []plugin.Plugin) error {
	if b.client == nil {
		return fmt.Errorf("bot is not initialized")
	}

	pending := make([]pendingPlugin, 0, len(plugins))

	pluginNames := make(map[string]struct{}, len(b.plugins)+len(plugins))
	for name := range b.plugins {
		pluginNames[name] = struct{}{}
	}

	allCommands := make([]discord.ApplicationCommandCreate, 0, len(b.commands))
	allCommands = append(allCommands, b.commands...)

	commandOwners := make(map[string]string, len(b.commandOwners))
	for name, owner := range b.commandOwners {
		commandOwners[name] = owner
	}

	for _, p := range plugins {
		meta := p.Metadata()
		if meta.Name == "" {
			return fmt.Errorf("plugin has an empty name")
		}
		if _, found := pluginNames[meta.Name]; found {
			return fmt.Errorf("plugin name already taken: %s", meta.Name)
		}
		pluginNames[meta.Name] = struct{}{}

		pluginRouter := handler.New()
		slashCommands := p.Commands(&PluginHelpers{bot: b})
		for _, command := range slashCommands {
			name := command.Definition.Name
			if name == "" {
				return fmt.Errorf("plugin %s returned a slash command without a name", meta.Name)
			}
			if command.Handler == nil {
				return fmt.Errorf("plugin %s returned no handler for command %s", meta.Name, name)
			}

			if owner, found := commandOwners[name]; found {
				if owner == meta.Name {
					return fmt.Errorf("plugin %s returned duplicate command prefix: /%s", meta.Name, name)
				}
				return fmt.Errorf("command prefix already taken: /%s (plugins %s and %s)", name, owner, meta.Name)
			}
			commandOwners[name] = meta.Name

			pluginRouter.SlashCommand("/"+name, command.Handler)
			allCommands = append(allCommands, command.Definition)
		}

		pending = append(pending, pendingPlugin{
			meta:     meta,
			plugin:   p,
			router:   pluginRouter,
			commands: slashCommands,
		})
	}

	// Syncing replaces the complete command set, so do it once for all plugins.
	if err := handler.SyncCommands(b.client, allCommands, b.guildIDs); err != nil {
		return err
	}

	for _, entry := range pending {
		b.router.Mount("", entry.router)
		b.plugins[entry.meta.Name] = botPlugin{
			plugin: entry.plugin,
		}
		slog.Info("registered handler",
			"plugin", entry.meta.Name,
		)

		for _, command := range entry.commands {
			slog.Info("registered command",
				"plugin", entry.meta.Name,
				"command", command.Definition.Name,
				"description", command.Definition.Description,
			)
		}
	}

	b.commands = allCommands
	b.commandOwners = commandOwners

	return nil
}
