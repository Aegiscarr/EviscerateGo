package cmdsServer

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type EightBallCommand struct {
	_ ken.SlashCommand
	_ ken.DmCapable
}

func (c *EightBallCommand) Name() string {
	return "askevi"
}

func (c *EightBallCommand) Description() string {
	return "ask a question and get an answer, magic 8-ball style. i swear this thing isnt sentient. please help."
}

func (c *EightBallCommand) Version() string {
	return "1.0.0"
}

func (c *EightBallCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *EightBallCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "question",
			Description: "question to ask",
			Required:    true,
		},
	}
}

func (c *EightBallCommand) IsDmCapable() bool {
	return true
}

func (c *EightBallCommand) Run(ctx ken.Context) (err error) {

	rand.Seed(time.Now().UnixNano())

	var (
		color      int
		query      string
		responseID int
	)

	query = ctx.Options().GetByName("question").StringValue()

	responses := [20]string{"It is certain.", "It is decidedly so.", "Without a doubt.", "Yes, definitely.", "You may rely on it.", "As I see it, yes.", "Most likely.", "Outlook good.", "Yes.", "Signs point to yes", "Reply hazy, try again.", "Ask again later.", "Better not tell you now.", "Cannot predict now.", "Concentrate and ask again.", "Don't count on it.", "My reply is no.", "My sources say no.", "Outlook not so good.", "Very doubtful."}

	responseID = rand.Intn(19)
	fmt.Println(query)

	if responseID < 10 {
		color = 0x00ae00 //affirmative
	} else if responseID > 9 && responseID < 14 {
		color = 0xfcb103 //non-commital
	} else {
		color = 0xff0000 //negative
	}

	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       responses[responseID],
					Description: `Question: "` + query + `"`,
					Color:       color,
					Footer: &discordgo.MessageEmbedFooter{
						Text: "Disclaimer: This answer is based on absolutely nothing. At least, Aegis thinks so.",
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
