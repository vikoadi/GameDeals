package main

import (
	"cheapshark"
	"launchpad.net/go-unityscopes/v2"
	"log"
	"math"
	"strconv"
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
func (s *Query) AddEmptyQueryResults(reply *scopes.SearchReply, query string, settings Settings) error {
	var cs cheapshark.CheapShark

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
	if deals, e := cs.Deals(&bestDealsReq); e != nil {
		log.Println(e)
		return e
	} else if err := s.registerCategory(reply, "bestDeals", "Best Deals", bestDealsCategoryTemplate, deals, true); err != nil {
		log.Println(e)
		return err
	}

	savingDealsReq := cheapshark.DealsRequest{SortBy: "Savings", OnSale: true, Steamworks: steamworks, UpperPrice: max_price, PageSize: queryLimit}
	if deals, e := cs.Deals(&savingDealsReq); e != nil {
		log.Println(e)
		return e
	} else if err := s.registerCategory(reply, "saving", "Most Saving", savingCategoryTemplate, deals, true); err != nil {
		log.Println(e)
		return err
	}

	cheapestDealsReq := cheapshark.DealsRequest{SortBy: "Price", OnSale: true, Steamworks: steamworks, UpperPrice: max_price, PageSize: queryLimit}
	if deals, e := cs.Deals(&cheapestDealsReq); e != nil {
		log.Println(e)
		return e
	} else if err := s.registerCategory(reply, "cheapest", "Cheapest", cheapestCategoryTemplate, deals, true); err != nil {
		log.Println(e)
		return err
	}

	bestGameDealsReq := cheapshark.DealsRequest{SortBy: "Metacritic", OnSale: true, Steamworks: steamworks, UpperPrice: max_price, PageSize: queryLimit}
	if deals, e := cs.Deals(&bestGameDealsReq); e != nil {
		log.Println(e)
		return e
	} else if err := s.registerCategory(reply, "best", "Popular Games", bestGameCategoryTemplate, deals, true); err != nil {
		log.Println(e)
		return err
	}

	//if err := registerCategory(reply, "bundles", "Game Bundles", bundleCategoryTemplate, "query"); err != nil {
	//return err
	//}

	return nil
}

var stores cheapshark.Store

func (s *Query) AddSearchResults(reply *scopes.SearchReply, query string) error {
	var cs cheapshark.CheapShark

	if stores == nil {
		stores = cs.Stores()
	}

	for _, store := range stores {
		searchReq := cheapshark.DealsRequest{Title: query, StoreID: store.StoreID}
		if deals, err := cs.Deals(&searchReq); err != nil {
			log.Println(err)
			return err
		} else if err := s.registerCategory(reply, store.StoreID, store.StoreName, searchCategoryTemplate, deals, false); err != nil {
			return err
		}
	}
	return nil
}

func (s *Query) registerCategory(reply *scopes.SearchReply, id string, title string, template string, deals cheapshark.Deal, completeDetail bool) error {
	category := reply.RegisterCategory(id, title, "", template)

	result := scopes.NewCategorisedResult(category)

	for _, d := range deals {
		savingsF, _ := d.Savings.Float64()

		storeIcon := s.getStoreIcon(d.StoreID)
		log.Println(storeIcon)
		if completeDetail {
			if info, err := gb.GetInfo(d.Title); err == nil {
				if Filter(info.Platforms, platformFilter) {
					addCategorisedGameResult(result, "http://www.cheapshark.com/redirect?dealID="+d.DealID, d.Title, d.Title, d.NormalPrice.String(), d.SalePrice.String(), strconv.Itoa(int(math.Floor(savingsF))), d.MetacriticScore.String(), d.DealRating.String(), info.Image.ThumbURL, storeIcon, info.Description, "release", "icon.png")
					if err := reply.Push(result); err != nil {
						return err
					}
				}
				continue
			}
		}
		// cant find data from GB database, use cheapshark one
		if platformFilter&8 > 0 { // only add if unknown platform is enabled
			addCategorisedGameResult(result, "http://www.cheapshark.com/redirect?dealID="+d.DealID, d.Title, d.Title, d.NormalPrice.String(), d.SalePrice.String(), strconv.Itoa(int(math.Floor(savingsF))), d.MetacriticScore.String(), d.DealRating.String(), d.Thumb, storeIcon, "", "released on", "")
			if err := reply.Push(result); err != nil {
				return err
			}
		}
	}

	return nil
}

func addCategorisedGameResult(result *scopes.CategorisedResult, uri string, dndUri string, title string, normalPrice string, salePrice string, savings string, metacriticScore string, dealRating string, art string, storeIcon string, description string, releaseDate string, platformsIcon string) error {

	result.SetURI(uri)
	result.SetDndURI(dndUri)
	result.SetTitle(title)
	result.SetArt(art)
	result.Set("normalPrice", normalPrice)
	result.Set("salePrice", "<b>$"+salePrice+"</b> from $"+normalPrice)
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

	attr := []Attr{
		{Value: savings + "%", Icon: storeIcon},
		{Value: metacriticScore, Icon: "image://theme/starred"},
	}

	result.Set("attributes", attr)

	completeAttr := attr

	if releaseDate != "" {
		completeAttr = append(completeAttr, Attr{Value: releaseDate})
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
