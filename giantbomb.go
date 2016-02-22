package main

import (
	//"bytes"
	"encoding/json"
	"errors"
	"github.com/GitbookIO/diskache"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var api_key = [insert-your-api-key]
var GiantBombCache *diskache.Diskache

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

func (g *GiantBomb) SetCacheDirectory(directory string) {
	// cache directory
	opts := &diskache.Opts{
		Directory: directory + "/giantbomb",
	}
	GiantBombCache, _ = diskache.New(opts)
}

func (g *GiantBomb) GetInfo(gameName string) (r Results, err error) {

	var res Result
	log.Println("getinfo for ", gameName)
	if content, inCache := GiantBombCache.Get(gameName); inCache {
		if er := parseJson(string(content), &res); er != nil {
			log.Println(er)
			err = er
			return
		}
		if result, e := getExactName(res, gameName); e != nil {
			err = e
			return
		} else {
			r = result
			return
		}
	} else {
		url := "http://www.giantbomb.com/api/search/?api_key=" + api_key + "&field_list=image,name,genres,platforms,description&limit=1&format=json&resources=game&query=" + gameName
		resp, e1 := http.Get(url)
		if e1 != nil {
			err = e1
			defer resp.Body.Close()
			return
		}
		defer resp.Body.Close()

		if data, er2 := ioutil.ReadAll(resp.Body); er2 != nil {
			log.Println("cant readall", er2)
			return
		} else {
			if er := parseJson(string(data), &res); er != nil {
				log.Println(er)
				err = er
				return
			}

			GiantBombCache.Set(gameName, []byte(data))

			if result, e := getExactName(res, gameName); e != nil {
				err = e
				return
			} else {
				r = result
				return
			}
		}
	}

	return
}

func getExactName(result Result, gameName string) (res Results, err error) {
	if (len(result.Res)) <= 0 {
		err = errors.New("not enough data")
		return
	}
	for _, g := range result.Res {
		if g.Name == gameName {
			res = g
			return
		}
	}
	err = errors.New("Game Name not found in list")
	return
}

func parseJson(jsonStream string, target interface{}) error {
	return json.NewDecoder(strings.NewReader(jsonStream)).Decode(target)
}

func GetPlatformsFilter(platformsInfo Platforms) int {
	platform := 0
	for _, platformInfo := range platformsInfo {
		if platformInfo.Abbreviation == "LIN" {
			platform += 1
		}
		if platformInfo.Abbreviation == "MAC" {
			platform += 2
		}
		if platformInfo.Abbreviation == "PC" {
			platform += 4
		}
	}

	// found nothing, set it to unknown
	if platform == 0 {
		platform = 8
	}
	return platform
}

func Filter(platformsInfo Platforms, filter int) bool {
	return (filter & GetPlatformsFilter(platformsInfo)) > 0
}
