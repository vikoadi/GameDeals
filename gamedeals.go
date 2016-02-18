package main

import (
	"cheapshark"
	"launchpad.net/go-unityscopes/v2"
	"log"
)

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
    "dealRating": "dealRating"
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
    "subtitle":"salePrice"
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
    "subtitle":"salePrice"
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

// SCOPE ***********************************************************************

var scope_interface scopes.Scope

type GameDealsScope struct {
	base *scopes.ScopeBase
}

func (s *GameDealsScope) Preview(result *scopes.Result, metadata *scopes.ActionMetadata, reply *scopes.PreviewReply, cancelled <-chan bool) error {
	layout1col := scopes.NewColumnLayout(1)
	layout2col := scopes.NewColumnLayout(2)
	layout3col := scopes.NewColumnLayout(3)

	// Single column layout
	layout1col.AddColumn("image", "header", "summary", "actions")

	// Two column layout
	layout2col.AddColumn("image")
	layout2col.AddColumn("header", "summary", "actions")

	// Three cokumn layout
	layout3col.AddColumn("image")
	layout3col.AddColumn("header", "summary", "actions")
	layout3col.AddColumn()

	// Register the layouts we just created
	reply.RegisterLayout(layout1col, layout2col, layout3col)

	header := scopes.NewPreviewWidget("header", "header")

	// It has title and a subtitle properties
	header.AddAttributeMapping("title", "title")
	header.AddAttributeMapping("subtitle", "subtitle")

	// Define the image section
	image := scopes.NewPreviewWidget("image", "image")
	// It has a single source property, mapped to the result's art property
	image.AddAttributeMapping("source", "art")

	// Define the summary section
	description := scopes.NewPreviewWidget("summary", "text")
	// It has a text property, mapped to the result's description property
	description.AddAttributeMapping("text", "description")

	// build variant map.
	tuple1 := make(map[string]interface{})
	tuple1["id"] = "open"
	tuple1["label"] = "Open"
	tuple1["uri"] = "application:///tmp/non-existent.desktop"

	tuple2 := make(map[string]interface{})
	tuple1["id"] = "download"
	tuple1["label"] = "Download"

	tuple3 := make(map[string]interface{})
	tuple1["id"] = "hide"
	tuple1["label"] = "Hide"

	actions := scopes.NewPreviewWidget("actions", "actions")
	actions.AddAttributeValue("actions", []interface{}{tuple1, tuple2, tuple3})

	var scope_data string
	metadata.ScopeData(scope_data)
	if len(scope_data) > 0 {
		extra := scopes.NewPreviewWidget("extra", "text")
		extra.AddAttributeValue("text", "test Text")
		reply.PushWidgets(header, image, description, actions, extra)
	} else {
		reply.PushWidgets(header, image, description, actions)
	}

	return nil
}

func (s *GameDealsScope) Search(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply, cancelled <-chan bool) error {
	//root_department := s.CreateDepartments(query, metadata, reply)
	//reply.RegisterDepartments(root_department)

	// test incompatible features in RTM version of libunity-scopes
	filter1 := scopes.NewOptionSelectorFilter("f1", "Options", false)
	var filterState scopes.FilterState
	// for RTM version of libunity-scopes we should see a log message
	reply.PushFilters([]scopes.Filter{filter1}, filterState)

	return s.AddQueryResults(reply, query.QueryString())
}

func (s *GameDealsScope) SetScopeBase(base *scopes.ScopeBase) {
	s.base = base
}

// RESULTS *********************************************************************

func (s *GameDealsScope) AddQueryResults(reply *scopes.SearchReply, query string) error {
	if query == "" {
		return s.AddEmptyQueryResults(reply, query)
	} else {
		return s.AddSearchResults(reply, query)
	}
}
func (s *GameDealsScope) AddEmptyQueryResults(reply *scopes.SearchReply, query string) error {
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

	//if err := registerCategory(reply, "bundles", "Game Bundles", bundleCategoryTemplate, "query"); err != nil {
	//return err
	//}

	return nil
}

func (s *GameDealsScope) AddSearchResults(reply *scopes.SearchReply, query string) error {
	return nil
}

func registerCategory(reply *scopes.SearchReply, id string, title string, template string, query cheapshark.Deal) error {
	category := reply.RegisterCategory(id, title, "", template)

	result := scopes.NewCategorisedResult(category)

	for _, d := range query {
		addCategorisedGameResult(result, d.Title, d.Title, d.Title, d.NormalPrice, d.SalePrice, d.Savings, d.MetacriticScore, d.DealRating, d.Thumb)
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
	result.Set("normalPrice", "<b>"+normalPrice+"</b>")
	result.Set("salePrice", "$"+salePrice+ " from $"+normalPrice+" ("+savings+"%")
	result.Set("savings", savings)
	result.Set("metacriticScore", metacriticScore)
	result.Set("dealRating", dealRating)

	return nil
}

// DEPARTMENTS *****************************************************************

func SearchDepartment(root *scopes.Department, id string) *scopes.Department {
	sub_depts := root.Subdepartments()
	for _, element := range sub_depts {
		if element.Id() == id {
			return element
		}
	}
	return nil
}

func (s *GameDealsScope) GetRockSubdepartments(query *scopes.CannedQuery,
	metadata *scopes.SearchMetadata,
	reply *scopes.SearchReply) *scopes.Department {
	active_dep, err := scopes.NewDepartment("Rock", query, "Rock Music")
	if err == nil {
		active_dep.SetAlternateLabel("Rock Music Alt")
		department, _ := scopes.NewDepartment("60s", query, "Rock from the 60s")
		active_dep.AddSubdepartment(department)

		department2, _ := scopes.NewDepartment("70s", query, "Rock from the 70s")
		active_dep.AddSubdepartment(department2)
	}

	return active_dep
}

func (s *GameDealsScope) GetSoulSubdepartments(query *scopes.CannedQuery,
	metadata *scopes.SearchMetadata,
	reply *scopes.SearchReply) *scopes.Department {
	active_dep, err := scopes.NewDepartment("Soul", query, "Soul Music")
	if err == nil {
		active_dep.SetAlternateLabel("Soul Music Alt")
		department, _ := scopes.NewDepartment("Motown", query, "Motown Soul")
		active_dep.AddSubdepartment(department)

		department2, _ := scopes.NewDepartment("New Soul", query, "New Soul")
		active_dep.AddSubdepartment(department2)
	}

	return active_dep
}

func (s *GameDealsScope) CreateDepartments(query *scopes.CannedQuery,
	metadata *scopes.SearchMetadata,
	reply *scopes.SearchReply) *scopes.Department {
	department, _ := scopes.NewDepartment("", query, "Browse Music")
	department.SetAlternateLabel("Browse Music Alt")

	rock_dept := s.GetRockSubdepartments(query, metadata, reply)
	if rock_dept != nil {
		department.AddSubdepartment(rock_dept)
	}

	soul_dept := s.GetSoulSubdepartments(query, metadata, reply)
	if soul_dept != nil {
		department.AddSubdepartment(soul_dept)
	}

	return department
}

// MAIN ************************************************************************

func main() {
	if err := scopes.Run(&GameDealsScope{}); err != nil {
		log.Fatalln(err)
	}
}
