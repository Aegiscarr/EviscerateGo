package cmdsServer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/zekrotja/ken"
)

type ColorCommand struct {
	ken.SlashCommand
	ken.DmCapable
}

var (
	_ ken.SlashCommand = (*ColorCommand)(nil)
	_ ken.DmCapable    = (*ColorCommand)(nil)
	_ ken.UserCommand
)

func (c *ColorCommand) Name() string {
	return "color"
}

func (c *ColorCommand) Description() string {
	return "convert a color of your choice to rgb, hex, and hsv, with a nice embed to boot"
}

func (c *ColorCommand) Version() string {
	return "1.0.0"
}

func (c *ColorCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *ColorCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "color",
			Description: "color to convert",
			Required:    true,
		},
	}
}

func (c *ColorCommand) IsDmCapable() bool {
	return true
}

func (c *ColorCommand) Run(ctx ken.Context) (err error) {

	var (
		colorType  string
		colorInput string

		colSepHue string
		colSepSat string
		colSepVal string

		colSepRed   string
		colSepGreen string
		colSepBlue  string

		colSepHueF64 float64
		colSepSatF64 float64
		colSepValF64 float64

		colSepRedInt   int
		colSepGreenInt int
		colSepBlueInt  int
		colSepRedF64   float64
		colSepGreenF64 float64
		colSepBlueF64  float64

		colorfulColor colorful.Color

		res error
	)

	colorInput = strings.ToLower(ctx.Options().GetByName("color").StringValue())

	if strings.Contains(colorInput, "#") || strings.Contains(colorInput, "hex") {
		colInTrunc := strings.ReplaceAll(strings.ReplaceAll(colorInput, "#", ""), "hex", "")
		colorfulColor, res = colorful.Hex("#" + colInTrunc)

		if res != nil {
			colorType = "invalidHexcode"
		} else {
			colorType = "hex"
		}
	} else if strings.Contains(colorInput, "[") || strings.Contains(colorInput, "rgb") {
		colorType = "rgb"
		colInTrunc := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(colorInput, "rgb", ""), ",", ""), "]", ""), "[", "")

		fmt.Sscanf(colInTrunc, "%s %s %s", &colSepRed, &colSepGreen, &colSepBlue)
		// string to int
		colSepRedInt, _ = strconv.Atoi(colSepRed)
		colSepGreenInt, _ = strconv.Atoi(colSepGreen)
		colSepBlueInt, _ = strconv.Atoi(colSepBlue)

		// int to float64
		colSepRedF64, _ = strconv.ParseFloat(strconv.Itoa(colSepRedInt), 64)
		colSepGreenF64, _ = strconv.ParseFloat(strconv.Itoa(colSepGreenInt), 64)
		colSepBlueF64, _ = strconv.ParseFloat(strconv.Itoa(colSepBlueInt), 64)
		colorfulColor.R = colSepRedF64 / 255
		colorfulColor.G = colSepGreenF64 / 255
		colorfulColor.B = colSepBlueF64 / 255

		if colorfulColor.R > 1 || colorfulColor.G > 1 || colorfulColor.B > 1 || colorfulColor.R < 0 || colorfulColor.G < 0 || colorfulColor.B < 0 {
			colorType = "invalidColorRGB"
		}

	} else if strings.Contains(colorInput, "%") || strings.Contains(colorInput, "hsv") || strings.Contains(colorInput, "deg") {
		colorType = "hsv"
		colInTrunc := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(colorInput, ",", ""), "]", ""), "[", ""), "deg", ""), "hsv", ""), "%", "")

		fmt.Sscanf(colInTrunc, "%s %s %s", &colSepHue, &colSepSat, &colSepVal)

		colSepHueF64, _ = strconv.ParseFloat(colSepHue, 64)
		colSepSatF64, _ = strconv.ParseFloat(colSepSat, 64)
		colSepValF64, _ = strconv.ParseFloat(colSepVal, 64)
		if colSepSatF64 > 1 {
			colSepSatF64 = colSepSatF64 / 100
		}
		if colSepValF64 > 1 {
			colSepValF64 = colSepValF64 / 100
		}
		colorfulColor = colorful.Hsv(colSepHueF64, colSepSatF64, colSepValF64)

		if colSepHueF64 > 360 || colSepHueF64 < 0 || colSepSatF64 < 0 || colSepSatF64 > 100 || colSepValF64 < 0 || colSepValF64 > 100 {
			colorType = "invalidColorHSV"
		}
	} else {
		colorType = "invalidColorType"
	}

	var HueStr, SatStr, ValStr = colorfulColor.Hsv()
	var hexcolor, _ = strconv.ParseInt(strings.ReplaceAll(colorfulColor.Hex(), "#", ""), 16, 64)
	if strings.Contains(colorType, "hsv") || strings.Contains(colorType, "rgb") || strings.Contains(colorType, "hex") {
		err = ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: "Here's your color!",
						Color: int(hexcolor),
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL:    "https://singlecolorimage.com/get/" + strconv.FormatInt(hexcolor, 16) + "/100x100",
							Width:  100,
							Height: 100,
						},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "Hex",
								Value: colorfulColor.Hex(),
							},
							{
								Name:  "RGB",
								Value: "[" + fmt.Sprintf("%.0f", colorfulColor.R*255) + ", " + fmt.Sprintf("%.0f", colorfulColor.G*255) + ", " + fmt.Sprintf("%.0f", colorfulColor.B*255) + "]",
							},
							{
								Name:  "HSV",
								Value: fmt.Sprintf("%.0f", HueStr) + "°, " + fmt.Sprintf("%.0f", SatStr*100) + "%, " + fmt.Sprintf("%.0f", ValStr*100) + "%",
							},
						},
					},
				},
			},
		})
	} else if colorType == "invalidColorType" {
		err = ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Sorry, I couldn't figure out what color type this is. Try prefixing with `hex`, `rgb`, or `hsv`.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	} else if colorType == "invalidColorRGB" {
		err = ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "One or more values in the color you entered are either too high or too low. Remember that RGB values have a range between 0 and 255.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	} else if colorType == "invalidColorHSV" {
		err = ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "One or more values in the color you entered are either too high or too low. Remember that hue is an angle between 0 and 360 degrees, and that saturation and value (brightness) lie between 0-100%.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	} else if colorType == "invalidHexcode" {
		err = ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "The hex code you entered is invalid. Remember that hex codes only use 0-9 and A-F as digits.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

	if err != nil {
		return
	}

	return err
}
