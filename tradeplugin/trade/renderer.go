package trade

import (
	"alaninnovates.com/trade-builder/common"
	"alaninnovates.com/trade-builder/common/loaders"
	"io"
	"math"
	"strconv"
	"strings"

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
	var offeringBeequips, lfBeequips []Beequip
	for _, itemRaw := range t.GetLookingFor() {
		switch item := itemRaw.(type) {
		case Sticker:
			lfStickers = append(lfStickers, item.Name)
			lfQuantities = append(lfQuantities, item.Quantity)
		case Beequip:
			lfBeequips = append(lfBeequips, item)
		}
	}
	for _, itemRaw := range t.GetOffering() {
		switch item := itemRaw.(type) {
		case Sticker:
			offeringStickers = append(offeringStickers, item.Name)
			offeringQuantities = append(offeringQuantities, item.Quantity)
		case Beequip:
			offeringBeequips = append(offeringBeequips, item)
		}
	}

	offeringRowCnt := int(math.Ceil(float64(len(offeringStickers)) / float64(4)))
	lfRowCnt := int(math.Ceil(float64(len(lfStickers)) / float64(4)))
	//fmt.Println(offeringRowCnt, lfRowCnt)

	offeringBeequipRowCnt := int(math.Ceil(float64(len(offeringBeequips)) / float64(2)))
	lfBeequipRowCnt := int(math.Ceil(float64(len(lfBeequips)) / float64(2)))
	//fmt.Println(offeringBeequipRowCnt, lfBeequipRowCnt)

	offeringHeight := float64(max(672, 20*offeringRowCnt+148*(offeringRowCnt-1)+128+20+(20*offeringBeequipRowCnt+520*offeringBeequipRowCnt)))
	lfHeight := float64(max(672, 20*lfRowCnt+148*(lfRowCnt-1)+128+20+(20*lfBeequipRowCnt+520*lfBeequipRowCnt)))
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

	lastYOffering := drawStickers(dc, 0, offeringRowCnt, offeringStickers, offeringQuantities)
	lastYLf := drawStickers(dc, 949, lfRowCnt, lfStickers, lfQuantities)
	drawBeequips(dc, 0, lastYOffering+20, offeringBeequipRowCnt, offeringBeequips)
	drawBeequips(dc, 949, lastYLf+20, lfBeequipRowCnt, lfBeequips)

	return common.ImageToPipe(dc.Image())
}

func drawStickers(dc *gg.Context, offset int, rowCnt int, stickers []string, quantities []int) int {
	padding := 20
	/*
		Each sticker has 168 of space: 20 + 128 + 20
	*/
	lastY := 0
	idx := 0
	for i := 1; i <= rowCnt; i++ {
		for j := 1; j <= 4; j++ {
			if idx >= len(stickers) {
				break
			}
			img := loaders.GetStickerImage(stickers[idx])
			posX := 191 + offset + padding*j + 148*(j-1)
			posY := 215 + padding*i + 148*(i-1)
			lastY = posY
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
	if lastY == 0 {
		lastY += 60
	}
	return lastY + 128
}

func drawBeequips(dc *gg.Context, offsetX int, offsetY, rowCnt int, beequips []Beequip) {
	padding := 20
	/*
		Each beequip has 336 of space: 20(pad) + 128+20+20+128 + 20(pad)
	*/
	idx := 0
	for i := 1; i <= rowCnt; i++ {
		for j := 1; j <= 2; j++ {
			if idx >= len(beequips) {
				break
			}
			img := loaders.GetBeequipImage(beequips[idx].Name)
			posX := 191 + offsetX + padding*j + 148*(j-1)
			posY := offsetY + padding*i + 148*(i-1)
			// draw border
			dc.SetHexColor("#00000055")
			dc.SetLineWidth(4)
			dc.DrawRoundedRectangle(float64(posX), float64(posY), 296, 500, 12)
			dc.Stroke()
			// draw beequip image
			dc.DrawImageAnchored(img, int(posX)+148, posY+padding, 0.5, 0)
			currY := posY + padding*2 + img.Bounds().Dy()
			filledStarImg, err := gg.LoadPNG("assets/star_filled.png")
			if err != nil {
				panic(err)
			}
			emptyStarImg, err := gg.LoadPNG("assets/star_empty.png")
			if err != nil {
				panic(err)
			}
			for k := 0; k < 5; k++ {
				if k < beequips[idx].Potential {
					dc.DrawImageAnchored(filledStarImg, posX+90+k*32, currY-20, 0.5, 0)
				} else {
					dc.DrawImageAnchored(emptyStarImg, posX+90+k*32, currY-20, 0.5, 0)
				}
			}
			currY += 40
			// info
			currY += drawTextSet(dc, posX+padding, currY, convertTextSet(dc, beequips[idx].Buffs), "#16a34a")
			currY += drawTextSet(dc, posX+padding, currY, convertTextSet(dc, beequips[idx].Debuffs), "#ef4444")
			currY += drawTextSet(dc, posX+padding, currY, convertBooleanSet(beequips[idx].Ability), "#ca8a04")
			currY += drawTextSet(dc, posX+padding, currY, convertTextSet(dc, beequips[idx].Bonuses), "#ca8a04")
			currY += 30
			centerX := posX + padding + 128
			for i, waxType := range beequips[idx].Waxes {
				waxPosX := 0
				add := 0
				if len(beequips[idx].Waxes)%2 == 0 {
					add = 25
				}
				if i >= len(beequips[idx].Waxes)/2+1 {
					waxPosX = centerX + 50*(i-len(beequips[idx].Waxes)/2) + add
				} else if i == len(beequips[idx].Waxes)/2 {
					waxPosX = centerX + add
				} else {
					waxPosX = centerX - 50*(i+1) + add
				}
				//fmt.Println(waxPosX)
				waxImg := loaders.GetWaxImage(waxType)
				dc.DrawImageAnchored(waxImg, waxPosX, currY-50, 0.5, 0)
			}
			idx++
		}
	}
}

func convertTextSet(dc *gg.Context, textSet map[string]int) []string {
	var finalText []string
	for k, v := range textSet {
		if v > 0 {
			textLen, _ := dc.MeasureString(strconv.Itoa(v) + k)
			if textLen > 128 {
				finalText = append(finalText, splitText(strconv.Itoa(v)+k)...)
			} else {
				finalText = append(finalText, strconv.Itoa(v)+k)
			}
		}
	}
	return finalText
}

func splitText(text string) []string {
	words := strings.Split(text, " ")
	var finalText []string
	currText := ""
	for _, word := range words {
		if len(currText)+len(word) > 20 {
			finalText = append(finalText, currText)
			currText = word
		} else {
			currText += " " + word
		}
	}
	finalText = append(finalText, currText)
	return finalText
}

func convertBooleanSet(booleanSet map[string]bool) []string {
	var finalText []string
	for k, v := range booleanSet {
		if v {
			finalText = append(finalText, k)
		}
	}
	return finalText
}

func drawTextSet(dc *gg.Context, posXStart int, posYStart int, text []string, color string) int {
	lastY := 0
	for _, v := range text {
		dc.SetHexColor(color)
		dc.DrawStringWrapped(v, float64(posXStart), float64(posYStart+lastY), 0, 0, 128, 1.5, gg.AlignLeft)
		// modify last y, accounting for the height of the text due to wrap
		lines := dc.WordWrap(v, 128)
		lastY += len(lines) * 30
	}
	return lastY
}
