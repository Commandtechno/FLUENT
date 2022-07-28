package main

import (
	"fmt"
	"os"
)

func main() {
	/*
		Flipgrid
	*/
	os.Mkdir("flipgrid", os.ModePerm)
	fmt.Println("Downloading Flipgrid stickers")
	getFlipgridStickers()

	/*
		Microsoft Teams
	*/
	os.Mkdir("teams", os.ModePerm)
	fmt.Println("Downloading Microsoft Teams animated emojis")
	getTeamsEmojis()

	/*
		Emojipedia
	*/
	os.Mkdir("emojipedia", os.ModePerm)
	fmt.Println("Downloading Emojipedia animated emojis")
	getEmojipediaEmojis()
}