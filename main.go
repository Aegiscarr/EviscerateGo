package main

import (
	"EviscerateGo/auxiliary/presence"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/ken/examples/basic/commands"
	"github.com/zekrotja/ken/store"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

var BotToken = flag.String("token", "", "Bot token")
var buildstring = "b20250227-ken"

func main() {

	*BotToken = ReadTokenFromFile("token-dev.txt")
	//*BotToken = ReadTokenFromFile("token.txt")
	if *BotToken != "" {
		log.Println("Token read from file")
	} else {
		log.Println("Token not read from file, fetching from env")
		*BotToken = os.Getenv("TOKEN")
	}

	session, err := discordgo.New("Bot " + *BotToken)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	k, err := ken.New(session, ken.Options{
		CommandStore: store.NewDefault(),
	})
	must(err)

	must(k.RegisterCommands(new(commands.TestCommand)))

	defer k.Unregister()

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	must(session.Open())
	must(presence.SetStatusOnLaunch(session))

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

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
