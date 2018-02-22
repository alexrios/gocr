package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"

	"gocv.io/x/gocv"
)

const templateChar = "ล ู ก ค ิ ด ม า เ ล ้ ว อ ฟ ห ่ ไ โ บ ซ บ ใ จ"
const templateDir = "templates/"
const outputDir = "outputs/"

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Invalid argument: gocr [mode, filename]")
	} else {
		if os.Args[1] == "--gentemp" {
			// Generate template
			GenTemplate(templateChar, templateDir)
		} else {
			// OCR

			// Read image
			fmt.Printf("Opening %v...\n", os.Args[1])

			img := gocv.IMRead(os.Args[1], gocv.IMReadGrayScale)

			// Apply auto threshold
			fmt.Println("Applying auto threshold...")

			newImg := AutoThreshold(img)

			gocv.IMWrite(outputDir+"02_auto_threshold.jpg", newImg)

			// Read template
			fmt.Println("Loading templates...")

			templates := ReadTemplate(templateChar, templateDir)

			// Row segmentation
			fmt.Println("Rows segmenting...")

			imgArr := GetImgArray(newImg)
			start, end := SplitLine(imgArr)

			DrawRowSegment(newImg, start, end)

			gocv.IMWrite(outputDir+"03_row_segment.jpg", newImg)

			// Character segmentation
			fmt.Println("Characters segmenting and template mathching...")
			fmt.Println(">>")

			for i := range start {
				row := CropImgArr(imgArr, image.Rectangle{image.Point{0, start[i]}, image.Point{len(imgArr[0]), end[i]}})
				rectTable := GetSegmentChar(row)

				rowImg := GetImgMat(row)

				for _, rect := range rectTable {
					gocv.Rectangle(rowImg, rect, color.RGBA{255, 0, 0, 0}, 1)
				}

				gocv.IMWrite(outputDir+"04_character_segment_"+strconv.Itoa(i)+".jpg", rowImg)

				for b := range rectTable {
					cropImg := CropImgArr(row, rectTable[b])

					fmt.Printf("%v", MatchTemplate(cropImg, templates[GetRatioBin(len(cropImg), len(cropImg[b]))])[0].char)
				}

				println()
			}

			fmt.Println("<<")

			fmt.Println("DONE!")

		}
	}

}
