package main

import (
	"EviscerateGo/auxiliary/presence"
	tokens "EviscerateGo/auxiliary/tokens"
	cmdsServer "EviscerateGo/cmds/server"

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

	tokens.GetBotToken()
	//fmt.Println(tokens.BotToken)
	tokens.GetRapidApiToken()
	//fmt.Println(tokens.RapidSzToken)

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
		new(cmdsServer.TestCommand),
		new(cmdsServer.PingCommand),
		new(cmdsServer.DevExcuse),
		new(cmdsServer.PkeCommand),
		new(cmdsServer.RandomColor),
		new(cmdsServer.ColorCommand),
		new(cmdsServer.D20Command),
		new(cmdsServer.EchoCommand),
		new(cmdsServer.EightBallCommand),
		new(cmdsServer.SongInfoCommand),
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
