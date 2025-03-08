package cmdsServer

import (
	"EviscerateGo/lib/api"
	"EviscerateGo/lib/structs"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type UnsplashCommand struct {
	_ ken.SlashCommand
	_ ken.DmCapable
}

var (
	_ ken.SlashCommand = (*UnsplashCommand)(nil)
	_ ken.DmCapable    = (*UnsplashCommand)(nil)
)

func (c *UnsplashCommand) Name() string {
	return "unsplash"
}

func (c *UnsplashCommand) Description() string {
	return "find a random image on unsplash"
}

func (c *UnsplashCommand) Version() string {
	return "1.0.0"
}

func (c *UnsplashCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *UnsplashCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "query",
			Description: "search query for unsplash",
			Required:    true,
		},
	}
}

func (c *UnsplashCommand) IsDmCapable() bool {
	return true
}

func (c *UnsplashCommand) Run(ctx ken.Context) (err error) {
	var (
		query    string
		unsplash *structs.UnsplashRandom
	)

	query = ctx.Options().GetByName("query").StringValue()

	unsplash = api.UnsplashImageFromApi(query)

	if unsplash.URLs.Small != "" {

		fmt.Println(unsplash.URLs.Small)
		fmt.Println(unsplash.Links.DownloadLocation)

		err = ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: "A random photo matching '" + query + "' (click to view on Unsplash)",
						URL:   fmt.Sprintf("%s?utm_source=Eviscerate-The-Synth&utm_medium=referral", unsplash.Links.HTML),
						Image: &discordgo.MessageEmbedImage{
							URL:    unsplash.URLs.Small,
							Width:  512,
							Height: 512,
						},
						Description: fmt.Sprintf("ALL IMAGES ARE TAKEN FROM UNSPLASH.COM\n\n%s", unsplash.Description),
						Author: &discordgo.MessageEmbedAuthor{
							Name: fmt.Sprintf("Photo by %s", unsplash.User.Name),
							URL:  unsplash.User.Links.HTML,
						},
					},
				},
			},
		})
		if err != nil {
			return
		}
	} else {
		err := ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Something happened, I couldn't find an image matching '" + query + "'! (Small Image URL is Blank)",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			return err
		}
	}

	return err
}
