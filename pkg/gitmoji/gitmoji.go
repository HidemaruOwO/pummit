package gitmoji

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Gitmoji struct {
	Emoji       string `json:"emoji"`
	Entity      string `json:"entity"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type GitmojiResponse struct {
	Gitmojis []Gitmoji `json:"gitmojis"`
}

func GetAllGitmojis() ([]Gitmoji, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("https://raw.githubusercontent.com/carloscuesta/gitmoji/master/packages/gitmojis/src/gitmojis.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch gitmojis")
	}

	var response GitmojiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Gitmojis, nil
}

func FindByCode(code string, gitmojis []Gitmoji) (Gitmoji, bool) {
	for _, gitmoji := range gitmojis {
		if gitmoji.Code == code {
			return gitmoji, true
		}
	}
	return Gitmoji{}, false
}

func FindByName(name string, gitmojis []Gitmoji) (Gitmoji, bool) {
	for _, gitmoji := range gitmojis {
		if gitmoji.Name == name {
			return gitmoji, true
		}
	}
	return Gitmoji{}, false
}
