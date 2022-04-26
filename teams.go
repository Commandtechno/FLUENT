package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

const TEAMS_EMOTICONS_HASH = "8479a21d47934aed85d6d6f236847484"

type TeamsResposne struct {
	Categories []TeamsCategory `json:"categories"`
}

type TeamsCategory struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Emoticons   []TeamsEmoticon `json:"emoticons"`
}

type TeamsEmoticon struct {
	ID          string         `json:"id"`
	Description string         `json:"description"`
	Shortcuts   []string       `json:"shortcuts"`
	Unicode     string         `json:"unicode"`
	Etag        string         `json:"etag"`
	Diverse     bool           `json:"diverse"`
	Animation   TeamsAnimation `json:"animation"`
	Keywords    []string       `json:"keywords"`
}

type TeamsAnimation struct {
	Fps         int `json:"fps"`
	FramesCount int `json:"framesCount"`
	FirstFrame  int `json:"firstFrame"`
}

func getTeamsEmojis() {
	res, err := http.Get(fmt.Sprintf("https://statics.teams.cdn.office.net/evergreen-assets/personal-expressions/v1/metadata/%s/en-us.json", TEAMS_EMOTICONS_HASH))
	if err != nil {
		panic(err)
	}

	var body TeamsResposne
	json.NewDecoder(res.Body).Decode(&body)

	for _, category := range body.Categories {
		os.Mkdir(filepath.Join("teams", category.Title), os.ModePerm)
		fmt.Println("Downloading category", category.Title)

		wg := &sync.WaitGroup{}
		category.downloadEmojis(wg)
		wg.Wait()
	}

	for _, category := range body.Categories {
		fmt.Println("Rendering category", category.Title)
		category.renderEmojis()
	}
}

func (category *TeamsCategory) downloadEmojis(wg *sync.WaitGroup) {
	for _, emoticon := range category.Emoticons {
		wg.Add(1)
		go func(emoticon TeamsEmoticon) {
			img, err := http.Get(fmt.Sprintf("https://statics.teams.cdn.office.net/evergreen-assets/personal-expressions/v2/assets/emoticons/%s/default/100_anim_f.png", emoticon.ID))
			if err != nil {
				panic(err)
			}

			file, err := os.Create(filepath.Join("teams", category.Title, emoticon.Description+".png"))
			if err != nil {
				panic(err)
			}

			if _, err := io.Copy(file, img.Body); err != nil {
				panic(err)
			}

			fmt.Println(" | Downloaded emoji", emoticon.Description)
			wg.Done()
		}(emoticon)
	}

}

func (category *TeamsCategory) renderEmojis() {
	magick, err := exec.LookPath("magick")
	if err != nil {
		magick, err = exec.LookPath("convert")
	}

	if err != nil {
		panic(err)
	}

	for _, emoticon := range category.Emoticons {
		delay := 100 / emoticon.Animation.Fps
		path := filepath.Join("teams", category.Title, emoticon.Description+".png")
		cmd := exec.Command(magick, "-delay", fmt.Sprintf("%d", delay), "-dispose", "previous", path, "-crop", "100x100", "+repage", "APNG:"+path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}

		fmt.Println(" | Rendered", emoticon.Description)
	}
}
