package cmds

import (
	"github.com/zekrotja/ken"
)

// type PingCommand struct{}

var (
	_ ken.SlashCommand = (*TestCommand)(nil)
	_ ken.DmCapable    = (*TestCommand)(nil)
)

func (c *TestCommand) Name() string {
	return "ping"
}

func (c *TestCommand) Description() string {
	return "Pong!"
}

func (c *TestCommand) Version() string {
	return "1.0.0"
}

func (c *TestCommand) IsDmCapable() bool {
	return true
}

func (c *TestCommand) Run(ctx ken.Context) (err error) {

}
