package main

import (
	"fmt"
	"os"
)

func main() {
	/*
		Flipgrid
	*/
	os.RemoveAll("flipgrid")
	os.Mkdir("flipgrid", os.ModePerm)

	fmt.Println("Downloading Flipgrid stickers")
	getFlipgridStickers()

	/*
		Microsoft Teams
	*/
	os.RemoveAll("teams")
	os.Mkdir("teams", os.ModePerm)

	fmt.Println("Downloading Microsoft Teams animated emojis")
	getTeamsEmojis()
}
