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

func (c *PkeCommand) IsDmCapable() bool {
	return false
}

func (c *PkeCommand) Run(ctx ken.Context) (err error) {
	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Why does that person have the [APP] tag?",
					Color:       0x8c1bb1,
					Description: "This server is bridged to decentralised chat platform Matrix using a bot. Due to Discord limitations, Matrix users show up with the [APP] tag on the Discord side.",
				},
			},
		},
	})
	return
}
