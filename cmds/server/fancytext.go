package cmdsServer

import (
	"EviscerateGo/lib/txt"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type FancyTextCommand struct {
	_ ken.SlashCommand
	_ ken.DmCapable
}

var (
	_ ken.SlashCommand = (*FancyTextCommand)(nil)
	_ ken.DmCapable    = (*FancyTextCommand)(nil)
)

func (c *FancyTextCommand) Name() string {
	return "fancytext"
}

func (c *FancyTextCommand) Description() string {
	return "makes your text adequately 𝒻𝒶𝓃𝒸𝓎, if you like to pretend youre better than the rest of us."
}

func (c *FancyTextCommand) Version() string {
	return "1.0.0"
}

func (c *FancyTextCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *FancyTextCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "text",
			Description: "text to fancify",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "fontface",
			Description: "unicode pseudofont to use",
			Required:    false,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "ⓑⓤⓑⓑⓛⓔⓢ (bubbles)",
					Value: "1|bubbles",
				},
				{
					Name:  "🅑🅤🅑🅑🅛🅔🅑🅛🅐🅒🅚 (bubbleblack)",
					Value: "2|bubbleblack",
				},
				{
					Name:  "⒫⒜⒭⒠⒩⒯⒣⒠⒮⒤⒮ (parenthesis)",
					Value: "3|parenthesis",
				},
				{
					Name:  "ˢᵘᵖᵉʳˢᶜʳⁱᵖᵗ (superscript)",
					Value: "4|superscript",
				},
				{
					Name:  "ｆｕｌｌｗｉｄｔｈ (fullwidth)",
					Value: "5|fullwidth",
				},
				{
					Name:  "𝐌𝐚𝐭𝐡𝐁𝐨𝐥𝐝 (mathbold)",
					Value: "6|mathbold",
				},
				{
					Name:  "𝑴𝒂𝒕𝒉𝑰𝒕𝒂𝒍𝒊𝒄 (mathitalic)",
					Value: "7|mathitalic",
				},
				{
					Name:  "𝖬𝖺𝗍𝗁𝖲𝖺𝗇𝗌 (mathsans)",
					Value: "8|mathsans",
				},
				{
					Name:  "𝘔𝘢𝘵𝘩𝘚𝘢𝘯𝘴𝘐𝘵𝘢𝘭𝘪𝘤 (mathsansitalic)",
					Value: "9|mathsansital",
				},
				{
					Name:  "𝗠𝗮𝘁𝗵𝗦𝗮𝗻𝘀𝗕𝗼𝗹𝗱 (mathsansbold)",
					Value: "10|mathsansbold",
				},
				{
					Name:  "𝙈𝙖𝙩𝙝𝙎𝙖𝙣𝙨𝘽𝙤𝙡𝙙𝙄𝙩𝙖𝙡𝙞𝙘 (mathsansbolditalic)",
					Value: "11|mathsansboldital",
				},
				{
					Name:  "𝔉𝔯𝔞𝔨𝔱𝔲𝔯 (fraktur)",
					Value: "12|fraktur",
				},
				{
					Name:  "𝕱𝖗𝖆𝖐𝖙𝖚𝖗𝕭𝖔𝖑𝖉 (frakturbold)",
					Value: "13|frakturbold",
				},
				{
					Name:  "Яц$$їап (russian)",
					Value: "14|russian",
				},
				{
					Name:  " ﾌ卂卩卂几乇丂乇 (japanese)",
					Value: "15|japanese",
				},
				{
					Name:  "คгค๒เς (arabic)",
					Value: "16|arabic",
				},
				{
					Name:  "ᎦᏗᎥᏒᎩ (fairy)",
					Value: "17|fairy",
				},
				{
					Name:  "աɨʐǟʀɖ (wizard)",
					Value: "18|wizard",
				},
				{
					Name:  "𝙼𝚘𝚗𝚘𝚜𝚙𝚊𝚌𝚎 (monospace)",
					Value: "19|monospace",
				},
				{
					Name:  "𝒮𝒸𝓇𝒾𝓅𝓉 (script)",
					Value: "20|script",
				},
				{
					Name:  "𝓢𝓬𝓻𝓲𝓹𝓽𝓑𝓸𝓵𝓭 (scriptbold)",
					Value: "21|scriptbold",
				},
				{
					Name:  "𝔻𝕠𝕦𝕓𝕝𝕖𝕊𝕥𝕣𝕦𝕔𝕜 (doublestruck)",
					Value: "22|doublestruck",
				},
				{
					Name:  "🅂🅀🅄🄰🅁🄴🄳 (squared)",
					Value: "23|squared",
				},
				{
					Name:  "ƒυηку (funky)",
					Value: "24|funky",
				},
				{
					Name:  "Áćúté (acute)",
					Value: "25|acute",
				},
				//{
				//	Name:  "ṚöċḳḊöẗṡ (rockdots)",
				//	Value: "26|rockdots",
				//},
				//{
				//	Name:  "Sŧɍøꝁɇđ (stroked)",
				//	Value: "27|stroked",
				//},
				//{
				//	Name:  "Iuʌǝɹʇǝp (inverted)",
				//	Value: "28|inverted",
				//},
				//{
				//	Name:  "1337 [3><7?3/V\\3] (1337extreme)",
				//	Value: "29|1337extreme",
				//},
				//{
				//	Name:  "Ｈｅａｖｙ (heavy)",
				//	Value: "30|heavy",
				//},
				//{
				//	Name:  "Lιƚƚʅҽ Fαɳƈყ (littlefancy)",
				//	Value: "31|littlefancy",
				//},
				//{
				//	Name:  "ʄąცƖɛ (fable)",
				//	Value: "32|fable",
				//},
				//{
				//	Name:  "ŞຟirlŞ (swirls)",
				//	Value: "33|swirls",
				//},
				//{
				//	Name:  "Ä¢¢êñ† (accent)",
				//	Value: "34|accent",
				//},
				//{
				//	Name:  "ᄂIПΣΛЯ (linear)",
				//	Value: "35|linear",
				//},
				//{
				//	Name:  "₴₵Ɽł฿฿ⱠɆ₴ (scribbles)",
				//	Value: "36|scribbles",
				//},
				//{
				//	Name:  "ﾌﾑｱﾑ刀乇丂乇 丂ᄃ尺ﾉｱｲ (japanesescript)",
				//	Value: "37|jpscript",
				//},
				//{
				//	Name:  "【S】【o】【l】【i】【t】【u】【d】【e】(solitude)",
				//	Value: "38|solitude",
				//},
				//{
				//	Name:  "『B』『r』『a』『c』『k』『e』『t』『s』(brackets)",
				//	Value: "39|brackets",
				//},
				//{
				//	Name:  "[̲̅B][̲̅o][̲̅x] [̲̅L][̲̅i][̲̅n][̲̅e][̲̅s] (boxlines)",
				//	Value: "40|boxlines",
				//},
				//{
				//	Name:  "ϚվʍҍօӀìç (symbolic)",
				//	Value: "41|symbolic",
				//},
				//{
				//	Name:  "ᗷᘿᘉᖶ (bent)",
				//	Value: "42|bent",
				//},
				//{
				//	Name:  "D̶̶a̶s̶h̶e̶s̶  (dashes)",
				//	Value: "43|dashes",
				//},
				//{
				//	Name:  "S̴i̴d̴e̴S̴q̴u̴i̴g̴g̴l̴e̴s̴  (sidesquiggles)",
				//	Value: "44|sidesquiggles",
				//},
				//{
				//	Name:  "S̷i̷d̷e̷S̷l̷a̷s̷h̷e̷s̷  (sideslashes)",
				//	Value: "45|sideslashes",
				//},
				//{
				//	Name:  "D̳o̳u̳b̳l̳e̳U̳n̳d̳e̳r̳l̳i̳n̳e̳ (doubleunderline)",
				//	Value: "46|doubleunderline",
				//},
				//{
				//	Name:  "T̾o̾p̾S̾q̾u̾i̾g̾g̾l̾e̾s̾ (topsquiggles)",
				//	Value: "47|topsquiggles",
				//},
				//{
				//	Name:  "A͎r͎r͎o͎w͎U͎p͎ (arrowup)",
				//	Value: "48|arrowup",
				//},
				//{
				//	Name:  "E͓̽x͓̽e͓̽s͓̽ (exes)",
				//	Value: "49|exes",
				//},
			},
		},
	}
}

func (c *FancyTextCommand) IsDmCapable() bool {
	return true
}

func (c *FancyTextCommand) Run(ctx ken.Context) (err error) {

	var (
		text     string
		fontface string
	)

	text = ctx.Options().GetByName("text").StringValue()

	if v, ok := ctx.Options().GetByNameOptional("fontface"); ok {
		fontface = v.StringValue()
	}

	if fontface == "" { // for the love of god make this less terrible
		fontface = strconv.FormatInt(int64(rand.Intn(25)), 10) + "| "
		fmt.Println(fontface)
	}
	fancifiedText := txt.FancyTextReplace(text, fontface)

	_ = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       fancifiedText,
					Color:       0x8c1bb1,
					Description: `it's dangerous to try be stylish, take this and watch it fall flat, probably.`,
				},
			},
		},
	})

	return err
}
