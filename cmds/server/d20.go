package cmdsServer

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type D20Command struct {
	ken.SlashCommand
}

var (
	_ ken.DmCapable    = (*D20Command)(nil)
	_ ken.SlashCommand = (*D20Command)(nil)
)

func (c *D20Command) Name() string {
	return "d20"
}

func (c *D20Command) Description() string {
	return "decide your fate at the hands of the iconic 20-sided dice!"
}

func (c *D20Command) Version() string {
	return "1.0.0"
}

func (c *D20Command) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *D20Command) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "rollfor",
			Description: "the thing youre rolling for",
		},
	}
}

func (c *D20Command) IsDmCapable() bool {
	return true
}

func (c *D20Command) Run(ctx ken.Context) (err error) {

	var rolledFor string
	if v, ok := ctx.Options().GetByNameOptional("rollfor"); ok {
		rolledFor = v.StringValue()
	}

	dicerolls := [20]string{`Luck really isn't on your side today, huh? It's a 1.`, `A 2. Couldn't have been much worse.`, `It's a tree!- Oh, wait. A 3.`, `A four-se. Of course. That didn't work, did it?`, `A 5. Nothing funny here.`, `A 6. The devil, anyone?`, `Lucky number 7! Now can you get two more?`, `8, not bad.`, `Just under halfway up. A 9`, `A 10! Halfway up the scale!`, `11. Decent.`, `12. Could have been much worse. Could've also been better, though.`, `13. Feelin' lucky?`, `Aand it's come up 14!`, `15! Getting up there!`, `16, solid.`, `17. Rolling real high now, aren't you?`, `18! You're old eno- wait this isn't a birthday.`, `19! So CLOSE!`, `NAT 20 BAYBEE!`}

	if rolledFor == "" {
		rolledFor = "no thing"
	}

	u := ctx.GetEvent().User.Username

	rand.Seed(time.Now().UnixNano())

	var diceroll int = rand.Intn(20)

	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: dicerolls[diceroll],
					Color: 0x8c1bb1,
					Footer: &discordgo.MessageEmbedFooter{
						Text: u + " rolled for " + rolledFor,
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
