package cmdsServer

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type EchoCommand struct {
	_ ken.SlashCommand
}

var (
	_ ken.SlashCommand = (*EchoCommand)(nil)
)

func (c *EchoCommand) Name() string {
	return "echo"
}

func (c *EchoCommand) Description() string {
	return "echo's a message to a different (or current) channel"
}

func (c *EchoCommand) Version() string {
	return "1.0.0"
}

func (c *EchoCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *EchoCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "message",
			Description: "message to echo",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        "channel",
			Description: "channel to send message to",
		},
	}
}

func (c *EchoCommand) Run(ctx ken.Context) (err error) {

	var (
		messagecontent string
		channelobj     discordgo.Channel
	)

	s := ctx.GetSession()

	messagecontent = ctx.Options().GetByName("message").StringValue()

	if v, ok := ctx.Options().GetByNameOptional("channel"); ok {
		channelobj = *v.ChannelValue(ctx)

	}
	if channelobj.ID == "" {
		channelobj.ID = ctx.GetEvent().ChannelID
	}
	fmt.Println(channelobj)

	_, _ = s.ChannelMessageSendComplex(channelobj.ID, &discordgo.MessageSend{
		Content: messagecontent,
	})

	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Echo sent to %s:\n```%s```", channelobj.Mention(), messagecontent),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	return err
}
