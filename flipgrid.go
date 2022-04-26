package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FlipgridCategory struct {
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

	StickerCount int               `json:"sticker_count"`
	Stickers     []FlipgridSticker `json:"stickers"`
}

type FlipgridSticker struct {
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

type FlipgridResponse[T interface{}] struct {
	Metadata FlipgridResponseMetadata `json:"metadata"`
	Data     []T                      `json:"data"`
}

type FlipgridResponseMetadata struct {
	Pagination FlipgridResponseMetadataPagination `json:"pagination"`
}

type FlipgridResponseMetadataPagination struct {
	Total       int  `json:"total"`
	TotalPages  int  `json:"total_pages"`
	FirstPage   bool `json:"first_page"`
	LastPage    bool `json:"last_page"`
	CurrentPage int  `json:"current_page"`
	Limit       int  `json:"limit"`
	Offset      int  `json:"offset"`
}

func getFlipgridStickers() {
	res, err := http.Get("https://api.flipgrid.com/api/sticker_categories")
	if err != nil {
		panic(err)
	}

	var body FlipgridResponse[FlipgridCategory]
	json.NewDecoder(res.Body).Decode(&body)

	wg := sync.WaitGroup{}
	for _, category := range body.Data {
		os.Mkdir(filepath.Join("flipgrid", category.Name), os.ModePerm)
		fmt.Println("Downloading category", category.Name)
		downloadStickers(wg, category, 0)
	}

	wg.Wait()
}

func downloadStickers(wg sync.WaitGroup, category FlipgridCategory, page int) {
	res, err := http.Get(fmt.Sprintf("https://api.flipgrid.com/api/sticker_categories/%d/stickers?page=%d", category.ID, page))
	if err != nil {
		panic(err)
	}

	var stickers FlipgridResponse[FlipgridSticker]
	json.NewDecoder(res.Body).Decode(&stickers)

	for _, sticker := range stickers.Data {
		wg.Add(1)
		go func(sticker FlipgridSticker) {
			png, err := http.Get(sticker.Assets.Png)
			if err != nil {
				panic(err)
			}

			file, err := os.Create(filepath.Join("flipgrid", category.Name, sticker.Name+".png"))
			if err != nil {
				panic(err)
			}

			if _, err := io.Copy(file, png.Body); err != nil {
				panic(err)
			}

			fmt.Println(" | Downloaded emoji", sticker.Name)
			wg.Done()
		}(sticker)
	}

	if !stickers.Metadata.Pagination.LastPage {
		downloadStickers(wg, category, stickers.Metadata.Pagination.CurrentPage+1)
	}
}
