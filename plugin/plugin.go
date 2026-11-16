package plugin

import (
	"errors"
	"io/fs"
	"log/slog"
	"path/filepath"
	"plugin"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

const (
	PLUGIN_ENTRY_FUNC = "Entry"
)

type PluginMetadata struct {
	Name    string
	Version string
	Author  string
}

// The bot passes an implementation of this interface to every plugin.
// They access information and helper functions through this.
// Handlers can access the Discord client directly via event.Client().
type PluginHelpers interface {
	Plugins() map[string]*Plugin
}

// Plugins implement this interface.
// Every plugin must have an Entry function that returns a struct implementing this interface.
type Plugin interface {
	Metadata() PluginMetadata
	Commands(helpers PluginHelpers) []SlashCommand
}

// SlashCommand couples a slash command definition with its handler.
// The route path is derived from Definition.Name. Subcommands are part of the
// definition's option tree; the top-level route also receives subcommand
// interactions, so a single handler can dispatch them using the interaction data.
type SlashCommand struct {
	Definition discord.SlashCommandCreate
	Handler    handler.SlashCommandHandler
}

// Load a plugin from a shared library file.
func Load(path string) (Plugin, error) {
	slog.Info("loading plugin",
		"path", path,
	)

	lib, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	entrySym, err := lib.Lookup(PLUGIN_ENTRY_FUNC)
	if err != nil {
		return nil, err
	}

	entry, ok := entrySym.(func() (Plugin, error))
	if !ok {
		return nil, errors.New("failed to cast plugin entry")
	}

	p, err := entry()
	if err != nil {
		return nil, err
	}

	meta := p.Metadata()
	slog.Info("loaded plugin",
		"path", path,
		"name", meta.Name,
		"version", meta.Version,
		"author", meta.Author,
	)

	return p, nil
}

// Load all plugins from a directory. Recursive.
// Returns at the first error.
func LoadDir(path string) ([]Plugin, error) {
	var plugins []Plugin

	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(d.Name()) != ".so" {
			return nil
		}

		plug, err := Load(p)
		if err != nil {
			return err
		}

		plugins = append(plugins, plug)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return plugins, nil
}
