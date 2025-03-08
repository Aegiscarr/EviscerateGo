package cmdsServer

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type UserInfoCommand struct {
	_ ken.SlashCommand
}

var (
	_ ken.SlashCommand = (*UserInfoCommand)(nil)
)

func (c *UserInfoCommand) Name() string {
	return "userinfo"
}

func (c *UserInfoCommand) Description() string {
	return "shows a bunch of info about a user"
}

func (c *UserInfoCommand) Version() string {
	return "1.0.0"
}

func (c *UserInfoCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *UserInfoCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "the user you want to get info of",
			Required:    true,
		},
	}
}

func (c *UserInfoCommand) Run(ctx ken.Context) (err error) {
	var (
		user        discordgo.User
		color       int
		colorstring string
		colortext   string
		avatarURL   string
		bannerURL   string
		uid         string
		uname       string
		discrim     string
		mention     string
		bannervalue string
	)

	user = *ctx.Options().GetByName("user").UserValue(ctx)

	color = user.AccentColor
	avatarURL = user.AvatarURL(strconv.Itoa(1024))
	bannerURL = user.BannerURL(strconv.Itoa(1024))
	uname = user.Username
	discrim = user.Username
	uid = user.ID
	mention = user.Mention()

	if bannerURL == "" {
		bannervalue = uname + " does not have a banner set."
	} else {
		bannervalue = "[Link](" + bannerURL + ")"
	}

	// check and fix leading zeroes

	colorstring = strconv.FormatInt(int64(color), 16)

	if color < 1048576 && color > 65535 {
		colortext = "#0" + colorstring
	} else if color < 65536 && color > 4095 {
		colortext = "#00" + colorstring
	} else if color < 4096 && color > 255 {
		colortext = "#000" + colorstring
	} else if color < 256 && color > 15 {
		colortext = "#0000" + colorstring
	} else if color < 16 && color > 0 {
		colortext = "#00000" + colorstring
	} else if color == 0 {
		colortext = "#000000"
	} else {
		colortext = "#" + colorstring
	}

	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Some info about " + uname,
					Color: color,
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL:    avatarURL,
						Width:  128,
						Height: 128,
					},
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Username",
							Value:  uname,
							Inline: true,
						},
						{
							Name:   "Discriminator",
							Value:  discrim,
							Inline: true,
						},
						{
							Name:   "Mention",
							Value:  mention,
							Inline: true,
						},
						{
							Name:   "User ID",
							Value:  uid,
							Inline: true,
						},
						{
							Name:   "Accent Color",
							Value:  colortext,
							Inline: true,
						},
						{
							Name:   "Avatar",
							Value:  `[Link](` + avatarURL + `)`,
							Inline: true,
						},
						{
							Name:   "Banner",
							Value:  bannervalue,
							Inline: true,
						},
					},
					Image: &discordgo.MessageEmbedImage{
						URL:    bannerURL,
						Width:  512,
						Height: 512,
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
