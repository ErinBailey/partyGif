package main

import (
	"flag"
	"fmt"
	"image/color"
	"image/gif"
	"os"
	"strings"

	colorful "github.com/lucasb-eyer/go-colorful"
)

func readGif(fileIn, fileOut string) *gif.GIF {
	notGif := !strings.Contains(fileIn, ".gif")

	if notGif {
		fmt.Println("This is embarrassing...right now I can only convert gifs")
		return nil
	}

	originalImage, err := os.Open(fileIn)
	if err != nil {
		fmt.Println("Error reading image", err)
		return nil
	}
	defer originalImage.Close()

	partyImage, err := gif.DecodeAll(originalImage)

	if err != nil {
		fmt.Println("Error decoding image", err)
		return nil
	}
	frameCount := 360.0 / float64(len(partyImage.Image))
	for i := 0; i < len(partyImage.Image); i++ {
		for j := 0; j < len(partyImage.Image[i].Palette); j++ {
			_, _, _, a := partyImage.Image[i].Palette[j].RGBA()
			if a == 0 {
				continue
			}
			colorfulColor, _ := colorful.MakeColor(partyImage.Image[i].Palette[j])
			h, c, l := colorfulColor.Hcl()
			h += float64(i) * frameCount
			newHueColor := colorful.Hcl(h, c, l)
			partyImage.Image[i].Palette[j] = alphaOverride{color: newHueColor.Clamped(), alpha: a}
		}
	}

	newGrayPartyImage, err := os.Create(fileOut)
	if err != nil {
		fmt.Println("Error creating new image", err)
		return nil
	}
	defer newGrayPartyImage.Close()

	finalPartyImage := gif.EncodeAll(newGrayPartyImage, partyImage)

	if finalPartyImage != nil {
		fmt.Println("Error writing file", finalPartyImage)
		return nil
	}

	return partyImage
}

type alphaOverride struct {
	color color.Color
	alpha uint32
}

func (a alphaOverride) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b, _ := a.color.RGBA()
	return r, g, b, a.alpha
}

func main() {
	var fileIn string
	var fileOut string

	flag.StringVar(&fileIn, "in", "", "filepath of the gif you want to party")
	flag.StringVar(&fileOut, "out", "", "filepath of where you want to save your new party gif")

	flag.Parse()
	readGif(fileIn, fileOut)
}
