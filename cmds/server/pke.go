package cmdsServer

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type PkeCommand struct {
	ken.SlashCommand
	ken.DmCapable
}

var (
	_ ken.SlashCommand = (*PkeCommand)(nil)
	_ ken.DmCapable    = (*PkeCommand)(nil)
)

func (c *PkeCommand) Name() string {
	return "botexplain"
}

func (c *PkeCommand) Description() string {
	return "an explainer on bot tags."
}

func (c *PkeCommand) Version() string {
	return "1.0.0"
}

func (c *PkeCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *PkeCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}
