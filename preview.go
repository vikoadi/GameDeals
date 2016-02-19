package main

import (
	"launchpad.net/go-unityscopes/v2"
)

type Preview struct {
}

func (p *Preview) AddPreviewResult(result *scopes.Result, metadata *scopes.ActionMetadata, reply *scopes.PreviewReply) error {
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

	var gb GiantBomb
	info := gb.GetInfo(result.Title())
	// Define the summary section
	description := scopes.NewPreviewWidget("summary", "text")
	// It has a text property, mapped to the result's description property
	//description.AddAttributeMapping("text", "description")
	description.AddAttributeValue("text", info.Description)

	// build variant map.
	tuple1 := make(map[string]interface{})
	tuple1["id"] = "open"
	tuple1["label"] = "Go to Store"
	tuple1["uri"] = result.URI()

	//tuple2 := make(map[string]interface{})
	//tuple1["id"] = "download"
	//tuple1["label"] = "Download"

	//tuple3 := make(map[string]interface{})
	//tuple1["id"] = "hide"
	//tuple1["label"] = "Hide"

	actions := scopes.NewPreviewWidget("actions", "actions")
	actions.AddAttributeValue("actions", []interface{}{tuple1})

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
