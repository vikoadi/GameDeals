package main

import (
	"cheapshark"
	"launchpad.net/go-unityscopes/v2"
	"math"
	"strconv"
)

type Query struct {
}

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
					"aspect-ratio":2.5,
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

func (s *Query) AddQueryResults(reply *scopes.SearchReply, query string) error {
	if query == "" {
		return s.AddEmptyQueryResults(reply, query)
	} else {
		return s.AddSearchResults(reply, query)
	}
}
func (s *Query) AddEmptyQueryResults(reply *scopes.SearchReply, query string) error {
	var cs cheapshark.CheapShark

	bestDealsReq := cheapshark.DealsRequest{SortBy: "Deal Rating", OnSale: true}
	if err := registerCategory(reply, "bestDeals", "Best Deals", bestDealsCategoryTemplate, cs.Deals(&bestDealsReq)); err != nil {
		return err
	}

	savingDealsReq := cheapshark.DealsRequest{SortBy: "Savings", OnSale: true}
	if err := registerCategory(reply, "saving", "Most Saving", savingCategoryTemplate, cs.Deals(&savingDealsReq)); err != nil {
		return err
	}

	cheapestDealsReq := cheapshark.DealsRequest{SortBy: "Price", OnSale: true}
	if err := registerCategory(reply, "cheapest", "Cheapest", cheapestCategoryTemplate, cs.Deals(&cheapestDealsReq)); err != nil {
		return err
	}

	bestGameDealsReq := cheapshark.DealsRequest{SortBy: "Metacritic", OnSale: true}
	if err := registerCategory(reply, "best", "Best Games", bestGameCategoryTemplate, cs.Deals(&bestGameDealsReq)); err != nil {
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
		if err := registerCategory(reply, store.StoreID, store.StoreName, searchCategoryTemplate, cs.Deals(&searchReq)); err != nil {
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

		addCategorisedGameResult(result, "http://www.cheapshark.com/redirect?dealID="+d.DealID, d.Title, d.Title, d.NormalPrice.String(), d.SalePrice.String(), strconv.Itoa(int(math.Floor(savingsF))), d.MetacriticScore.String(), d.DealRating.String(), d.Thumb)
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
