package cmdsServer

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/zekrotja/ken"
)

type RandomColor struct {
	ken.SlashCommand
	ken.DmCapable
}

var (
	_ ken.SlashCommand = (*RandomColor)(nil)
	_ ken.DmCapable    = (*RandomColor)(nil)
)

func (c *RandomColor) Name() string {
	return "randomcolor"
}

func (c *RandomColor) Description() string {
	return "generate a random color!"
}

func (c *RandomColor) Version() string {
	return "1.0.0"
}

func (c *RandomColor) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *RandomColor) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *RandomColor) IsDmCapable() bool {
	return true
}

func (c *RandomColor) Run(ctx ken.Context) (err error) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var randomcolorRGB colorful.Color
	randomcolorRGB.R = rand.Float64()
	randomcolorRGB.G = rand.Float64()
	randomcolorRGB.B = rand.Float64()
	var randomRedReadable string = fmt.Sprintf("%.0f", randomcolorRGB.R*255)
	var randomGreenReadable string = fmt.Sprintf("%.0f", randomcolorRGB.G*255)
	var randomBlueReadable string = fmt.Sprintf("%.0f", randomcolorRGB.B*255)

	// HEX
	var randomColorHexTruncated string = strings.ReplaceAll(randomcolorRGB.Hex(), "#", "")
	var RandomColorHexInt64, res = strconv.ParseInt(randomColorHexTruncated, 16, 32)

	// HSV
	var randomHue, randomSat, randomVal = randomcolorRGB.Hsv()
	var randomHueReadable string = fmt.Sprintf("%.0f", randomHue)
	var randomSatReadable string = fmt.Sprintf("%.0f", randomSat*100)
	var randomValReadable string = fmt.Sprintf("%.0f", randomVal*100)

	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Here's a random color!",
					Color: int(RandomColorHexInt64),
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL:    "https://singlecolorimage.com/get/" + randomColorHexTruncated + "/100x100",
						Width:  100,
						Height: 100,
					},
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Hex",
							Value: randomcolorRGB.Hex(),
						},
						{
							Name:  "RGB",
							Value: "[" + randomRedReadable + ", " + randomGreenReadable + ", " + randomBlueReadable + "]",
						},
						{
							Name:  "HSV",
							Value: randomHueReadable + "°, " + randomSatReadable + "%, " + randomValReadable + "%",
						},
					},
				},
			},
		},
	})

	if err != nil {
		return
	}
	if res != nil {
		return
	}

	return
}
