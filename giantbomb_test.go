package main

import (
	"log"
	"testing"
)

func TestGames(t *testing.T) {
	var gb GiantBomb

	info := gb.GetInfo("xcom")

	if info.Name == "" {
		t.Log("Name is empty")
		t.Fail()
	}
	log.Println(info.Description)
}
