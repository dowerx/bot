package bot

import (
	"bot/plugin"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type botPlugin struct {
	plugin        plugin.Plugin
	removeHandler func()
}

// The bots internal data struct.
type Bot struct {
	session *discordgo.Session
	guildID string
	plugins map[string]botPlugin
}

func (b *Bot) Connect(token string, guildID string) error {
	session, err := discordgo.New(token)
	if err != nil {
		return err
	}
	b.session = session
	b.guildID = guildID
	b.plugins = make(map[string]botPlugin)

	if err := session.Open(); err != nil {
		return err
	}

	return nil
}

func (b *Bot) Close() {
	b.session.Close()
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
	slog.Info("session info",
		"username", b.session.State.User.Username,
		"displayName", b.session.State.User.DisplayName(),
	)
}

// Register a plugins handler and commands.
func (b *Bot) Register(p plugin.Plugin) error {
	meta := p.Metadata()

	_, found := b.plugins[meta.Name]
	if found {
		return fmt.Errorf("plugin name already taken: %s", meta.Name)
	}

	removeHandler := b.session.AddHandler(p.Handle)
	slog.Info("registered handler",
		"plugin", meta.Name,
	)

	for _, command := range p.Register(&PluginHelpers{bot: b}) {
		if _, err := b.session.ApplicationCommandCreate(
			b.session.State.User.ID,
			b.guildID,
			command,
		); err != nil {
			removeHandler()
			return err
		}

		slog.Info("registered command",
			"plugin", meta.Name,
			"command", command.Name,
			"description", command.Description,
		)
	}

	b.plugins[meta.Name] = botPlugin{
		plugin:        p,
		removeHandler: removeHandler,
	}

	return nil
}

// Load a list of plugins. Return at the first error.
func (b *Bot) RegisterAll(plugins []plugin.Plugin) error {
	for _, p := range plugins {
		if err := b.Register(p); err != nil {
			return err
		}
	}
	return nil
}
