package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var api_key = [insert-your-api-key]

type GiantBomb struct {
}

type Results struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Platforms    Platforms `json:"platforms"`
	Image        Image     `json:"image"`
	ResourceType string    `json:"resource_type"`
}
type Image struct {
	IconURL   string `json:"icon_url"`
	MediumURL string `json:"medium_url"`
	ScreenURL string `json:"screen_url"`
	SmallURL  string `json:"small_url"`
	SuperURL  string `json:"super_url"`
	ThumbURL  string `json:"thumb_url"`
	TinyURL   string `json:"tiny_url"`
}
type Platforms []struct {
	APIDetailURL  string `json:"api_detail_url"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	SiteDetailURL string `json:"site_detail_url"`
	Abbreviation  string `json:"abbreviation"`
}
type Result struct {
	Error                string    `json:"error"`
	Limit                int       `json:"limit"`
	Offset               int       `json:"offset"`
	NumberOfPageResults  int       `json:"number_of_page_results"`
	NumberOfTotalResults int       `json:"number_of_total_results"`
	StatusCode           int       `json:"status_code"`
	Res                  []Results `json:"results"`
	Version              string    `json:"version"`
}

func (g *GiantBomb) GetInfo(gameName string) (r Results, err error) {
	var res Result
	if e := getJson("http://www.giantbomb.com/api/search/?api_key="+api_key+"&field_list=image,name,genres,platforms,description&limit=1&format=json&resources=game&query="+gameName, &res); e != nil {
		log.Fatal(e)
		err = e
		return
	}
	r = res.Res[0]
	return
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
