package main

import (
	"launchpad.net/go-unityscopes/v2"
	"log"
	//"fmt"
)

// SCOPE ***********************************************************************

var scope_interface scopes.Scope

type GameDealsScope struct {
	base *scopes.ScopeBase
}

var preview Preview

func (s *GameDealsScope) Preview(result *scopes.Result, metadata *scopes.ActionMetadata, reply *scopes.PreviewReply, cancelled <-chan bool) error {

	return preview.AddPreviewResult(result, metadata, reply)
}

var qu Query

func (s *GameDealsScope) Search(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply, cancelled <-chan bool) error {
	//root_department := s.CreateDepartments(query, metadata, reply)
	//reply.RegisterDepartments(root_department)

	// test incompatible features in RTM version of libunity-scopes
	filter1 := scopes.NewOptionSelectorFilter("f1", "Options", false)
	var filterState scopes.FilterState
	// for RTM version of libunity-scopes we should see a log message
	reply.PushFilters([]scopes.Filter{filter1}, filterState)

	var settings Settings
	s.base.Settings(&settings)

	return qu.AddQueryResults(reply, query.QueryString(), settings)
}

type Settings struct {
	Steamworks bool `json:"steamworks"`
	Localized  bool `json:"localized"`
	MaxPrice   int  `json:"max_price"`
}

func (s *GameDealsScope) SetScopeBase(base *scopes.ScopeBase) {
	s.base = base

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
