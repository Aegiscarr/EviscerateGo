package cmdsServer

import (
	"EviscerateGo/auxiliary/api"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type DevExcuse struct {
	ken.SlashCommand
	ken.DmCapable
}

var (
	_ ken.SlashCommand = (*DevExcuse)(nil)
	_ ken.DmCapable    = (*DevExcuse)(nil)
)

func (c *DevExcuse) Name() string {
	return "devexcuse"
}

func (c *DevExcuse) Description() string {
	return "for if you're *really* starved for excuses in your IT job."
}

func (c *DevExcuse) Version() string {
	return "1.0.0"
}

func (c *DevExcuse) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *DevExcuse) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *DevExcuse) IsDmCapable() bool {
	return true
}

func (c *DevExcuse) Run(ctx ken.Context) (err error) {

	var randomexcuse = api.GetDevExcuse()
	var randomexcuseQuoted = `"` + randomexcuse.Data + `"`

	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: randomexcuseQuoted,
		},
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	return
}
