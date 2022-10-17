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

	"github.com/Aegiscarr/randomcolor"
	"github.com/bwmarrin/discordgo"
)

// BotToken Flags
var (
	BotToken = flag.String("token", "", "Bot token")
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

//var dicerolls = [`Luck really isn't on your side today, huh? It's a 1.`, `A 2. Couldn't have been much worse.`, `It's a tree!- Oh, wait. A 3.`, `A four-se. Of course. That didn't work, did it?`, `A 5. Nothing funny here.`, `A 6. The devil, anyone?`, `Lucky number 7! Now can you get two more?`, `8, not bad.`, `Just under halfway up. A 9`, `A 10! Halfway up the scale!`, `11. Decent.`, `12. Could have been much worse. Could've also been better, though.`, `13. Feelin' lucky?`, `Aand it's come up 14!`, `15! Getting up there!`, `16, solid.`, `17. Rolling real high now, aren't you?`, `18! You're old eno- wait this isn't a birthday.`, `19! So CLOSE!`, `NAT 20 BAYBEE!`]

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

func GetDevExcuse() *RandomDevExcuse {
	resp, err := http.Get(fmt.Sprintf("https://api.tabliss.io/v1/developer-excuses"))
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

			var RandomColorHex string = randomcolor.GetRandomColorInHex()
			var RandomColorRGB randomcolor.RGBColor = randomcolor.GetRandomColorInRgb()
			var red string = strconv.Itoa(RandomColorRGB.Red)
			var green string = strconv.Itoa(RandomColorRGB.Green)
			var blue string = strconv.Itoa(RandomColorRGB.Blue)
			var RandomColorHSV randomcolor.HSVColor = randomcolor.RGBToHSV(RandomColorRGB)
			var hue string = strconv.FormatFloat(RandomColorHSV.Hue, 'f', 0, 64)
			var sat string = strconv.FormatFloat(RandomColorHSV.Saturation, 'f', 0, 64)
			var val string = strconv.FormatFloat(RandomColorHSV.Value, 'f', 0, 64)
			var RandomColorHexInt64, res = strconv.ParseInt(RandomColorHex, 16, 32)
			var err error

			log.Println("Red: " + red)
			log.Println("Green: " + green)
			log.Println("Blue: " + blue)
			log.Println("Hue: " + hue)
			log.Println("Sat: " + sat)
			log.Println("Val: " + val)
			log.Println("---")

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "Random color",
							Color: int(RandomColorHexInt64),
							Thumbnail: &discordgo.MessageEmbedThumbnail{
								URL:    "https://singlecolorimage.com/get/" + RandomColorHex + "/100x100",
								Width:  100,
								Height: 100,
							},
							Fields: []*discordgo.MessageEmbedField{
								&discordgo.MessageEmbedField{
									Name:  "Hex",
									Value: "#" + RandomColorHex,
								},
								&discordgo.MessageEmbedField{
									Name:  "RGB",
									Value: "[" + red + ", " + green + ", " + blue + "]",
								},
								&discordgo.MessageEmbedField{
									Name:  "HSV",
									Value: hue + "°, " + sat + "%, " + val + "%",
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

//func getRandomColor() {
//	var randHex string = randomcolor.GetRandomColorInHex()
//	var randRGB randomcolor.RGBColor = randomcolor.GetRandomColorInRgb()
//	var randHSV randomcolor.HSVColor = randomcolor.GetRandomColorInHSV()
//}
