package presence

import "github.com/bwmarrin/discordgo"

//var buildstring string = "b250308-ken"
var err error

func SetStatusOnLaunch(s *discordgo.Session) error {
	err = s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{{Type: 3, Name: "over the Den"}},
	})
	return err
}

func SetStatus(s *discordgo.Session) error {

	return err
}
