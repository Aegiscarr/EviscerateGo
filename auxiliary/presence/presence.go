package presence

import "github.com/bwmarrin/discordgo"

var buildstring string = "b250227-ken"
var err error

func SetStatusOnLaunch(s *discordgo.Session) error {
	err = s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{{Type: 3, Name: "over the Den // " + buildstring}},
	})
	return err
}
