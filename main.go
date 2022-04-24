package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Position  int       `json:"position"`
	PerRow    int       `json:"per_row"`
	UpdatedAt time.Time `json:"updated_at"`
	Icons     struct {
		Svg string `json:"svg"`
		Pdf string `json:"pdf"`
		Png string `json:"png"`
	} `json:"icons"`

	StickerCount int       `json:"sticker_count"`
	Stickers     []Sticker `json:"stickers"`
}

type Sticker struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
	Assets   struct {
		Svg string `json:"svg"`
		Pdf string `json:"pdf"`
		Png string `json:"png"`
	} `json:"assets"`
	Name      string    `json:"name"`
	Position  int       `json:"position"`
	UpdatedAt time.Time `json:"updated_at"`
	Width     string    `json:"width"`
	Height    string    `json:"height"`
}

type Response[T interface{}] struct {
	Metadata ResponseMetadata `json:"metadata"`
	Data     []T              `json:"data"`
}

type ResponseMetadata struct {
	Pagination ResponseMetadataPagination `json:"pagination"`
}

type ResponseMetadataPagination struct {
	Total       int  `json:"total"`
	TotalPages  int  `json:"total_pages"`
	FirstPage   bool `json:"first_page"`
	LastPage    bool `json:"last_page"`
	CurrentPage int  `json:"current_page"`
	Limit       int  `json:"limit"`
	Offset      int  `json:"offset"`
}

func getCategories() []Category {
	res, err := http.Get("https://api.flipgrid.com/api/sticker_categories")
	if err != nil {
		panic(err)
	}

	var stickerCategories Response[Category]
	json.NewDecoder(res.Body).Decode(&stickerCategories)

	return stickerCategories.Data
}

func getStickers(category Category, page int) {
	res, err := http.Get(fmt.Sprintf("https://api.flipgrid.com/api/sticker_categories/%d/stickers?page=%d", category.ID, page))
	if err != nil {
		panic(err)
	}

	var stickers Response[Sticker]
	json.NewDecoder(res.Body).Decode(&stickers)

	for _, sticker := range stickers.Data {
		fmt.Println(" | Downloading sticker", sticker.Name)
		file, err := os.Create(filepath.Join("stickers", category.Name, sticker.Name+".png"))
		if err != nil {
			panic(err)
		}

		png, err := http.Get(sticker.Assets.Png)
		if err != nil {
			panic(err)
		}

		io.Copy(file, png.Body)
	}

	if stickers.Metadata.Pagination.LastPage {
		fmt.Println()
	} else {
		getStickers(category, stickers.Metadata.Pagination.CurrentPage+1)
	}
}

func main() {
	os.RemoveAll("stickers")
	os.Mkdir("stickers", os.ModePerm)

	categories := getCategories()
	for _, category := range categories {
		fmt.Println("Downloading category", category.Name)
		os.Mkdir(filepath.Join("stickers", category.Name), os.ModePerm)
		getStickers(category, 0)
	}
}