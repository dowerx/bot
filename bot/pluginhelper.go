package bot

import (
	"bot/plugin"
)

// This implementation is passed to plugins to expose helper functions.
type PluginHelpers struct {
	bot *Bot
}

func (h *PluginHelpers) Plugins() map[string]*plugin.Plugin {
	plugins := make(map[string]*plugin.Plugin)
	for k, p := range h.bot.plugins {
		plugins[k] = &p.plugin
	}
	return plugins
}
