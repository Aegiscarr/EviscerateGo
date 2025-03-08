package main

import (
	cmdsAdmin "EviscerateGo/cmds/admin"
	cmdsServer "EviscerateGo/cmds/server"
	"EviscerateGo/lib/conf"
	"EviscerateGo/lib/presence"
	"EviscerateGo/lib/tokens"
	"EviscerateGo/lib/txt"

	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/ken/store"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	conf.GetAdminId()
	txt.InitFTReplacers()
	tokens.GetBotToken()
	tokens.GetRapidApiToken()
	tokens.GetUploaderToken()
	tokens.GetUnsplashToken()

	session, err := discordgo.New("Bot " + tokens.BotToken)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	k, err := ken.New(session, ken.Options{
		CommandStore: store.NewDefault(),
	})
	must(err)

	must(k.RegisterCommands(
		new(cmdsServer.PingCommand),
		new(cmdsServer.DevExcuse),
		new(cmdsServer.PkeCommand),
		new(cmdsServer.RandomColor),
		new(cmdsServer.ColorCommand),
		new(cmdsServer.D20Command),
		new(cmdsServer.EchoCommand),
		new(cmdsServer.EightBallCommand),
		new(cmdsServer.SongInfoCommand),
		new(cmdsServer.GoogleCommand),
		new(cmdsServer.MagickCommand),
		new(cmdsServer.UnsplashCommand),
		new(cmdsServer.FancyTextCommand),
		new(cmdsServer.TimestampCommand),
		new(cmdsServer.UserInfoCommand),

		new(cmdsAdmin.StatusSetCommand),
	),
	)

	defer k.Unregister()

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	must(session.Open())
	must(presence.SetStatusOnLaunch(session))

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}
