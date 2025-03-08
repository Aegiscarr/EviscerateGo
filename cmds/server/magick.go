package cmdsServer

import (
	"EviscerateGo/lib/structs"
	"EviscerateGo/lib/tokens"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cavaliergopher/grab/v3"
	"github.com/disintegration/gift"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/zekrotja/ken"
	"golang.org/x/image/webp"
)

type MagickCommand struct {
	_ ken.SlashCommand
}

var (
	_ ken.SlashCommand = (*MagickCommand)(nil)
)

func (c *MagickCommand) Name() string {
	return "magick"
}

func (c *MagickCommand) Description() string {
	return "mess with an image through imagemagick"
}

func (c *MagickCommand) Version() string {
	return "1.0.0"
}

func (c *MagickCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "url",
			Description: "the image you want to mess with",
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
	}
}

func (c *MagickCommand) Run(ctx ken.Context) (err error) {

	var (
		imgUrl         string
		method         string
		dst            *image.RGBA
		src            image.Image
		fEdit          *os.File
		execFolder     string
		execPath       string
		imagePath      string
		uploadResponse structs.CumulonimbusResponse
		colorfulCol    colorful.Color
		intensity      float32
	)

	s := ctx.GetSession()
	i := ctx.GetEvent()

	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "image editing in process... please stand by.",
		},
	})

	imgUrl = ctx.Options().GetByName("url").StringValue()
	method = ctx.Options().GetByName("method").StringValue()

	if v, ok := ctx.Options().GetByNameOptional("intensity"); ok {
		intensity = float32(v.FloatValue())
	}

	execPath, err = os.Executable()
	execFolder = filepath.Dir(execPath)

	resp, err := grab.Get(".", imgUrl)
	if err != nil {
		fmt.Print(err)
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title: fmt.Sprintf("Error: no image found at URL %v", imgUrl),
					Color: 0xdd0000,
				},
			},
		})

		if err != nil {
			//ChannelLog(fmt.Sprintf("an error occurred while sending embed update: %v", err))
		}
	}

	//ChannelLog(fmt.Sprintf("Download saved to %v ", resp.Filename))
	file := resp.Filename
	buf := make([]byte, 512)

	openFile, err := os.Open(file)
	if err != nil {
		//ChannelLog(fmt.Sprintf("An error occurred during file read: %v", err))
	}
	_, err = openFile.Read(buf)
	//ChannelLog(fmt.Sprintf("bytes read: `%v`", n))
	if err != nil {
		//ChannelLog(fmt.Sprintf("An error occurred during file read to buffer: %v", err))
	}
	contentType := http.DetectContentType(buf)
	//ChannelLog(fmt.Sprintf("Content type is `%v`", contentType))
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
		//ChannelLog(fmt.Sprintf("An error occurred during decode: %v", err))
		return
	}

	if method == "blur" {
		if intensity == 0 {
			intensity = 1.5
		}

		if intensity < 0 {
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title: "intensity should be higher than 0",
						Color: 0xdd0000,
					},
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

	//avgcol, err := prominentcolor.Kmeans(dst)
	//fmt.Printf("%v %v %v", avgcol[0].Color.R, avgcol[0].Color.G, avgcol[0].Color.B)

	// COLOR CONVERSION YAAAAAAAAAAAAAAAAAAAAAAAAA
	//	avgColRInt, err := strconv.ParseInt(strconv.FormatUint(uint64(avgcol[0].Color.R), 10), 10, 64)
	//	avgColR, err := strconv.ParseFloat(strconv.Itoa(int(avgColRInt)), 64)
	//
	//	avgColGInt, err := strconv.ParseInt(strconv.FormatUint(uint64(avgcol[0].Color.G), 10), 10, 64)
	//	avgColG, err := strconv.ParseFloat(strconv.Itoa(int(avgColGInt)), 64)
	//
	//	avgColBInt, err := strconv.ParseInt(strconv.FormatUint(uint64(avgcol[0].Color.B), 10), 10, 64)
	//	avgColB, err := strconv.ParseFloat(strconv.Itoa(int(avgColBInt)), 64)
	//
	//	avgColRDiv := avgColR / 255
	//	avgColGDiv := avgColG / 255
	//	avgColBDiv := avgColB / 255
	//
	//	colorfulCol.R = avgColRDiv
	//	colorfulCol.G = avgColGDiv
	//	colorfulCol.B = avgColBDiv

	fEdit, err = os.Create("image.png")
	if err != nil {
		//ChannelLog(fmt.Sprintf("error while creating image file: %v", err))
	}
	png.Encode(fEdit, dst)

	hexcol, err := strconv.ParseInt(strings.ReplaceAll(colorfulCol.Hex(), "#", ""), 16, 64)

	imagePath = execFolder + "\\image.png"
	//ChannelLog(fmt.Sprintf(imagePath))

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
	req.Header.Set("Authorization", tokens.UploaderToken)
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

	fmt.Println(hexcol)
	fmt.Println(uploadResponse.URL)
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title: "Here's the edited image!",
				Color: 0x000000, //int(hexcol),
				Image: &discordgo.MessageEmbedImage{
					URL: uploadResponse.URL,
				},
				Description: "Upload powered by [Cumulonimbus](https://alekeagle.me)",
			},
		},
	})

	if err != nil {
		fmt.Printf("an error occurred while sending embed update: %v", err)
	}

	os.Remove(file)
	os.Remove(imagePath)
	if err != nil {
		return
	}

	return err
}
