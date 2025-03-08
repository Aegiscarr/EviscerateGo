package cmdsServer

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type GoogleCommand struct {
	_ ken.SlashCommand
}

var (
	_ ken.SlashCommand = (*GoogleCommand)(nil)
)

func (c *GoogleCommand) Name() string {
	return "google"
}

func (c *GoogleCommand) Description() string {
	return "googles something for you. definitely."
}

func (c *GoogleCommand) Version() string {
	return "1.0.0"
}

func (c *GoogleCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *GoogleCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "searchquery",
			Description: "the thing you want to oh so snarkily google.",
			Required:    true,
		},
	}
}

func (c *GoogleCommand) Run(ctx ken.Context) (err error) {

	//link = "https://lmgtfy.app/?q=" + strings.ReplaceAll(ctx.Options().GetByName("searchquery").StringValue(), " ", "+")

	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Here you go: https://google.com/?q=" + strings.ReplaceAll(ctx.Options().GetByName("searchquery").StringValue(), " ", "+"),
		},
	})

	if err != nil {
		return
	}
	return err
}
