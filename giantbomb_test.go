package main

import (
	"testing"
)

func TestGames(t *testing.T) {
	var gb GiantBomb

	gb.SetCacheDirectory("cache")

	info, _ := gb.GetInfo("xcom")

	if info.Name == "" {
		t.Log("Name is empty")
		t.Fail()
	}
	if len(info.Platforms) < 1 {
		t.Log("platform is empty")
		t.Fail()
	}
}
