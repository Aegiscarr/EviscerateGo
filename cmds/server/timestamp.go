package cmdsServer

import (
	"EviscerateGo/lib/time"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type TimestampCommand struct {
	_ ken.SlashCommand
	_ ken.DmCapable
}

var (
	_ ken.SlashCommand = (*TimestampCommand)(nil)
	_ ken.DmCapable    = (*TimestampCommand)(nil)
)

func (c *TimestampCommand) Name() string {
	return "timestamp"
}

func (c *TimestampCommand) Description() string {
	return "generate a discord timestamp"
}

func (c *TimestampCommand) Version() string {
	return "1.0.0"
}

func (c *TimestampCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *TimestampCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "time",
			Description: "time to use // format as 01/02 03:04:05PM '06 -0700 (america style so mm/dd",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "display",
			Description: "how the timestamp gets displayed",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "short time // hh:mm",
					Value: "t",
				},
				{
					Name:  "long time // hh:mm:ss",
					Value: "T",
				},
				{
					Name:  "short date // dd-mm-yyyy",
					Value: "d",
				},
				{
					Name:  "long date // dd [name of month] yyyy",
					Value: "D",
				},
				{
					Name:  "long date with short time // dd [name of month] yyyy at hh:mm",
					Value: "f",
				},
				{
					Name:  "long date with day of week and short time // [name of day] dd [name of month] yyyy at hh:mm",
					Value: "F",
				},
				{
					Name:  "relative // (in) x hours/minutes/seconds (ago)",
					Value: "R",
				},
			},
		},
	}
}

func (c *TimestampCommand) IsDmCapable() bool {
	return true
}

func (c *TimestampCommand) Run(ctx ken.Context) (err error) {
	var (
		input   string
		display string
	)

	input = ctx.Options().GetByName("time").StringValue()
	display = ctx.Options().GetByName("display").StringValue()

	t, err := time.Str2Utc(input)
	if err != nil {
		fmt.Printf("An error occurred while converting input to UTC time: %v", err)
		err := ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("An error occurred while converting input to UTC time: %v", err),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			return err
		}
	}
	//snowflake := Utc2snowflake(t)
	//fmt.Println(snowflake)

	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Here is a snowflake timestamp!",
					Color: 0x8c1bb1,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Display",
							Value: fmt.Sprintf("<t:%v:%v>", t, display),
						},
						{
							Name:  "Raw",
							Value: fmt.Sprintf("`<t:%v:%v>`", t, display),
						},
					},
				},
			},
		},
	})
	if err != nil {
		return
	}

	return err
}
