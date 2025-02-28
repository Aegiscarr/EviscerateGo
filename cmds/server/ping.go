package cmdsServer

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type PingCommand struct {
	ken.SlashCommand
}

var (
	_ ken.SlashCommand = (*PingCommand)(nil)
	_ ken.DmCapable    = (*PingCommand)(nil)
)

func (c *PingCommand) Name() string {
	return "ping"
}

func (c *PingCommand) Description() string {
	return "Pong!"
}

func (c *PingCommand) Version() string {
	return "1.0.0"
}

func (c *PingCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *PingCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *PingCommand) IsDmCapable() bool {
	return true
}

func (c *PingCommand) Run(ctx ken.Context) (err error) {
	s := ctx.GetSession()

	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong! Latency: " + s.HeartbeatLatency().String(),
		},
	})

	if err != nil {
		return
	}
	return err
}
