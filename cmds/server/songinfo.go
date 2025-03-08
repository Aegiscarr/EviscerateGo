package cmdsServer

import (
	"EviscerateGo/lib/api"
	"EviscerateGo/lib/color"
	"EviscerateGo/lib/structs"
	"EviscerateGo/lib/tokens"
	"fmt"
	"image/jpeg"
	"os"
	"strings"
	"time"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/bwmarrin/discordgo"
	"github.com/cavaliergopher/grab/v3"
	"github.com/zekrotja/ken"
)

type SongInfoCommand struct {
	_ ken.SlashCommand
	_ ken.DmCapable
}

var (
	_ ken.SlashCommand = (*SongInfoCommand)(nil)
	_ ken.DmCapable    = (*SongInfoCommand)(nil)
)

func (c *SongInfoCommand) Name() string {
	return "songinfo"
}

func (c *SongInfoCommand) Description() string {
	return "grabs songinfo from spotify. might not work as well as you'd hope, spotify's search is just sh*t now"
}

func (c *SongInfoCommand) Version() string {
	return "1.0.0"
}

func (c *SongInfoCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *SongInfoCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "query",
			Description: "song you want to search for. artist - songtitle format works best.",
			Required:    true,
		},
	}
}

func (c *SongInfoCommand) IsDmCapable() bool {
	return true
}

func (c *SongInfoCommand) Run(ctx ken.Context) (err error) {
	var (
		query        string
		songdata     *structs.RapidSzResponse
		sDuration    int64
		t            time.Time
		artistString string
	)

	s := ctx.GetSession()
	i := ctx.GetEvent()

	_ = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "song lookup in progress",
		},
	})

	query = strings.ReplaceAll(ctx.Options().GetByName("query").StringValue(), " ", "%20")

	songdata = api.GetRapidAPICall(query, "tracks", tokens.RapidSzToken)

	if songdata.Tracks.TotalCount == 0 {
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title: "No results.",
				},
			},
		})
		if err != nil {
			fmt.Printf("An error occurred while sending failure response: %v", err)
		}
	} else {
		parsedURL := "https://open.spotify.com/track/" + songdata.Tracks.Items[0].Data.ID
		parsedAlbumURL := "https://open.spotify.com/album/" + songdata.Tracks.Items[0].Data.AlbumOfTrack.ID

		for _, artistName := range songdata.Tracks.Items[0].Data.Artists.Items {
			artistString = artistString + ", [" + artistName.Profile.Name + "](https://open.spotify.com/artist/" + strings.ReplaceAll(artistName.URI, "spotify:artist:", "") + ")"
			fmt.Println(artistString)
		}

		artistString = strings.Replace(artistString, ", ", "", 1)

		resp, _ := grab.Get(".", songdata.Tracks.Items[0].Data.AlbumOfTrack.CoverArt.Sources[2].URL)
		fmt.Printf("Download saved to %v", resp.Filename)
		file := resp.Filename
		f, _ := os.Open(file)
		src, _ := jpeg.Decode(f)
		avgcol, _ := prominentcolor.Kmeans(src)
		fmt.Printf("%v %v %v", avgcol[0].Color.R, avgcol[0].Color.G, avgcol[0].Color.B)

		hexcol := color.ConvertColorInt64(avgcol)

		err = os.Remove(file)
		if err != nil {
			fmt.Printf("An error occurred during file deletion: %v", err)
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
			fmt.Printf("An error occurred sending the embed: %v", err)
		}

	}
	if err != nil {
		return
	}

	return err
}
