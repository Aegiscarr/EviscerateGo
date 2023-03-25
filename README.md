# Eviscerate The Synth
*i'll get a logo for this dingus someday, i swear*\
Discord bot I'm developing as a long term "get good at problem solving and golang" project

#### *Setup*
##### REQUIREMENTS (if you want to have feature parity anyway)
- [Go 1.19.2](https://go.dev/dl) or higher
- a 64 bit computer, which you probably do if youve bought a computer in the last decade
- a [discord api token](https://discord.com/developers)
- a rapidapi token + being subscribed to the [unofficial spotify api](https://rapidapi.com/Glavier/api/spotify23)
- a [cumulonimbus](https://alekeagle.me) api token

store your api tokens in:
- token.txt (discord token)
- rapid-sz-token.txt (rapidapi token)
- uploader-token.txt (cumulonimbus token)

at the bottom of main.go, there's a function called ChannelLog which needs a channel ID your bot will have access to. I'll set up a better way of doing it soon.

```
go get
go build
```
run executable\
profit

#### *Will there be a public version of this bot?*
Probably not, at least not for the forseeable future. This is currently a passion project and not something I intend for to become the next big "thing".

#### *Can I watch you develop it, at least?*
Yes! I'll usually be in my [Discord server](https://discord.gg/SJcAWEynbj) streaming it, and sometimes on [Twitch](https://twitch.tv/aegiscarr).\
My Discord has a channel dedicated to me ranting about it too.

#### *This wouldn't have been the same without:*
- [Deepfried Chips](https://github.com/Deepfried-Chips) for helping me get this thing off the ground in the first place
- [My entire community](https://discord.gg/SJcAWEynbj) for helping me with ideas and testing (yall have found so many little things i wouldn't have lol)

#### *Libraries used:*
- [Colorful](https://github.com/lucasb-eyer/go-colorful) by [Lucas Beyer](https://github.com/lucasb-eyer)
- [ProminentColor](https://github.com/EdlinOrg/prominentcolor) by [EdlinOrg](https://github.com/EdlinOrg)
- [DiscordGo](https://github.com/bwmarrin/discordgo) by [bwmarrin](https://github.com/bwmarrin)
- [GRAB](https://github.com/cavaliergopher/grab) by [Ryan Armstrong](https://github.com/cavaliercoder)
- [htgo-tts](https://github.com/hegedustibor/htgo-tts) by [Tibor Hegedűs](https://github.com/hegedustibor)
- [A DCA implementation](https://github.com/jonas747/dca) by [jonas747](https://github.com/jonas747)
- [GO IMAGE FILTERING TOOLKIT](https://github.com/disintegration/gift) by [Grigory Dryapak](https://github.com/disintegration)
##### *Indirectly used:*
- [An MP3 implementation](https://github.com/hajimehoshi/go-mp3) by [Hajime Hoshi](https://github.com/hajimehoshi) // used in htgo-tts
- [OTO](https://github.com/hajimehoshi/oto) by [Hajime Hoshi](https://github.com/hajimehoshi) // also, used in htgo-tts
- [An OGG implementation](https://github.com/jonas747/ogg) by [jonas747](https://github.com/jonas747) // used in dca
- [Resize](https://github.com/nfnt/resize) by [Jan Schlicht](https://github.com/nfnt) // used in prominentcolor
- [Cutter](https://github.com/oliamb/cutter) by [Olivier Amblet](https://github.com/oliamb) // used in prominentcolor
###
###
*damn, aegis actually writing a readme for once? yeah im freaked out too*

