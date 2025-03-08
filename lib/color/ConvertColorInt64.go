package color

import (
	"strconv"
	"strings"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/lucasb-eyer/go-colorful"
)

var colorfulCol colorful.Color

func ConvertColorInt64(colInput []prominentcolor.ColorItem) (hexcol int64) {
	avgColRInt, _ := strconv.ParseInt(strconv.FormatUint(uint64(colInput[0].Color.R), 10), 10, 64)
	avgColR, _ := strconv.ParseFloat(strconv.Itoa(int(avgColRInt)), 64)

	avgColGInt, _ := strconv.ParseInt(strconv.FormatUint(uint64(colInput[0].Color.G), 10), 10, 64)
	avgColG, _ := strconv.ParseFloat(strconv.Itoa(int(avgColGInt)), 64)

	avgColBInt, _ := strconv.ParseInt(strconv.FormatUint(uint64(colInput[0].Color.B), 10), 10, 64)
	avgColB, _ := strconv.ParseFloat(strconv.Itoa(int(avgColBInt)), 64)

	avgColRDiv := avgColR / 255
	avgColGDiv := avgColG / 255
	avgColBDiv := avgColB / 255

	colorfulCol.R = avgColRDiv
	colorfulCol.G = avgColGDiv
	colorfulCol.B = avgColBDiv

	hexcol, _ = strconv.ParseInt(strings.ReplaceAll(colorfulCol.Hex(), "#", ""), 16, 64)

	return hexcol
}
