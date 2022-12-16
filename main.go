package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/bwmarrin/discordgo"
	"github.com/lucasb-eyer/go-colorful"
)

// IDEAS
// fun facts
// generate bitmap from image
// autogenerate css font change

// BotToken Flags
var (
	BotToken = flag.String("token", "", "Bot token")
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

type ExtendedCommand struct {
	applicationcommand *discordgo.ApplicationCommand
}

type RandomDevExcuse struct {
	Data string `json:"data"`
}

func init() {
	*BotToken = ReadTokenFromFile("token.txt")
	if *BotToken != "" {
		log.Println("Token read from file")
	} else {
		log.Println("Token not read from file, fetching from env")
		*BotToken = os.Getenv("TOKEN")
	}
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func ReadTokenFromFile(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
		}
	}(f)
	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil {

	}
	buf = buf[:n]
	return string(buf)
}

func init() {
	flag.Parse()
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

var (
	commands = []*ExtendedCommand{
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "ping",
				Description: "Ping",
			},
		},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "echo",
				Description: "make the bot relay a message to a channel",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionChannel,
						Name:        "channel",
						Description: "The channel to send the message to",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "content",
						Description: "The content of the message",
						Required:    true,
					},
				},
			},
		},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "devexcuse",
				Description: "Grabs a random excuse from developerexcuses.com",
			},
		},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "d20",
				Description: "Roll a d20",
			},
		},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "pke",
				Description: "Shows an explanation about PluralKit and multiplicity",
			},
		},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "randomcolor",
				Description: "Generates a random color",
			},
		},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "color",
				Description: "Generates an embed with info about a color",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "color",
						Description: "The color you want to convert (Hex/RGB/HSV)",
						Required:    true,
					},
				},
			},
		},
		//{
		//	applicationcommand: &discordgo.ApplicationCommand{
		//		Name:        "coinflip",
		//		Description: "Flips a coin.",
		//		//Options: []*discordgo.ApplicationCommandOption{
		//		//	{
		//		//		Type: discordgo.ApplicationCommandOptionString,
		//		//		Name: "coin",
		//		//		Description: "The coin you wish to flip, defaults to EUR",
		//		//		Required: true,
		//		//	},
		//		//},
		//	},
		//},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "owoify",
				Description: "owoify a string",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "string",
						Description: "text to owoify",
						Required:    true,
					},
				},
			},
		},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "bitmap",
				Description: "convert an image to a text bitmap (JPG/PNG/GIF, not animated)",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionAttachment,
						Name:        "image",
						Description: "image to convert",
						Required:    true,
					},
				},
			},
		},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "userinfo",
				Description: "shows a whole bunch of info about a user in a neatly formatted embed",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "the user you want to get info of",
						Required:    true,
					},
				},
			},
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong! Latency: " + s.HeartbeatLatency().String(),
				},
			})
			if err != nil {
				return
			}
		},
		"echo": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			if i.Interaction.Member.Permissions&discordgo.PermissionManageMessages == discordgo.PermissionManageMessages {
				options := i.ApplicationCommandData().Options

				optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
				for _, opt := range options {
					optionMap[opt.Name] = opt
				}

				var (
					channelobj     *discordgo.Channel
					messagecontent string
				)

				// Get the value from the option map.
				// When the option exists, ok = true
				if option, ok := optionMap["channel"]; ok {
					// Option values must be type asserted from interface{}.
					// Discordgo provides utility functions to make this simple.
					channelobj = option.ChannelValue(s)
				}

				if opt, ok := optionMap["content"]; ok {
					messagecontent = opt.StringValue()
				}

				_, err := s.ChannelMessageSendComplex(channelobj.ID, &discordgo.MessageSend{
					Content: messagecontent,
				})
				if err != nil {
					return
				}

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Echo sent to %s:\n```%s```", channelobj.Mention(), messagecontent),
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					return
				}
			} else {
				return
			}

		},
		"devexcuse": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var randomexcuse = GetDevExcuse()
			var randomexcuseQuoted = `"` + randomexcuse.Data + `"`
			var err error
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: randomexcuseQuoted,
				},
			})
			if err != nil {
				return
			}
		},
		"d20": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			rand.Seed(time.Now().UnixNano())
			var err error
			var diceroll int = rand.Intn(20)
			dicerolls := [20]string{`Luck really isn't on your side today, huh? It's a 1.`, `A 2. Couldn't have been much worse.`, `It's a tree!- Oh, wait. A 3.`, `A four-se. Of course. That didn't work, did it?`, `A 5. Nothing funny here.`, `A 6. The devil, anyone?`, `Lucky number 7! Now can you get two more?`, `8, not bad.`, `Just under halfway up. A 9`, `A 10! Halfway up the scale!`, `11. Decent.`, `12. Could have been much worse. Could've also been better, though.`, `13. Feelin' lucky?`, `Aand it's come up 14!`, `15! Getting up there!`, `16, solid.`, `17. Rolling real high now, aren't you?`, `18! You're old eno- wait this isn't a birthday.`, `19! So CLOSE!`, `NAT 20 BAYBEE!`}
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: dicerolls[diceroll],
				},
			})
			if err != nil {
				return
			}
		},
		"pke": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var err error
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Why does that person have the [BOT] tag?",
							Color:       0x8c1bb1,
							Description: "On this server we cater to people who experience plurality (or multiplicity) \nBots like PluralKit help these people express themselves easier. \n \nDue to discord limitations, these people show up with the [BOT] tag, but rest assured they are not bots.",
							Fields: []*discordgo.MessageEmbedField{
								&discordgo.MessageEmbedField{
									Name:  "A brief explanation of plurality",
									Value: "There's a lot more to it than what I'm showing here, but essentially plurality is the existence of multiple self-aware entities (they don't necessarily have to be people) in one brain. It's like having roommates inside your head.",
								},
								&discordgo.MessageEmbedField{
									Name:  "Better explanations and more resources",
									Value: `[MoreThanOne](https://morethanone.info/)`,
								},
							},
						},
					},
				},
			})
			if err != nil {
				return
			}
		},
		"randomcolor": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			// RGB base to convert from
			rand.Seed(time.Now().UnixNano())
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

			var err error

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
								&discordgo.MessageEmbedField{
									Name:  "Hex",
									Value: randomcolorRGB.Hex(),
								},
								&discordgo.MessageEmbedField{
									Name:  "RGB",
									Value: "[" + randomRedReadable + ", " + randomGreenReadable + ", " + randomBlueReadable + "]",
								},
								&discordgo.MessageEmbedField{
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
		},
		"color": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var err, res error

			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

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
			)

			// Get the value from the option map.
			// When the option exists, ok = true
			if option, ok := optionMap["color"]; ok {
				// Option values must be type asserted from interface{}.
				// Discordgo provides utility functions to make this simple.
				colorInput = strings.ToLower(option.StringValue())
			}

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
				colSepRedInt, res = strconv.Atoi(colSepRed)
				colSepGreenInt, res = strconv.Atoi(colSepGreen)
				colSepBlueInt, res = strconv.Atoi(colSepBlue)

				// int to float64
				colSepRedF64, res = strconv.ParseFloat(strconv.Itoa(colSepRedInt), 64)
				colSepGreenF64, res = strconv.ParseFloat(strconv.Itoa(colSepGreenInt), 64)
				colSepBlueF64, res = strconv.ParseFloat(strconv.Itoa(colSepBlueInt), 64)
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

				colSepHueF64, res = strconv.ParseFloat(colSepHue, 64)
				colSepSatF64, res = strconv.ParseFloat(colSepSat, 64)
				colSepValF64, res = strconv.ParseFloat(colSepVal, 64)
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
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
									&discordgo.MessageEmbedField{
										Name:  "Hex",
										Value: colorfulColor.Hex(),
									},
									&discordgo.MessageEmbedField{
										Name:  "RGB",
										Value: "[" + fmt.Sprintf("%.0f", colorfulColor.R*255) + ", " + fmt.Sprintf("%.0f", colorfulColor.G*255) + ", " + fmt.Sprintf("%.0f", colorfulColor.B*255) + "]",
									},
									&discordgo.MessageEmbedField{
										Name:  "HSV",
										Value: fmt.Sprintf("%.0f", HueStr) + "°, " + fmt.Sprintf("%.0f", SatStr*100) + "%, " + fmt.Sprintf("%.0f", ValStr*100) + "%",
									},
								},
							},
						},
					},
				})
			} else if colorType == "invalidColorType" {
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry, I couldn't figure out what color type this is. Try prefixing with `hex`, `rgb`, or `hsv`.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			} else if colorType == "invalidColorRGB" {
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "One or more values in the color you entered are either too high or too low. Remember that RGB values have a range between 0 and 255.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			} else if colorType == "invalidColorHSV" {
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "One or more values in the color you entered are either too high or too low. Remember that hue is an angle between 0 and 360 degrees, and that saturation and value (brightness) lie between 0-100%.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			} else if colorType == "invalidHexcode" {
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
		},
		"owoify": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var err error

			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			var originalString string
			if option, ok := optionMap["string"]; ok {
				// Option values must be type asserted from interface{}.
				// Discordgo provides utility functions to make this simple.
				originalString = strings.ToLower(option.StringValue())
			}

			var owoifiedString = `"` + strings.ReplaceAll(originalString, "o", "owo") + `"`

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: owoifiedString,
				},
			})

			if err != nil {
				return
			}
		},
		"userinfo": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var (
				err error

				user        discordgo.User
				color       int
				avatarURL   string
				bannerURL   string
				uid         string
				uname       string
				discrim     string
				bannervalue string
				//userflags discordgo.UserFlags
			)

			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if option, ok := optionMap["user"]; ok {
				user = *option.UserValue(s)
				color = user.AccentColor
				avatarURL = user.AvatarURL(strconv.Itoa(1024))
				bannerURL = user.BannerURL(strconv.Itoa(1024))
				uname = user.Username
				discrim = user.Discriminator
				//userflags = user.PublicFlags
				uid = user.ID
			}

			if bannerURL == "" {
				bannervalue = uname + " does not have a banner set."
			} else {
				bannervalue = "[Link](" + bannerURL + ")"
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
								&discordgo.MessageEmbedField{
									Name:   "Username",
									Value:  uname,
									Inline: true,
								},
								&discordgo.MessageEmbedField{
									Name:   "Discriminator",
									Value:  discrim,
									Inline: true,
								},
								&discordgo.MessageEmbedField{
									Name:   "User ID",
									Value:  uid,
									Inline: true,
								},
								&discordgo.MessageEmbedField{
									Name:  "Avatar",
									Value: `[Link](` + avatarURL + `)`,
								},
								&discordgo.MessageEmbedField{
									Name:  "Banner",
									Value: bannervalue,
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

		},
	}
)

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v.applicationcommand)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.applicationcommand.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer func(s *discordgo.Session) {
		err := s.Close()
		if err != nil {

		}
	}(s)
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{{Type: 3, Name: "over the Den"}},
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Graceful shutdown")
}

func GetDevExcuse() *RandomDevExcuse {
	resp, err := http.Get(("https://api.tabliss.io/v1/developer-excuses"))
	if err != nil {
		log.Fatalf("An error occured: %v", err)
	}
	var randomexcuse *RandomDevExcuse
	err = json.NewDecoder(resp.Body).Decode(&randomexcuse)
	if err != nil {
		var invalid *RandomDevExcuse
		return invalid
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)
	return randomexcuse
}
