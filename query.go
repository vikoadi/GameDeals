package main

import (
	"cheapshark"
	"launchpad.net/go-unityscopes/v2"
	"log"
	"math"
	"strconv"
	"time"
)

type Query struct {
	scopeDirectory string
}

const windowsQueryLimit = 20

const bestDealsCategoryTemplate = `{
  "schema-version": 1,
  "template": {
    "category-layout": "carousel",
                "card-layout": "vertical",
                "card-size": "large"
  },
  "components": {
    "title": "title",
                "art" : {
					"aspect-ratio":0.725,
                    "field": "art"
                },
    "subtitle":"salePrice",
    "attributes":"attributes"
  }
}`

const savingCategoryTemplate = `{
  "schema-version": 1,
  "template": {
    "category-layout": "horizontal-list",
    "card-size": "medium",
    "overlay":true
  },
  "components": {
    "title": "title",
    "art":  "art",
    "subtitle":"salePrice",
    "attributes":"attributes"
  }
}`

const cheapestCategoryTemplate = `{
  "schema-version": 1,
  "template": {
    "category-layout": "horizontal-list",
    "card-size": "medium",
    "overlay":true
  },
  "components": {
    "title": "title",
    "art":  "art",
    "subtitle":"salePrice",
    "attributes":"attributes"
  }
}`

const bestGameCategoryTemplate = `{
  "schema-version": 1,
  "template": {
    "category-layout": "horizontal-list",
    "card-size": "medium",
    "overlay":true
  },
  "components": {
    "title": "title",
    "art":  "art",
    "subtitle":"salePrice",
    "attributes":"attributes"
  }
}`

const bundleCategoryTemplate = `{
  "schema-version": 1,
  "template": {
    "category-layout": "vertical-journal",
    "card-layout" : "horizontal",
    "card-size": "small"
  },
  "components": {
    "title": "title",
    "art":  "art",
    "subtitle": "username"
  }
}`
const searchCategoryTemplate = `{
  "schema-version": 1,
  "template": {
    "category-layout": "vertical-journal",
    "card-layout" : "horizontal",
    "card-size": "20"
  },
  "components": {
    "title": "title",
    "art":  "art",
    "subtitle": "salePrice"
  }
}`

func (s *Query) AddQueryResults(reply *scopes.SearchReply, query string, settings Settings) error {
	if query == "" {
		return s.AddEmptyQueryResults(reply, query, settings)
	} else {
		return s.AddSearchResults(reply, query)
	}
}

type Category struct {
	DealsReq       cheapshark.Deal
	Id             string
	Title          string
	Template       string
	CompleteDetail bool
}

var cs cheapshark.CheapShark

func (s *Query) AddEmptyQueryResults(reply *scopes.SearchReply, query string, settings Settings) error {

	steamworks := settings.Steamworks
	max_price := settings.MaxPrice

	// cheapshark treat 50 as no limit, but we want different behaviour
	if max_price == 50 {
		max_price = 51
	} else if max_price == 0 {
		max_price = 50
	}

	// limit cheapshark request if Windows is enabled, because Windows
	// compatible games are too much
	queryLimit := 0
	if settings.Windows {
		queryLimit = windowsQueryLimit
	}

	bestDealsReq := cheapshark.DealsRequest{SortBy: "Deal Rating", OnSale: true, Steamworks: steamworks, UpperPrice: max_price, PageSize: queryLimit}
	savingDealsReq := cheapshark.DealsRequest{SortBy: "Savings", OnSale: true, Steamworks: steamworks, UpperPrice: max_price, PageSize: queryLimit}
	cheapestDealsReq := cheapshark.DealsRequest{SortBy: "Price", OnSale: true, Steamworks: steamworks, UpperPrice: max_price, PageSize: queryLimit}
	bestGameDealsReq := cheapshark.DealsRequest{SortBy: "Metacritic", OnSale: true, Steamworks: steamworks, UpperPrice: max_price, PageSize: queryLimit}

	deals := make(chan Category)
	go func() {
		deals <- createDeals(bestDealsReq, "deals", "Best Deals", bestDealsCategoryTemplate, true)
		deals <- createDeals(savingDealsReq, "saving", "Most Saving", savingCategoryTemplate, true)
		deals <- createDeals(cheapestDealsReq, "cheapest", "Cheapest", cheapestCategoryTemplate, true)
		deals <- createDeals(bestGameDealsReq, "best", "Popular Games", bestGameCategoryTemplate, true)
		close(deals)
	}()

	if err := s.registerCategory(reply, deals); err != nil {
		log.Println(err)
		return err
	}

	//if err := registerCategory(reply, "bundles", "Game Bundles", bundleCategoryTemplate, "query"); err != nil {
	//return err
	//}

	return nil
}

func createDeals(dealsReq cheapshark.DealsRequest, id string, title string, template string, complete bool) Category {
	cat := Category{
		Id:             id,
		Title:          title,
		Template:       template,
		CompleteDetail: complete,
	}
	log.Println("createDeals ", title)

	if deals, e := cs.Deals(&dealsReq); e != nil {
		log.Println(e)
		return cat
	} else {
		cat.DealsReq = deals
		return cat
	}
}

