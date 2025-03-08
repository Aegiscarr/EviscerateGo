package txt

import (
	"fmt"
	"strconv"
	"strings"
)

func FancyTextReplace(text string, fontface string) string {
	var (
		fancifiedText string
		fontID        string
		fontName      string
		foundSep      bool
	)

	fmt.Println(text)
	fontID, fontName, foundSep = strings.Cut(fontface, "|")
	if !foundSep {
		fmt.Println("Separator not found")
	} else {
		fmt.Printf("Font ID: %v // Font Name: %v", fontID, fontName)
	}

	ID, _ := strconv.ParseInt(fontID, 10, 0)

	switch {
	case ID == 1:
		{
			fancifiedText = bubblesReplacer.Replace(text)
		}
	case ID == 2:
		{
			fancifiedText = bubblesBlackReplacer.Replace(text)
		}
	case ID == 3:
		{
			fancifiedText = parenthesisReplacer.Replace(text)
		}
	case ID == 4:
		{
			fancifiedText = superscriptReplacer.Replace(text)
		}
	case ID == 5:
		{
			fancifiedText = fullWidthReplacer.Replace(text)
		}
	case ID == 6:
		{
			fancifiedText = mathBoldReplacer.Replace(text)
		}
	case ID == 7:
		{
			fancifiedText = mathBoldItalicReplacer.Replace(text)
		}
	case ID == 8:
		{
			fancifiedText = mathSansReplacer.Replace(text)
		}
	case ID == 9:
		{
			fancifiedText = mathSansItalicReplacer.Replace(text)
		}
	case ID == 10:
		{
			fancifiedText = mathSansBoldReplacer.Replace(text)
		}
	case ID == 11:
		{
			fancifiedText = mathSansBoldItalicReplacer.Replace(text)
		}
	case ID == 12:
		{
			fancifiedText = frakturReplacer.Replace(text)
		}
	case ID == 13:
		{
			fancifiedText = frakturBoldReplacer.Replace(text)
		}
	case ID == 14:
		{
			fancifiedText = russianReplacer.Replace(text)
		}
	case ID == 15:
		{
			fancifiedText = japaneseReplacer.Replace(text)
		}
	case ID == 16:
		{
			fancifiedText = arabicReplacer.Replace(text)
		}
	case ID == 17:
		{
			fancifiedText = fairyReplacer.Replace(text)
		}
	case ID == 18:
		{
			fancifiedText = wizardReplacer.Replace(text)
		}
	case ID == 19:
		{
			fancifiedText = monospaceReplacer.Replace(text)
		}
	case ID == 20:
		{
			fancifiedText = scriptReplacer.Replace(text)
		}
	case ID == 21:
		{
			fancifiedText = scriptBoldReplacer.Replace(text)
		}
	case ID == 22:
		{
			fancifiedText = doubleStruckReplacer.Replace(text)
		}
	case ID == 23:
		{
			fancifiedText = squaredReplacer.Replace(text)
		}
	case ID == 24:
		{
			fancifiedText = funkyReplacer.Replace(text)
		}
	case ID == 25:
		{
			fancifiedText = acuteReplacer.Replace(text)
		}
	}

	return fancifiedText
}
