package cmdsServer

import (
	"EviscerateGo/auxiliary/api"
	"EviscerateGo/auxiliary/structs"
	"EviscerateGo/auxiliary/tokens"
	"fmt"
	"image/jpeg"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/bwmarrin/discordgo"
	"github.com/cavaliergopher/grab/v3"
	"github.com/lucasb-eyer/go-colorful"
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
		colorfulCol  colorful.Color
		artistString string
	)

	s := ctx.GetSession()
	i := ctx.GetEvent()

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
			fmt.Sprintf("An error occurred while sending failure response: %v", err)
		}
	} else {
		parsedURL := "https://open.spotify.com/track/" + songdata.Tracks.Items[0].Data.ID
		parsedAlbumURL := "https://open.spotify.com/album/" + songdata.Tracks.Items[0].Data.AlbumOfTrack.ID

		for _, artistName := range songdata.Tracks.Items[0].Data.Artists.Items {
			artistString = artistString + ", [" + artistName.Profile.Name + "](https://open.spotify.com/artist/" + strings.ReplaceAll(artistName.URI, "spotify:artist:", "") + ")"
			fmt.Sprintf(artistString)
		}

		artistString = strings.Replace(artistString, ", ", "", 1)

		resp, err := grab.Get(".", songdata.Tracks.Items[0].Data.AlbumOfTrack.CoverArt.Sources[2].URL)
		fmt.Sprintf("Download saved to %v", resp.Filename)
		file := resp.Filename
		f, err := os.Open(file)
		src, err := jpeg.Decode(f)
		avgcol, err := prominentcolor.Kmeans(src)
		fmt.Sprintf("%v %v %v", avgcol[0].Color.R, avgcol[0].Color.G, avgcol[0].Color.B)

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
			fmt.Sprintf("An error occurred during file deletion: %v", err)
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
