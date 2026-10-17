# bot

A plugin-based Discord bot written in Go. Slash commands are provided by
plugins compiled as Go `plugin` shared libraries (`.so`) and loaded at
runtime.

## How it works

- `cmd/server` — entry point. Parses config, connects to Discord, loads
  plugins from the plugin directory, registers their handlers and slash
  commands, then blocks until `SIGINT`/`SIGTERM`.
- `bot/` — wraps a `discordgo` session; registers plugin handlers and
  application commands.
- `plugin/` — plugin interface and loading. Plugins are `.so` files built
  with `-buildmode=plugin`; each must export a `Register() (Plugin, error)`
  function. Plugin directories are scanned recursively.
- `builtins/` — plugins that ship with the repo (e.g. `echo`).
- `config/` — configuration from environment variables.

## Configuration

Environment variables:

| Variable         | Required | Default      | Description                                      |
| ---------------- | -------- | ------------ | ------------------------------------------------ |
| `BOT_TOKEN`      | yes      | —            | Discord token, prefixed with `Bot ` or `Bearer ` |
| `BOT_GUILD_ID`   | no       | —            | Guild to register commands in; empty = all guilds |
| `BOT_PLUGIN_DIR` | no       | `./plugins`  | Directory to load `.so` plugins from             |

## Building

```sh
make          # builds bin/server and plugins/*.so from builtins/*
make clean    # removes built binaries and plugins
```

Each subdirectory of `builtins/` is built into a matching `.so` in
`plugins/`.

## Running

```sh
export BOT_TOKEN="Bot <token>"
make run      # builds everything, then starts bin/server
```

The server shuts down gracefully on `SIGINT`/`SIGTERM`.

## Writing a plugin

Create a directory (e.g. `builtins/myplugin/`) with a `main` package:

```go
package main

import (
	"bot/plugin"

	"github.com/bwmarrin/discordgo"
)

type MyPlugin struct{}

func (p *MyPlugin) Metadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name:    "myplugin",
		Version: "0.1.0",
		Author:  "you",
	}
}

func (p *MyPlugin) Commands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{Name: "hello", Description: "Says hello."},
	}
}

func (p *MyPlugin) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// handle interactions
}

func Register() (plugin.Plugin, error) {
	return &MyPlugin{}, nil
}

func main() {}
```

Build it with `go build -buildmode=plugin -o plugins/myplugin.so ./builtins/myplugin`
(or just run `make`) and place the `.so` under the plugin directory.

Note: Go plugins must be built from the same source tree and toolchain as
the server binary.

## License

Copyright © 2026 dowerx

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/agpl-3.0.html>.
