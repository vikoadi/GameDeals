package main

import (
	"testing"
)

func TestGames(t *testing.T) {
	var gb GiantBomb

	gb.SetCacheDirectory("cache")

	info, _ := gb.GetInfo("XCOM 2")

	if info.Name == "" {
		t.Log("Name is empty")
		t.Fail()
	}

	if info.Name != "XCOM 2" {
		t.Log("Name is not precise")
		t.Fail()
	}

	if len(info.Platforms) < 1 {
		t.Log("platform is empty")
		t.Fail()
	}

	if _, err := gb.GetInfo("unknown games"); err == nil {
		t.Log("shouldnt return any game")
		t.Fail()
	}

	if info, err := gb.GetInfo("worms 2 armageddon"); err == nil {
		if info.Name != "worms 2 armageddon" {
			t.Log("shouldnt return exact name")
			t.Fail()
		}
	}
}
