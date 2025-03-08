package cmdsAdmin

import (
	"EviscerateGo/lib/conf"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type StatusSetCommand struct {
	_ ken.SlashCommand
	_ ken.DmCapable
}

var (
	_ ken.SlashCommand = (*StatusSetCommand)(nil)
	_ ken.DmCapable    = (*StatusSetCommand)(nil)
)

func (c *StatusSetCommand) Name() string {
	return "statusset"
}

func (c *StatusSetCommand) Description() string {
	return "set the bot status. bot admin eyes only."
}

func (c *StatusSetCommand) Version() string {
	return "1.0.0"
}

func (c *StatusSetCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *StatusSetCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "type",
			Description: "change between playing/watching/stuff",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Playing",
					Value: 0,
				},
				{
					Name:  "Streaming",
					Value: 1,
				},
				{
					Name:  "Listening",
					Value: 2,
				},
				{
					Name:  "Watching",
					Value: 3,
				},
				{
					Name:  "Competing",
					Value: 5,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "string",
			Description: "second part of status",
			Required:    true,
		},
	}
}

func (c *StatusSetCommand) IsDmCapable() bool {
	return true
}

func (c *StatusSetCommand) Run(ctx ken.Context) (err error) {
	var (
		actType   int
		actString string
	)

	s := ctx.GetSession()

	if ctx.User().ID == conf.AdminId {
		actType = int(ctx.Options().GetByName("type").IntValue())
		actString = ctx.Options().GetByName("string").StringValue()
		s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{{Type: discordgo.ActivityType(actType), Name: actString}},
		})
		err = ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Status Set",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		if err != nil {
			return
		}
	} else {
		err = ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Heh, yeah, I wasn't gonna give anyone control over the *status* of this thing. -aegis",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			return
		}
	}

	return err
}