var stores cheapshark.Store

func (s *Query) AddSearchResults(reply *scopes.SearchReply, query string) error {
	var cs cheapshark.CheapShark

	if stores == nil {
		stores = cs.Stores()
	}

	for _, store := range stores {
		searchReq := cheapshark.DealsRequest{Title: query, StoreID: store.StoreID}
		if _, err := cs.Deals(&searchReq); err != nil {
			log.Println(err)
			return err
		}
		//else if err := s.registerCategory(reply, store.StoreID, store.StoreName, searchCategoryTemplate, deals, false); err != nil {
		//return err
		//}
	}
	return nil
}

func getDateString(unixDate int64) string {
	return time.Unix(unixDate, 0).Format("2 January 2006")
}

func (s *Query) registerCategory(reply *scopes.SearchReply, cats <-chan Category) error {
	for cat := range cats {
		log.Println("registerCategory ", cat.Title)
		category := reply.RegisterCategory(cat.Id, cat.Title, "", cat.Template)
		res := make(chan *scopes.CategorisedResult)
		go func() {
			result := scopes.NewCategorisedResult(category)

			for _, d := range cat.DealsReq {
				savingsF, _ := d.Savings.Float64()
				releaseDate, _ := d.ReleaseDate.Int64()
				releaseDateStr := ""
				if r := getDateString(releaseDate); releaseDate != 0 {
					releaseDateStr = r
				}

				storeIcon := s.getStoreIcon(d.StoreID)
				if cat.CompleteDetail {
					if info, err := gb.GetInfo(d.Title); err == nil {
						if Filter(info.Platforms, platformFilter) {
							addCategorisedGameResult(result, "http://www.cheapshark.com/redirect?dealID="+d.DealID, d.Title, d.Title, d.NormalPrice.String(), d.SalePrice.String(), strconv.Itoa(int(math.Floor(savingsF))), d.MetacriticScore.String(), d.DealRating.String(), info.Image.ThumbURL, info.Image.SmallURL, storeIcon, info.Description, releaseDateStr, s.getStoreIcon("platform"+strconv.Itoa(GetPlatformsFilter(info.Platforms))))
							res <- result
							//reply.Push(result)
						}
						continue
					}
				}
				// cant find data from GB database, use cheapshark one
				if platformFilter&8 > 0 { // only add if unknown platform is enabled
					addCategorisedGameResult(result, "http://www.cheapshark.com/redirect?dealID="+d.DealID, d.Title, d.Title, d.NormalPrice.String(), d.SalePrice.String(), strconv.Itoa(int(math.Floor(savingsF))), d.MetacriticScore.String(), d.DealRating.String(), d.Thumb, d.Thumb, storeIcon, "", releaseDateStr, "")
					//reply.Push(result)
					res <- result
				}
			}
			close(res)
		}()
		for r := range res {
			reply.Push(r)
		}
	}

	return nil
}

func addCategorisedGameResult(result *scopes.CategorisedResult, uri string, dndUri string, title string, normalPrice string, salePrice string, savings string, metacriticScore string, dealRating string, art string, bigArt string, storeIcon string, description string, releaseDate string, platformsIcon string) error {

	result.SetURI(uri)
	result.SetDndURI(dndUri)
	result.SetTitle(title)
	result.SetArt(art)
	result.Set("bigArt", bigArt)
	result.Set("normalPrice", normalPrice)
	if salePrice == "0" {
		result.Set("salePrice", "<b>FREE</b> from $"+normalPrice)
	} else {
		result.Set("salePrice", "<b>$"+salePrice+"</b> from $"+normalPrice)
	}
	result.Set("uri", uri)
	if description != "" {
		result.Set("description", description)
	} else {
		result.Set("description", "<h1>No Description</h1>")
	}

	type Attr struct {
		Value string `json:"value"`
		Icon  string `json:"icon"`
	}

	attr := []Attr{}
	if savings != "0" {
		attr = append(attr, Attr{Value: savings + "%", Icon: storeIcon})
	}
	if metacriticScore != "0" {
		attr = append(attr, Attr{Value: metacriticScore, Icon: "image://theme/starred"})
	}

	result.Set("attributes", attr)

	completeAttr := attr

	if releaseDate != "" {
		completeAttr = append(completeAttr, Attr{Value: "released at " + releaseDate})
	}
	if platformsIcon != "" {
		completeAttr = append(completeAttr, Attr{Icon: platformsIcon})
	}

	result.Set("completeAttributes", completeAttr)

	return nil
}

var platformFilter = 0

func (s *Query) SetPlatformFilter(linux bool, osx bool, windows bool, unknown bool) {
	platformFilter = 0
	if linux {
		platformFilter += 1
	}
	if osx {
		platformFilter += 2
	}
	if windows {
		platformFilter += 4
	}
	if unknown {
		platformFilter += 8
	}
}

func (s *Query) SetScopeDirectory(dir string) {
	s.scopeDirectory = dir
}
func (s *Query) getStoreIcon(storeID string) string {
	return s.scopeDirectory + "/" + storeID + ".png"
}
