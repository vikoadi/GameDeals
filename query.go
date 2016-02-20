package main

import (
	"cheapshark"
	"launchpad.net/go-unityscopes/v2"
	"log"
	"math"
	"strconv"
)

type Query struct {
}

const queryLimit = 30

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
    "card-size": "medium"
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
    "card-size": "medium"
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
    "card-size": "medium"
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
    "card-size": "small"
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

	bestDealsReq := cheapshark.DealsRequest{SortBy: "Deal Rating", OnSale: true, Steamworks: steamworks, UpperPrice: max_price, PageSize: queryLimit}
	if deals, e := cs.Deals(&bestDealsReq); e != nil {
		log.Println(e)
		return e
	} else if err := registerCategory(reply, "bestDeals", "Best Deals", bestDealsCategoryTemplate, deals); err != nil {
		log.Println(e)
		return err
	}

	savingDealsReq := cheapshark.DealsRequest{SortBy: "Savings", OnSale: true, Steamworks: steamworks, UpperPrice: max_price, PageSize: queryLimit}
	if deals, e := cs.Deals(&savingDealsReq); e != nil {
		log.Println(e)
		return e
	} else if err := registerCategory(reply, "saving", "Most Saving", savingCategoryTemplate, deals); err != nil {
		log.Println(e)
		return err
	}

	cheapestDealsReq := cheapshark.DealsRequest{SortBy: "Price", OnSale: true, Steamworks: steamworks, UpperPrice: max_price, PageSize: queryLimit}
	if deals, e := cs.Deals(&cheapestDealsReq); e != nil {
		log.Println(e)
		return e
	} else if err := registerCategory(reply, "cheapest", "Cheapest", cheapestCategoryTemplate, deals); err != nil {
		log.Println(e)
		return err
	}

	bestGameDealsReq := cheapshark.DealsRequest{SortBy: "Metacritic", OnSale: true, Steamworks: steamworks, UpperPrice: max_price, PageSize: queryLimit}
	if deals, e := cs.Deals(&bestGameDealsReq); e != nil {
		log.Println(e)
		return e
	} else if err := registerCategory(reply, "best", "Popular Games", bestGameCategoryTemplate, deals); err != nil {
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
		} else if err := registerCategory(reply, store.StoreID, store.StoreName, searchCategoryTemplate, deals); err != nil {
			return err
		}
	}
	return nil
}

func registerCategory(reply *scopes.SearchReply, id string, title string, template string, deals cheapshark.Deal) error {
	category := reply.RegisterCategory(id, title, "", template)

	result := scopes.NewCategorisedResult(category)

	for _, d := range deals {
		savingsF, _ := d.Savings.Float64()

		if info, err := gb.GetInfo(d.Title); err != nil {
			// cant find data from GB database, use cheapshark one
			addCategorisedGameResult(result, "http://www.cheapshark.com/redirect?dealID="+d.DealID, d.Title, d.Title, d.NormalPrice.String(), d.SalePrice.String(), strconv.Itoa(int(math.Floor(savingsF))), d.MetacriticScore.String(), d.DealRating.String(), d.Thumb)
		} else {
			addCategorisedGameResult(result, "http://www.cheapshark.com/redirect?dealID="+d.DealID, d.Title, d.Title, d.NormalPrice.String(), d.SalePrice.String(), strconv.Itoa(int(math.Floor(savingsF))), d.MetacriticScore.String(), d.DealRating.String(), info.Image.ThumbURL)
		}
		if err := reply.Push(result); err != nil {
			return err
		}
	}

	return nil
}

func addCategorisedGameResult(result *scopes.CategorisedResult, uri string, dndUri string, title string, normalPrice string, salePrice string, savings string, metacriticScore string, dealRating string, art string) error {

	result.SetURI(uri)
	result.SetDndURI(dndUri)
	result.SetTitle(title)
	result.SetArt(art)
	result.Set("normalPrice", normalPrice)
	result.Set("salePrice", "$"+salePrice+" from $"+normalPrice+" ("+savings+"%)")
	result.Set("savings", savings)
	result.Set("metacriticScore", metacriticScore)
	result.Set("uri", uri)

	type Attr struct {
		Value string `json:"value"`
		Icon  string `json:"value"`
	}

	result.Set("attributes", []Attr{
		Attr{"one", ""},
		Attr{"two", ""},
	})

	return nil
}
