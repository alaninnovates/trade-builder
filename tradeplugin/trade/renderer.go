package trade

import (
	"alaninnovates.com/trade-builder/common"
	"alaninnovates.com/trade-builder/common/loaders"
	"io"
	"math"
	"strconv"

	"github.com/fogleman/gg"
)

/*
Center an element, returning the left/top coordinate (i.e. anchor to right)

fullSize: size of the screen (i.e. w/h)
desiredWidth: desired width of the element

example:
centerAnchorL(gg.Width(), 1000)
-> centers an element of size 1000, returning the left coordinate
*/
func centerAnchorL(fullSize, desiredWidth int) float64 {
	center := fullSize / 2
	return float64(center - (desiredWidth / 2))
}

func RenderTrade(t *Trade) *io.PipeReader {
	var offeringStickers, lfStickers []string
	var offeringQuantities, lfQuantities []int
	// todo: support beequips
	for _, stickerRaw := range t.GetLookingFor() {
		sticker := stickerRaw.(Sticker)
		lfStickers = append(lfStickers, sticker.Name)
		lfQuantities = append(lfQuantities, sticker.Quantity)
	}
	for _, stickerRaw := range t.GetOffering() {
		sticker := stickerRaw.(Sticker)
		offeringStickers = append(offeringStickers, sticker.Name)
		offeringQuantities = append(offeringQuantities, sticker.Quantity)
	}

	offeringRowCnt := int(math.Ceil(float64(len(offeringStickers)) / float64(4)))
	lfRowCnt := int(math.Ceil(float64(len(lfStickers)) / float64(4)))
	//fmt.Println(offeringRowCnt, lfRowCnt)

	offeringHeight := float64(max(672, 20*offeringRowCnt+148*(offeringRowCnt-1)+128+20))
	lfHeight := float64(max(672, 20*lfRowCnt+148*(lfRowCnt-1)+128+20))
	//fmt.Println(offeringHeight, lfHeight)

	addtlHeight := -672 - 168 + int(max(offeringHeight, lfHeight))
	dc := gg.NewContext(2000, 1200+addtlHeight)

	// background
	bg := loaders.GetRandomBackground()
	// keep drawing bg until we cover the whole image
	y := 0
	for y < 1200+addtlHeight {
		dc.DrawImage(bg, 0, y)
		y += bg.Bounds().Dy()
	}

	// big rectangle box
	dc.SetHexColor("#FEC200")
	dc.DrawRoundedRectangle(centerAnchorL(dc.Width(), 1800), 102, 1800, 1000+float64(addtlHeight), 36)
	dc.Fill()

	// header: offering/lf
	tradeHeader, err := gg.LoadPNG("assets/trade_header.png")
	if err != nil {
		panic(err)
	}
	dc.DrawImageAnchored(tradeHeader, dc.Width()/2, 17, 0.5, 0)

	// center: arrows
	tradeCenter, err := gg.LoadPNG("assets/trade_center.png")
	if err != nil {
		panic(err)
	}
	dc.DrawImageAnchored(tradeCenter, dc.Width()/2, 260, 0.5, 0)

	// item boxes: offering
	dc.SetHexColor("#FEBC2B")
	dc.DrawRoundedRectangle(191, 215, 672, offeringHeight, 24)
	dc.Fill()

	dc.SetLineWidth(4)
	dc.SetHexColor("#000000")
	dc.DrawRoundedRectangle(191, 215, 672, offeringHeight, 24)
	dc.Stroke()

	// item boxes: lf
	dc.SetHexColor("#FEBC2B")
	dc.DrawRoundedRectangle(1140, 215, 672, lfHeight, 24)
	dc.Fill()

	dc.SetLineWidth(4)
	dc.SetHexColor("#000000")
	dc.DrawRoundedRectangle(1140, 215, 672, lfHeight, 24)
	dc.Stroke()

	if err := dc.LoadFontFace("assets/Buycat.ttf", 28); err != nil {
		panic(err)
	}

	drawStickers(dc, 0, offeringRowCnt, offeringStickers, offeringQuantities)
	drawStickers(dc, 949, lfRowCnt, lfStickers, lfQuantities)

	return common.ImageToPipe(dc.Image())
}

func drawStickers(dc *gg.Context, offset int, rowCnt int, stickers []string, quantities []int) {
	padding := 20
	/*
		Each sticker has 168 of space: 20 + 128 + 20
	*/
	idx := 0
	for i := 1; i <= rowCnt; i++ {
		for j := 1; j <= 4; j++ {
			if idx >= len(stickers) {
				break
			}
			img := loaders.GetStickerImage(stickers[idx])
			posX := 191 + offset + padding*j + 148*(j-1)
			posY := 215 + padding*i + 148*(i-1)
			// draw border
			dc.SetHexColor("#00000055")
			dc.SetLineWidth(4)
			dc.DrawRoundedRectangle(float64(posX), float64(posY), 128, 128, 12)
			dc.Stroke()
			// draw sticker image
			dc.DrawImageAnchored(img, int(posX)+64, int(posY)+64, 0.5, 0.5)
			// draw qty
			// if quantities[idx] > 1 {
			dc.SetHexColor("#000000")
			qtyText := "x" + strconv.Itoa(quantities[idx])
			measuredX, _ := dc.MeasureString(qtyText)
			offset := float64(8)
			dc.DrawString(qtyText, float64(posX)+128-measuredX-offset, float64(posY)+128-offset)
			// }
			idx++
		}
	}
}
