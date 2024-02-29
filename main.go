package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"io"
	"log"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	//_ "image/draw"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	_ "image/png"
	_ "net/url"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/bwmarrin/discordgo"
	"github.com/cavaliergopher/grab/v3"
	"github.com/disintegration/gift"

	//
	//htgotts "github.com/hegedustibor/htgo-tts"
	//handlers "github.com/hegedustibor/htgo-tts/handlers"
	//voices "github.com/hegedustibor/htgo-tts/voices"
	//"github.com/jonas747/dca"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/webp"
)

// IDEAS
// gpt integration ohdeargod (agony while keeping 24/7)
// useless web
// randomproto
// server stats (meh, dont see the point)
// polls (needs db'ing)
// reminders (needs db'ing)
// mod stuffs
// customizable interaction accent color
// fun math things
// timestamp

// BotToken Flags
var (
	BotToken          = flag.String("token", "", "Bot token")
	CumulonimbusToken = flag.String("cumulonimbus-token", "", "Cumulonimbus Token")
	RapidSzToken      = flag.String("rapid-sz-token", "", "RapidSz Token")
	UnsplashToken     = flag.String("unsplash-token", "", "Unsplash token")
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var uploadClient http.Client

var buildstring string = " b240229"

//dw := imagick.NewDrawingWand()

var s *discordgo.Session
var vc *discordgo.VoiceConnection
var inVC bool

type ExtendedCommand struct {
	applicationcommand *discordgo.ApplicationCommand
}

type RandomDevExcuse struct {
	Data string `json:"data"`
}

func init() {
	*BotToken = ReadTokenFromFile("token-dev.txt")
	//*BotToken = ReadTokenFromFile("token.txt")
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
	*CumulonimbusToken = ReadTokenFromFile("uploader-token.txt")
	if *CumulonimbusToken != "" {
		log.Println("Uploader token read from file")
	} else {
		log.Println("Token not read from file, fetching from env")
		*BotToken = os.Getenv("UPLOADER_TOKEN")
	}
	*RapidSzToken = ReadTokenFromFile("rapid-sz-token.txt")
	if *RapidSzToken != "" {
		log.Println("RapidAPI token read from file")
	} else {
		log.Println("Token not read from file, fetching from env")
		*RapidSzToken = os.Getenv("RAPID_SZ_TOKEN")
	}
	*UnsplashToken = ReadTokenFromFile("unsplash-token.txt")
	if *UnsplashToken != "" {
		log.Println("Unsplash token read from file")
	} else {
		log.Println("Token not read from file, fetching from env")
		*UnsplashToken = os.Getenv("UNSPLASH_TOKEN")
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
		//{
		//	applicationcommand: &discordgo.ApplicationCommand{
		//		Name:        "bitmap",
		//		Description: "convert an image to a text bitmap (JPG/PNG/GIF, not animated)",
		//		Options: []*discordgo.ApplicationCommandOption{
		//			{
		//				Type:        discordgo.ApplicationCommandOptionAttachment,
		//				Name:        "image",
		//				Description: "image to convert",
		//				Required:    true,
		//			},
		//		},
		//	},
		//},
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
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "rps",
				Description: "play rock paper scissors with the bot",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "choice",
						Description: "your move",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{
								Name:  "rock",
								Value: "rock",
							},
							{
								Name:  "paper",
								Value: "paper",
							},
							{
								Name:  "scissors",
								Value: "scissors",
							},
						},
					},
				},
			},
		},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "askevi",
				Description: "ask a question and get an answer, magic 8-ball style",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "question",
						Description: "the question you want to ask me",
						Required:    true,
					},
				},
			},
		},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "magick",
				Description: "mess with an image through imagemagick",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "url",
						Description: "The image you want to blur",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "method",
						Description: "modification you want to make",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{
								Name:  "blur (intensity in px)",
								Value: "blur",
							},
							{
								Name:  "contrast (intensity in %, between -100 and 100)",
								Value: "contrast",
							},
							{
								Name:  "saturation (intensity in %, between -100 and 500)",
								Value: "saturation",
							},
							{
								Name:  "sobel (intensity not used)",
								Value: "sobel",
							},
							{
								Name:  "hue shift (intensity in degrees)",
								Value: "hueshift",
							},
							{
								Name:  "invert (intensity not used)",
								Value: "invert",
							},
							{
								Name:  "pixelate (intensity in px)",
								Value: "pixelate",
							},
							{
								Name:  "sepia (intensity in %, between 0 and 100)",
								Value: "sepia",
							},
							{
								Name:  "brightness (intensity in %, between -100 and 100)",
								Value: "brightness",
							},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionNumber,
						Name:        "intensity",
						Description: "Value of whatever filter you're using, highly context-sensitive",
					},
				},
			},
		},
		// music cmd's, oh feck
		//{
		//	applicationcommand: &discordgo.ApplicationCommand{
		//		Name:        "play",
		//		Description: "play a song",
		//		Options: []*discordgo.ApplicationCommandOption{
		//			{
		//				Type:        discordgo.ApplicationCommandOptionString,
		//				Name:        "url",
		//				Description: "the video you want to play",
		//				Required:    true,
		//			},
		//		},
		//	},
		//},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "join",
				Description: "joins your voice chat",
			},
		},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "leave",
				Description: "disconnects from the voice channel",
			},
		},
		//{
		//	applicationcommand: &discordgo.ApplicationCommand{
		//		Name:        "play",
		//		Description: "play a song",
		//		Options: []*discordgo.ApplicationCommandOption{
		//			{
		//				Type:        discordgo.ApplicationCommandOptionString,
		//				Name:        "url",
		//				Description: "the video you want to play",
		//				Required:    true,
		//			},
		//		},
		//	},
		//},
		//{
		//	applicationcommand: &discordgo.ApplicationCommand{
		//		Name:        "join",
		//		Description: "joins your voice chat",
		//	},
		//},
		//{
		//	applicationcommand: &discordgo.ApplicationCommand{
		//		Name:        "leave",
		//		Description: "disconnects from the voice channel",
		//	},
		//},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "songinfo",
				Description: "Returns song info from the Spotify API",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "query",
						Description: "search query",
						Required:    true,
					},
					//					{
					//						Type:        discordgo.ApplicationCommandOptionString,
					//						Name:        "type",
					//						Description: "type of search, defaults to all",
					//						Choices: []*discordgo.ApplicationCommandOptionChoice{
					//							{
					//								Name:  "multi",
					//								Value: "multi",
					//							},
					//							{
					//								Name:  "albums",
					//								Value: "albums",
					//							},
					//							{
					//								Name:  "artists",
					//								Value: "artists",
					//							},
					//							{
					//								Name:  "tracks",
					//								Value: "tracks",
					//							},
				},
			},
		},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "google",
				Description: "Let me google that for you!",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "query",
						Description: "search query",
						Required:    true,
					},
				},
			},
		},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "statusset",
				Description: "Set the bot's status",
				Options: []*discordgo.ApplicationCommandOption{
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
				},
			},
		},
		//{
		//	applicationcommand: &discordgo.ApplicationCommand{
		//		Name:        "settings accentcolor",
		//		Description: "Sets the accent color for commands without variable colors",
		//		Options: []*discordgo.ApplicationCommandOption{
		//			{
		//				Type:        discordgo.ApplicationCommandOptionString,
		//				Name:        "hexcode",
		//				Description: "hexcode of the accent color",
		//				Required:    true,
		//			},
		//		},
		//	},
		//},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "unsplash",
				Description: "Get an image from unsplash",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "query",
						Description: "search term",
						Required:    true,
					},
				},
			},
		},
		//{
		//	applicationcommand: &discordgo.ApplicationCommand{
		//		Name:        "math",
		//		Description: "do some math",
		//		Options: []*discordgo.ApplicationCommandOption{
		//			{
		//				Type:        discordgo.ApplicationCommandOptionString,
		//				Name:        "input",
		//				Description: "math input",
		//				Required:    true,
		//			},
		//		},
		//	},
		//},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "timestamp",
				Description: "generate a hoverable timestamp",
				Options: []*discordgo.ApplicationCommandOption{
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
				},
			},
		},
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "ticket",
				Description: "Open a ticket",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "title",
						Description: "ticket title",
						Required:    true,
					},
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ChannelLog(fmt.Sprintf("/// command `ping` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))

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

			ChannelLog(fmt.Sprintf("/// command `echo` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))

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

				//fmt.Sprintf("Echo sent to %v in %v: %v by %v#%v", channelobj.Name, channelobj.GuildID, messagecontent, i.User.Username, i.User.Discriminator)
				ChannelLog(fmt.Sprintf("Echo sent to %v in %v: %v by %v#%v", channelobj.Name, channelobj.GuildID, messagecontent, i.Member.User.Username, i.Member.User.Discriminator))

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
			ChannelLog(fmt.Sprintf("/// command `devexcuse` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
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
			ChannelLog(fmt.Sprintf("/// command `d20` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
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
			ChannelLog(fmt.Sprintf("/// command `pke` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
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
								{
									Name:  "A brief explanation of plurality",
									Value: "There's a lot more to it than what I'm showing here, but essentially plurality is the existence of multiple self-aware entities (they don't necessarily have to be people) in one brain. It's like having roommates inside your head.",
								},
								{
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
			ChannelLog(fmt.Sprintf("/// command `randomcolor` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
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
		},
		"color": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ChannelLog(fmt.Sprintf("/// command `color` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
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
			ChannelLog(fmt.Sprintf("/// command `owoify` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
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
			ChannelLog(fmt.Sprintf("/// command `userinfo` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
			var (
				err error

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
				//userflags       discordgo.UserFlags
				//userflagbitmask string
				//stringLength    int
				//charIndex       int
				//badgeList       []string
				//badgeCount      int
				//badgeString     string
				//bitmaskSplit    []string
				//bitmaskSeg      []int
				//badgeNum        int
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
				mention = user.Mention()
			}

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

			// badges ohno
			//userflagbitmask = fmt.Sprintf("%b", userflags) // bitmask to string
			//badgeCount = strings.Count(userflagbitmask, "1")
			//stringLength = len(userflagbitmask)
			//bitmaskSplit = strings.SplitAfterN(fmt.Sprint(userflagbitmask), "1", -1) // split after each "1", return as many as there are... this took three hours
			//ChannelLog(fmt.Sprintf(uname)
			//ChannelLog(fmt.Sprintf(userflagbitmask)
			//ChannelLog(fmt.Sprintf("length: " + strconv.FormatInt(int64(stringLength), 10))
			//badgeNum = 1
			//for badgeNum < badgeCount {
			//	for i, s := range bitmaskSplit {
			//		ChannelLog(fmt.Sprintf(i, s)
			//		bitmaskSeg = append(bitmaskSeg, strings.Count(s, "0"))
			//		badgeIndexReverse := [23]string{"<:activedev:1054431966960832542>", "<:activedev:1054431966960832542>", "nul", "nul", "nul", "<:modalumni:1054431972996427777>", "<:verifieddev:1054431968160399381>", "nul", "<:bughunter2:1054431970911862914>", "nul", "nul", "nul", "nul", "<:earlysupporter:1054431972002377728>", "<:balance:1054431978314797136>", "<:brilliance:1054432130391887952>", "<:bravery:1054432128856756224>", "nul", "nul", "<:bughunter1:1054431969674526781>", "<:hypesquad:1054431981255012483>", "<:partner:1054431975945011250>", "<:staff:1054432180803227740>"}
			//		badgeIndex := [23]string{"<:staff:1054432180803227740>", "<:partner:1054431975945011250>", "<:hypesquad:1054431981255012483>", "<:bughunter1:1054431969674526781>", "nul", "nul", "<:bravery:1054432128856756224>", "<:brilliance:1054432130391887952>", "<:balance:1054431978314797136>", "<:earlysupporter:1054431972002377728>", "nul", "nul", "nul", "nul", "<:bughunter2:1054431970911862914>", "nul", "nul", "<:verifieddev:1054431968160399381>", "<:modalumni:1054431972996427777>", "nul", "nul", "nul", "<:activedev:1054431966960832542>"}
			//		ChannelLog(fmt.Sprintf(bitmaskSeg)
			//		ChannelLog(fmt.Sprintf(bitmaskSeg[len(bitmaskSeg)-1])
			//		iMax := bitmaskSplit[len(bitmaskSplit)-1]
			//		if i == 0 {
			//			if stringLength == 23 {
			//				ChannelLog(fmt.Sprintf(badgeIndexReverse[bitmaskSeg[len(bitmaskSeg)-1]])
			//				if badgeIndexReverse[bitmaskSeg[len(bitmaskSeg)-1]] != "nul" {
			//					badgeList = append(badgeList, badgeIndexReverse[bitmaskSeg[len(bitmaskSeg)-1]])
			//				} else {
			//					badgeList = badgeList
			//				}
			//			}
			//		} else if i > 0 && strconv.FormatInt(int64(i), 10) < iMax {
			//			ChannelLog(fmt.Sprintf(badgeIndexReverse[bitmaskSeg[len(bitmaskSeg)-1]+1])
			//			if badgeIndexReverse[bitmaskSeg[len(bitmaskSeg)-1]+1] != "nul" {
			//				badgeList = append(badgeList, badgeIndexReverse[bitmaskSeg[len(bitmaskSeg)-1]])
			//				charIndex = charIndex + bitmaskSeg[len(bitmaskSeg)-1]
			//			} else {
			//				badgeList = badgeList
			//			}
			//		} else {
			//			ChannelLog(fmt.Sprintf(badgeIndex[bitmaskSeg[len(bitmaskSeg)-1]])
			//			if badgeIndex[bitmaskSeg[len(bitmaskSeg)-1]] != "nul" {
			//				badgeList = append(badgeList, badgeIndex[bitmaskSeg[len(bitmaskSeg)-1]])
			//				charIndex = charIndex + bitmaskSeg[len(bitmaskSeg)-1]
			//			} else {
			//				badgeList = badgeList
			//			}
			//		}
			//		charIndex = charIndex + bitmaskSeg[len(bitmaskSeg)-1]
			//		ChannelLog(fmt.Sprintf("charIndex " + strconv.FormatInt(int64(charIndex), 10))
			//	}
			//
			//	ChannelLog(fmt.Sprintf(badgeList)
			//	ChannelLog(fmt.Sprintf(badgeString)
			//
			//	badgeString = strings.Join(badgeList, " ")
			//	if badgeString == "" {
			//		badgeString = uname + " does not have any badges."
			//	}
			//	badgeNum += 1
			//}

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
								//{
								//	Name:   "Badges",
								//	Value:  badgeString,
								//	Inline: true,
								//},
								//{
								//	Name:   "⠀",
								//	Value:  "⠀",
								//	Inline: true,
								//},
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

		},
		"rps": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ChannelLog(fmt.Sprintf("/// command `rps` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
			rand.Seed(time.Now().UnixNano())
			var (
				err error

				botChoiceID int = rand.Intn(3)
				botChoice   string
				choice      string
				resultID    int // 0 = user win, 1 = bot win, 2 = tie
			)

			choices := [3]string{"rock", "paper", "scissors"}

			//responses := [20]string{`Luck really isn't on your side today, huh? It's a 1.`, `A 2. Couldn't have been much worse.`, `It's a tree!- Oh, wait. A 3.`, `A four-se. Of course. That didn't work, did it?`, `A 5. Nothing funny here.`, `A 6. The devil, anyone?`, `Lucky number 7! Now can you get two more?`, `8, not bad.`, `Just under halfway up. A 9`, `A 10! Halfway up the scale!`, `11. Decent.`, `12. Could have been much worse. Could've also been better, though.`, `13. Feelin' lucky?`, `Aand it's come up 14!`, `15! Getting up there!`, `16, solid.`, `17. Rolling real high now, aren't you?`, `18! You're old eno- wait this isn't a birthday.`, `19! So CLOSE!`, `NAT 20 BAYBEE!`}
			result := [3]string{"You win", "I win", "It's a tie"}

			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if option, ok := optionMap["choice"]; ok {
				choice = strings.ToLower(option.StringValue())
			}

			if choice != "rock" && choice != "paper" && choice != "scissors" {
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "invalid choice",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			}

			botChoice = choices[botChoiceID]

			if choice == "rock" {
				if botChoice == "rock" {
					resultID = 2
				} else if botChoice == "paper" {
					resultID = 1
				} else {
					resultID = 0
				}
			} else if choice == "paper" {
				if botChoice == "rock" {
					resultID = 0
				} else if botChoice == "paper" {
					resultID = 2
				} else {
					resultID = 1
				}
			} else {
				if botChoice == "rock" {
					resultID = 1
				} else if botChoice == "paper" {
					resultID = 0
				} else {
					resultID = 2
				}
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: result[resultID] + " // My pick: " + botChoice,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})

			if err != nil {
				return
			}
		},
		"askevi": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ChannelLog(fmt.Sprintf("/// command `askevi` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
			rand.Seed(time.Now().UnixNano())
			var err error
			var color int
			var query string
			var responseID int

			responses := [20]string{"It is certain.", "It is decidedly so.", "Without a doubt.", "Yes, definitely.", "You may rely on it.", "As I see it, yes.", "Most likely.", "Outlook good.", "Yes.", "Signs point to yes", "Reply hazy, try again.", "Ask again later.", "Better not tell you now.", "Cannot predict now.", "Concentrate and ask again.", "Don't count on it.", "My reply is no.", "My sources say no.", "Outlook not so good.", "Very doubtful."}

			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if option, ok := optionMap["question"]; ok {
				// Option values must be type asserted from interface{}.
				// Discordgo provides utility functions to make this simple.
				ChannelLog(fmt.Sprintf(option.StringValue()))
				query = strings.ToLower(option.StringValue())
			}

			//femboyBiasTriggerMox := [3]string{"i", "femboy", "am"}

			fmt.Println(i.Member.User.ID)
			fmt.Println(strings.Contains(query, "i"))
			fmt.Println(strings.Contains(query, "femboy"))

			//if i.Member.User.ID == "758810453622259743" && strings.Contains(query, "i") && strings.Contains(query, "femboy") {
			//	fmt.Println(strings.Count(query, "not") % 2)
			//	if strings.Count(query, "not")%2 == 0 {
			//		responseID = rand.Intn(9)
			//	} else {
			//		responseID = rand.Intn(5) + 14
			//	}
			//} else if strings.Contains(query, "mox") && strings.Contains(query, "femboy") {
			//	if strings.Count(query, "not")%2 == 0 {
			//		responseID = rand.Intn(9)
			//	} else {
			//		responseID = rand.Intn(5) + 14
			//	}
			responseID = rand.Intn(19)

			if responseID < 10 {
				color = 0x00ae00 //affirmative
			} else if responseID > 9 && responseID < 14 {
				color = 0xfcb103 //non-commital
			} else {
				color = 0xff0000 //negative
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: responses[responseID],
							Color: color,
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
		},
		"magick": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ChannelLog(fmt.Sprintf("/// command `magick` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
			var (
				err error

				imgUrl         string
				method         string
				dst            *image.RGBA
				src            image.Image
				fEdit          *os.File
				execFolder     string
				execPath       string
				imagePath      string
				uploadResponse CumulonimbusResponse
				colorfulCol    colorful.Color
				intensity      float32
			)

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "image editing in process... please stand by.",
				},
			})

			if err != nil {
				ChannelLog(fmt.Sprintf("an error occurred sending initial message: %v", err))
			}

			execPath, err = os.Executable()
			execFolder = filepath.Dir(execPath)

			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if option, ok := optionMap["url"]; ok {
				imgUrl = option.StringValue()
				ChannelLog(fmt.Sprintf(imgUrl))
			}
			if opt, ok := optionMap["method"]; ok {
				method = opt.StringValue()
				ChannelLog(fmt.Sprintf(method))
			}
			if opt2, ok := optionMap["intensity"]; ok {
				intensity = float32(opt2.FloatValue())
			}
			resp, err := grab.Get(".", imgUrl)
			if err != nil {
				log.Fatal(err)
			}

			ChannelLog(fmt.Sprintf("Download saved to %v ", resp.Filename))
			file := resp.Filename
			buf := make([]byte, 512)

			openFile, err := os.Open(file)
			if err != nil {
				ChannelLog(fmt.Sprintf("An error occurred during file read: %v", err))
			}
			n, err := openFile.Read(buf)
			ChannelLog(fmt.Sprintf("bytes read: `%v`", n))
			if err != nil {
				ChannelLog(fmt.Sprintf("An error occurred during file read to buffer: %v", err))
			}
			contentType := http.DetectContentType(buf)
			ChannelLog(fmt.Sprintf("Content type is `%v`", contentType))
			openFile.Close()
			openFile, err = os.Open(file)

			if contentType == "image/jpeg" {
				src, err = jpeg.Decode(openFile)
			}
			if contentType == "image/png" {
				src, err = png.Decode(openFile)
			}
			if contentType == "image/webp" {
				src, err = webp.Decode(openFile)
			}
			if err != nil {
				ChannelLog(fmt.Sprintf("An error occurred during decode: %v", err))
				return
			}

			if method == "blur" {
				if intensity == 0 {
					intensity = 1.5
				}

				if intensity < 0 {
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseDeferredMessageUpdate,
						Data: &discordgo.InteractionResponseData{
							Content: "what are you trying to do, break me!? (intensity should be more than 0)",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					})
				}
				g := gift.New(
					gift.GaussianBlur(intensity),
				)
				dst = image.NewRGBA(g.Bounds(src.Bounds()))
				g.Draw(dst, src)
			}

			if method == "contrast" {
				if intensity == 0 {
					intensity = 50
				}

				if intensity < -100 || intensity > 100 {
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseDeferredMessageUpdate,
						Data: &discordgo.InteractionResponseData{
							Content: "what are you trying to do, break me!? (intensity should be between -100 and 100)",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					})
				}

				g := gift.New(
					gift.Contrast(intensity),
				)
				dst = image.NewRGBA(g.Bounds(src.Bounds()))
				g.Draw(dst, src)
			}

			if method == "saturation" {
				if intensity == 0 {
					intensity = 50
				}

				if intensity < -100 || intensity > 500 {
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseDeferredMessageUpdate,
						Data: &discordgo.InteractionResponseData{
							Content: "what are you trying to do, break me!? (intensity should be between -100 and 500)",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					})
				}

				g := gift.New(
					gift.Saturation(intensity),
				)
				dst = image.NewRGBA(g.Bounds(src.Bounds()))
				g.Draw(dst, src)
			}

			if method == "sobel" {
				g := gift.New(
					gift.Sobel(),
				)
				dst = image.NewRGBA(g.Bounds(src.Bounds()))
				g.Draw(dst, src)
			}

			if method == "hueshift" {
				if intensity == 0 {
					intensity = 45
				}

				g := gift.New(
					gift.Hue(intensity),
				)
				dst = image.NewRGBA(g.Bounds(src.Bounds()))
				g.Draw(dst, src)
			}

			if method == "invert" {
				g := gift.New(
					gift.Invert(),
				)
				dst = image.NewRGBA(g.Bounds(src.Bounds()))
				g.Draw(dst, src)
			}

			if method == "pixelate" {
				intensity = float32(math.Round(float64(intensity)))
				if intensity == 0 {
					intensity = 5
				}

				g := gift.New(
					gift.Pixelate(int(intensity)),
				)
				dst = image.NewRGBA(g.Bounds(src.Bounds()))
				g.Draw(dst, src)
			}

			if method == "sepia" {
				if intensity == 0 {
					intensity = 50
				}

				if intensity < 0 {
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseDeferredMessageUpdate,
						Data: &discordgo.InteractionResponseData{
							Content: "what are you trying to do, break me!? (intensity should be between 0 and 100)",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					})
				}

				g := gift.New(
					gift.Sepia(intensity),
				)
				dst = image.NewRGBA(g.Bounds(src.Bounds()))
				g.Draw(dst, src)
			}

			if method == "brightness" {
				if intensity == 0 {
					intensity = 25
				}

				if intensity < -100 || intensity > 100 {
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseDeferredMessageUpdate,
						Data: &discordgo.InteractionResponseData{
							Content: "what are you trying to do, break me!? (intensity should be between -100 and 100)",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					})
				}

				g := gift.New(
					gift.Brightness(intensity),
				)
				dst = image.NewRGBA(g.Bounds(src.Bounds()))
				g.Draw(dst, src)

			}

			avgcol, err := prominentcolor.Kmeans(dst)
			fmt.Printf("%v %v %v", avgcol[0].Color.R, avgcol[0].Color.G, avgcol[0].Color.B)

			// COLOR CONVERSION YAAAAAAAAAAAAAAAAAAAAAAAAA
			avgColRInt, err := strconv.ParseInt(strconv.FormatUint(uint64(avgcol[0].Color.R), 10), 10, 64)
			avgColR, err := strconv.ParseFloat(strconv.Itoa(int(avgColRInt)), 64)

			avgColGInt, err := strconv.ParseInt(strconv.FormatUint(uint64(avgcol[0].Color.G), 10), 10, 64)
			avgColG, err := strconv.ParseFloat(strconv.Itoa(int(avgColGInt)), 64)

			avgColBInt, err := strconv.ParseInt(strconv.FormatUint(uint64(avgcol[0].Color.B), 10), 10, 64)
			avgColB, err := strconv.ParseFloat(strconv.Itoa(int(avgColBInt)), 64)

			avgColRDiv := avgColR / 255
			avgColGDiv := avgColG / 255
			avgColBDiv := avgColB / 255

			colorfulCol.R = avgColRDiv
			colorfulCol.G = avgColGDiv
			colorfulCol.B = avgColBDiv

			fEdit, err = os.Create("image.png")
			if err != nil {
				ChannelLog(fmt.Sprintf("error while creating image file: %v", err))
			}
			png.Encode(fEdit, dst)

			hexcol, err := strconv.ParseInt(strings.ReplaceAll(colorfulCol.Hex(), "#", ""), 16, 64)

			imagePath = execFolder + "\\image.png"
			ChannelLog(fmt.Sprintf(imagePath))

			_, err = os.Open(imagePath)

			form := new(bytes.Buffer)
			writer := multipart.NewWriter(form)
			fw, err := writer.CreateFormFile("file", filepath.Base("evi-edit-upload.png"))
			if err != nil {
				log.Fatal(err)
			}
			fd, err := os.Open("image.png")
			if err != nil {
				log.Fatal(err)
			}
			defer fd.Close()
			_, err = io.Copy(fw, fd)
			if err != nil {
				log.Fatal(err)
			}

			writer.Close()

			client := &http.Client{}
			req, err := http.NewRequest("POST", "https://alekeagle.me/api/upload", form)
			if err != nil {
				log.Fatal(err)
			}
			req.Header.Set("Authorization", *CumulonimbusToken)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			httpresp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer httpresp.Body.Close()
			bodyText, err := io.ReadAll(httpresp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", bodyText)
			err = json.Unmarshal(bodyText, &uploadResponse)

			if err != nil {
				return
			}

			//_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			//		Embeds: []*discordgo.MessageEmbed{
			//			{
			//				Title: "Here's the edited image!",
			//				Color: int(hexcol),
			//				Image: &discordgo.MessageEmbedImage{
			//					URL:    uploadResponse.URL,
			//					Width:  1024,
			//					Height: 1024,
			//				},
			//				Description: "Upload powered by [Cumulonimbus](https://alekeagle.me)",
			//			},
			//		},
			//	},
			//)
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title: "Here's the edited image!",
						Color: int(hexcol),
						Image: &discordgo.MessageEmbedImage{
							URL:    uploadResponse.URL,
							Width:  1024,
							Height: 1024,
						},
						Description: "Upload powered by [Cumulonimbus](https://alekeagle.me)",
					},
				},
			})

			if err != nil {
				ChannelLog(fmt.Sprintf("an error occurred while sending embed update: %v", err))
			}

			os.Remove(file)
			os.Remove(imagePath)
			if err != nil {
				return
			}

		},
		//"join": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		//	ChannelLog(fmt.Sprintf("/// command `join` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
		//	var (
		//		err      error
		//		vcID     string
		//		guild    *discordgo.Guild = new(discordgo.Guild)
		//		guildID  string
		//		msg      string
		//		msgID    string
		//		msgOldID string
		//		file     string
		//	)
		//
		//	guild, err = s.State.Guild(i.GuildID)
		//	if err != nil {
		//		ChannelLog(fmt.Sprintf("An error occurred during state loading: %v", err))
		//	}
		//	guildID = guild.ID
		//	ChannelLog(fmt.Sprintf(guildID))
		//
		//	for _, voiceState := range guild.VoiceStates {
		//		if voiceState.UserID == i.Member.User.ID {
		//			vcID = voiceState.ChannelID
		//		}
		//	}
		//
		//	vc, err = s.ChannelVoiceJoin(guildID, vcID, false, true)
		//	inVC = true
		//
		//	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		//		Type: discordgo.InteractionResponseChannelMessageWithSource,
		//		Data: &discordgo.InteractionResponseData{
		//			Content: "Joined your voice channel!",
		//			Flags:   discordgo.MessageFlagsEphemeral,
		//		},
		//	})
		//
		//	msgOldID = ""
		//
		//	for vc != nil {
		//		latestmessage, err := s.ChannelMessages("1063016193771982858", 1, "", "", "")
		//		speech := htgotts.Speech{Folder: "audio", Language: voices.English, Handler: &handlers.Native{}}
		//		msgID = latestmessage[0].ID
		//		ChannelLog(fmt.Sprintf(msgID))
		//
		//		if msgID != msgOldID {
		//
		//			msg = latestmessage[0].Content
		//			file, err = speech.CreateSpeechFile(msg, "tts")
		//			msgOldID = msgID
		//			ChannelLog(fmt.Sprintf(file))
		//			vc.Speaking(true)
		//			encodeSession, err := dca.EncodeFile(file, dca.StdEncodeOptions)
		//			defer encodeSession.Cleanup()
		//			output, err := os.Create("output.dca")
		//			if err != nil {
		//				return
		//			}
		//
		//			decoder := dca.NewDecoder(output)
		//
		//			for {
		//				frame, err := decoder.OpusFrame()
		//				if err != nil {
		//					if err != io.EOF {
		//						// Handle the error
		//					}
		//
		//					break
		//				}
		//
		//				// Do something with the frame, in this example were sending it to discord
		//				select {
		//				case vc.OpusSend <- frame:
		//				case <-time.After(time.Second):
		//					// We haven't been able to send a frame in a second, assume the connection is borked
		//					return
		//				}
		//			}
		//		}
		//
		//		if err != nil {
		//			ChannelLog(fmt.Sprintf("An error occurred trying to send OPUS audio data: %v", err))
		//		}
		//	}
		//
		//	if err != nil {
		//		return
		//	}
		//},
		//"leave": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		//	ChannelLog(fmt.Sprintf("/// command `leave` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
		//	if vc != nil {
		//		vc.Disconnect()
		//		vc = nil
		//		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		//			Type: discordgo.InteractionResponseChannelMessageWithSource,
		//			Data: &discordgo.InteractionResponseData{
		//				Content: "Left your voice channel!",
		//				Flags:   discordgo.MessageFlagsEphemeral,
		//			},
		//		})
		//
		//		if err != nil {
		//			return
		//		}
		//	} else {
		//		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		//			Type: discordgo.InteractionResponseChannelMessageWithSource,
		//			Data: &discordgo.InteractionResponseData{
		//				Content: "Not in a voice channel",
		//				Flags:   discordgo.MessageFlagsEphemeral,
		//			},
		//		})
		//
		//		if err != nil {
		//			return
		//		}
		//	}
		//
		//},
		"songinfo": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ChannelLog(fmt.Sprintf("/// command `songinfo` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
			var (
				err          error
				query        string
				parsedQuery  string
				songdata     *RapidSzResponse
				sDuration    int64
				t            time.Time
				colorfulCol  colorful.Color
				artistString string
				// Song info

			)

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "song lookup in progress",
				},
			})

			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if option, ok := optionMap["query"]; ok {
				// Option values must be type asserted from interface{}.
				// Discordgo provides utility functions to make this simple.
				query = option.StringValue()
				parsedQuery = strings.ReplaceAll(query, " ", "%20")
			}
			//if opt, ok := optionMap["type"]; ok {
			//	queryType = opt.StringValue()
			//}

			songdata = GetRapidAPICall(parsedQuery, "tracks")

			if songdata.Tracks.TotalCount == 0 {
				_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Embeds: &[]*discordgo.MessageEmbed{
						{
							Title: "No results.",
						},
					},
				})
				if err != nil {
					ChannelLog(fmt.Sprintf("An error occurred while sending failure response: %v", err))
				}
			} else {
				parsedURL := "https://open.spotify.com/track/" + songdata.Tracks.Items[0].Data.ID
				parsedAlbumURL := "https://open.spotify.com/album/" + songdata.Tracks.Items[0].Data.AlbumOfTrack.ID

				for _, artistName := range songdata.Tracks.Items[0].Data.Artists.Items {
					artistString = artistString + ", [" + artistName.Profile.Name + "](https://open.spotify.com/artist/" + strings.ReplaceAll(artistName.URI, "spotify:artist:", "") + ")"
					ChannelLog(fmt.Sprintf(artistString))
				}

				artistString = strings.Replace(artistString, ", ", "", 1)

				resp, err := grab.Get(".", songdata.Tracks.Items[0].Data.AlbumOfTrack.CoverArt.Sources[2].URL)
				ChannelLog(fmt.Sprintf("Download saved to %v", resp.Filename))
				file := resp.Filename
				f, err := os.Open(file)
				src, err := jpeg.Decode(f)
				avgcol, err := prominentcolor.Kmeans(src)
				ChannelLog(fmt.Sprintf("%v %v %v", avgcol[0].Color.R, avgcol[0].Color.G, avgcol[0].Color.B))

				// COLOR CONVERSION YAAAAAAAAAAAAAAAAAAAAAAAAA
				avgColRInt, err := strconv.ParseInt(strconv.FormatUint(uint64(avgcol[0].Color.R), 10), 10, 64)
				avgColR, err := strconv.ParseFloat(strconv.Itoa(int(avgColRInt)), 64)

				avgColGInt, err := strconv.ParseInt(strconv.FormatUint(uint64(avgcol[0].Color.G), 10), 10, 64)
				avgColG, err := strconv.ParseFloat(strconv.Itoa(int(avgColGInt)), 64)

				avgColBInt, err := strconv.ParseInt(strconv.FormatUint(uint64(avgcol[0].Color.B), 10), 10, 64)
				avgColB, err := strconv.ParseFloat(strconv.Itoa(int(avgColBInt)), 64)

				avgColRDiv := avgColR / 255
				avgColGDiv := avgColG / 255
				avgColBDiv := avgColB / 255

				colorfulCol.R = avgColRDiv
				colorfulCol.G = avgColGDiv
				colorfulCol.B = avgColBDiv

				hexcol, err := strconv.ParseInt(strings.ReplaceAll(colorfulCol.Hex(), "#", ""), 16, 64)

				err = os.Remove(file)
				if err != nil {
					ChannelLog(fmt.Sprintf("An error occurred during file deletion: %v", err))
				}
				sDuration = int64(songdata.Tracks.Items[0].Data.Duration.TotalMilliseconds)
				t = time.UnixMilli(sDuration)
				tParse := t.Format("04:05")

				_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Embeds: &[]*discordgo.MessageEmbed{
						{
							Title: songdata.Tracks.Items[0].Data.Name,
							URL:   parsedURL,
							Color: int(hexcol),
							Thumbnail: &discordgo.MessageEmbedThumbnail{
								URL:    songdata.Tracks.Items[0].Data.AlbumOfTrack.CoverArt.Sources[0].URL,
								Width:  128,
								Height: 128,
							},
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:   "Artist",
									Value:  artistString,
									Inline: true,
								},
								{
									Name:   "Album",
									Value:  "[" + songdata.Tracks.Items[0].Data.AlbumOfTrack.Name + "](" + parsedAlbumURL + ")",
									Inline: true,
								},
								//{
								//	Name: "Release",
								//	Value: songdata.Tracks.Items[0].Data.AlbumOfTrack.
								//},
								{
									Name:  "Duration",
									Value: tParse,
								},
							},
						},
					},
				},
				)

				if err != nil {
					ChannelLog(fmt.Sprintf("An error occurred sending the embed: %v", err))
				}

			}
			if err != nil {
				return
			}

		},
		"google": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ChannelLog(fmt.Sprintf("/// command `google` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
			var (
				err         error
				q           string
				link        string
				parsedQuery string
			)

			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if option, ok := optionMap["query"]; ok {
				// Option values must be type asserted from interface{}.
				// Discordgo provides utility functions to make this simple.
				q = option.StringValue()
				parsedQuery = strings.ReplaceAll(q, " ", "+")
			}

			link = "https://lmgtfy.app/?q=" + parsedQuery

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Here you go: " + link,
				},
			})

			if err != nil {
				return
			}
		},
		//"settings accentcolor": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		//	var (
		//		err error
		//		//udata     DBSettings
		//		uid       string
		//		input     string
		//		parsedCol string
		//		//hexcol    int
		//	)

		//	options := i.ApplicationCommandData().Options

		//	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		//	for _, opt := range options {
		//		optionMap[opt.Name] = opt
		//	}

		//	if option, ok := optionMap["hexcode"]; ok {
		//		// Option values must be type asserted from interface{}.
		//		// Discordgo provides utility functions to make this simple.
		//		input = option.StringValue()
		//		parsedCol = strings.ReplaceAll(input, "#", "0x")
		//	}

		//	uid = i.Interaction.Member.User.ID

		//	os.WriteFile("db/"+uid+".json", []byte("\"Color\": "+parsedCol), 0666)

		//	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		//		Type: discordgo.InteractionResponseChannelMessageWithSource,
		//		Data: &discordgo.InteractionResponseData{
		//			Content: "Color has been set",
		//			Flags:   discordgo.MessageFlagsEphemeral,
		//		},
		//	})

		//	if err != nil {
		//		return
		//	}
		//},
		"statusset": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var (
				err       error
				actType   int
				actString string
			)

			if i.Member.User.ID == "386539887123365909" {
				options := i.ApplicationCommandData().Options
				optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
				for _, opt := range options {
					optionMap[opt.Name] = opt
				}
				if option, ok := optionMap["type"]; ok {
					// Option values must be type asserted from interface{}.
					// Discordgo provides utility functions to make this simple.
					actType = int(option.IntValue())
				}

				if opt, ok := optionMap["string"]; ok {
					actString = opt.StringValue()
				}

				s.UpdateStatusComplex(discordgo.UpdateStatusData{
					Activities: []*discordgo.Activity{{Type: discordgo.ActivityType(actType), Name: actString}},
				})
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
		},
		"unsplash": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ChannelLog(fmt.Sprintf("/// command `unsplash` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
			var (
				err      error
				query    string
				unsplash *UnsplashRandom
			)

			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			if option, ok := optionMap["query"]; ok {
				// Option values must be type asserted from interface{}.
				// Discordgo provides utility functions to make this simple.
				query = option.StringValue()
			}
			unsplash = UnsplashImageFromApi(query)

			if unsplash.URLs.Small != "" {

				fmt.Println(unsplash.URLs.Small)
				fmt.Println(unsplash.Links.DownloadLocation)

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Something happened, I couldn't find an image matching '" + query + "'! (Small Image URL is Blank)",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					return
				}
			}
		},
		//"math": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		//	ChannelLog(fmt.Sprintf("/// command `math` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
		//	var input string
		//
		//	options := i.ApplicationCommandData().Options
		//	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		//	for _, opt := range options {
		//		optionMap[opt.Name] = opt
		//	}
		//	if option, ok := optionMap["input"]; ok {
		//		// Option values must be type asserted from interface{}.
		//		// Discordgo provides utility functions to make this simple.
		//		input = strings.ToLower(option.StringValue())
		//	}
		//
		//	DoMath(input)
		//
		//	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		//		Type: discordgo.InteractionResponseChannelMessageWithSource,
		//		Data: &discordgo.InteractionResponseData{
		//			Content: "Currently in development, no response will be given",
		//			Flags:   discordgo.MessageFlagsEphemeral,
		//		},
		//	})
		//	if err != nil {
		//		return
		//	}
		//},
		"timestamp": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ChannelLog(fmt.Sprintf("/// command `timestamp` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
			var (
				err     error
				input   string
				display string
			)

			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			if option, ok := optionMap["time"]; ok {
				input = strings.ToUpper(option.StringValue())
			}
			if opt, ok := optionMap["display"]; ok {
				display = opt.StringValue()
			}

			t, err := Str2utc(input)
			if err != nil {
				ChannelLog(fmt.Sprintf("An error occurred while converting input to UTC time: %v", err))
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("An error occurred while converting input to UTC time: %v", err),
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					return
				}
			}
			//snowflake := Utc2snowflake(t)
			//fmt.Println(snowflake)

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "Here is a snowflake timestamp!",
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

		},
		//"ticket": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		//	ChannelLog(fmt.Sprintf("/// command `ticket` used by %v#%v (%v) in server %v", i.Member.User.Username, i.Member.User.Discriminator, i.Member.User.ID, i.GuildID))
		//	var (
		//		err                  error
		//		title                string
		//		msgData              *discordgo.MessageSend = new(discordgo.MessageSend)
		//		dmChannel            *discordgo.Channel     = new(discordgo.Channel)
		//		channelData          discordgo.GuildChannelCreateData
		//		channelOverridesAll  discordgo.PermissionOverwrite
		//		channelOverridesMods discordgo.PermissionOverwrite
		//		//PermissionOverwrites []*discordgo.PermissionOverwrite
		//	)
		//
		//	options := i.ApplicationCommandData().Options
		//	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		//	for _, opt := range options {
		//		optionMap[opt.Name] = opt
		//	}
		//	if option, ok := optionMap["title"]; ok {
		//		title = option.StringValue()
		//		ChannelLog(fmt.Sprintf("Title set to %v", title))
		//	}
		//
		//	msgData.Content = fmt.Sprintf("This is a test response. The title you entered was: %v", title)
		//	dmChannel, err = s.UserChannelCreate(i.Member.User.ID)
		//	if err != nil {
		//		ChannelLog(fmt.Sprintf("Could not create DM channel: %v", err))
		//	}
		//
		//	_, err = s.ChannelMessageSendComplex(dmChannel.ID, msgData)
		//	if err != nil {
		//		ChannelLog(fmt.Sprintf("An error occurred sending DM: %v", err))
		//	}
		//
		//	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		//		Type: discordgo.InteractionResponseChannelMessageWithSource,
		//		Data: &discordgo.InteractionResponseData{
		//			Content: fmt.Sprintf("Check your DMs!"),
		//			Flags:   discordgo.MessageFlagsEphemeral,
		//		},
		//	})
		//
		//	channelOverridesAll.ID = i.GuildID
		//	channelOverridesAll.Type = discordgo.PermissionOverwriteTypeRole
		//	channelOverridesAll.Deny = 0x0000000000000400
		//
		//	channelOverridesMods.ID = "785469547587174421"
		//	channelOverridesMods.Type = discordgo.PermissionOverwriteTypeRole
		//	channelOverridesMods.Allow = 0x0000000000000400
		//
		//	channelData.Name = fmt.Sprintf("%v-%v", title, i.Member.User.Username)
		//	channelData.Type = discordgo.ChannelTypeGuildText
		//	channelData.PermissionOverwrites = []*discordgo.PermissionOverwrite{&channelOverridesAll, &channelOverridesMods}
		//	channelData.ParentID = "1101425457569730590"
		//
		//	s.GuildChannelCreateComplex(i.GuildID, channelData)
		//	if err != nil {
		//		return
		//	}
		//},
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
	s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{{Type: 3, Name: "over the Den // " + buildstring}},
	})

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

func ChannelLog(logInput string) {
	fmt.Println(logInput)

	_, err := s.ChannelMessageSendComplex("1087401958173847552", &discordgo.MessageSend{
		Content: logInput,
	})
	if err != nil {
		fmt.Println(err)
	}
}
