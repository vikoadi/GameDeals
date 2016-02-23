package main

import (
	"launchpad.net/go-unityscopes/v2"
	"log"
	//"fmt"
)

// SCOPE ***********************************************************************

var scope_interface scopes.Scope
var gb GiantBomb
var departments = [][2]string{
	{"action", "Action"},
	{"actionadventure", "Action-adventure"},
	{"adventure", "Adventure"},
	{"fighting", "Fighting"},
	{"firstpersonshooter", "First-Person Shooter"},
	{"platformer", "Platformer"},
	{"puzzle", "Puzzle"},
	{"realtimestrategy", "Real-Time Strategy"},
	{"roleplaying", "Role Playing"},
	{"strategy", "Strategy"},
}

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

	gb.SetCacheDirectory(s.base.CacheDirectory())
	qu.SetScopeDirectory(s.base.ScopeDirectory())

	var settings Settings
	s.base.Settings(&settings)
	qu.SetPlatformFilter(settings.Linux, settings.Osx, settings.Windows, settings.Unknown)

	// test incompatible features in RTM version of libunity-scopes
	filter1 := scopes.NewOptionSelectorFilter("f1", "Options", false)
	var filterState scopes.FilterState
	// for RTM version of libunity-scopes we should see a log message
	reply.PushFilters([]scopes.Filter{filter1}, filterState)

	return qu.AddQueryResults(reply, query.QueryString(), settings)
}

type Settings struct {
	Steamworks bool `json:"steamworks"`
	MaxPrice   int  `json:"max_price"`
	Linux      bool `json:"linux"`
	Osx        bool `json:"osx"`
	Windows    bool `json:"windows"`
	Unknown    bool `json:"unknown"`
}

func (s *GameDealsScope) SetScopeBase(base *scopes.ScopeBase) {
	s.base = base
}

// DEPARTMENTS *****************************************************************

func (s *GameDealsScope) CreateDepartments(query *scopes.CannedQuery,
	metadata *scopes.SearchMetadata,
	reply *scopes.SearchReply) *scopes.Department {

	all, _ := scopes.NewDepartment("", query, "All Genres")
	for _, genre := range departments {
		department, _ := scopes.NewDepartment(genre[0], query, genre[1])
		all.AddSubdepartment(department)
	}

	return all
}

// MAIN ************************************************************************

func main() {
	if err := scopes.Run(&GameDealsScope{}); err != nil {
		log.Fatalln(err)
	}
}
