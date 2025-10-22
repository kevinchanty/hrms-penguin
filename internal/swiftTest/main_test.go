package main

import (
	"slices"
	"testing"
)

func TestResetScore(t *testing.T) {
	scoreBoard := NewScoreBoard(&[]Player{
		{name: "Kimberley", score: 20},
		{name: "Carman", score: 0},
	})

	scoreBoard.ResetScore(0)

	for _, player := range *scoreBoard.players {
		want := 0
		get := player.score

		if want != get {
			t.Errorf("Want: %d, Got: %d", want, get)
		}
	}
}

func TestHighestScoreWin(t *testing.T) {
	scoreBoard := NewScoreBoard(&[]Player{
		{name: "Kimberley", score: 20},
		{name: "Carman", score: 0},
	})

	scoreBoard.doesHighestScoreWin = true
	scoreBoard.state = StateGameover

	want := []Player{{"Kimberley", 20}}
	get := *scoreBoard.Winners()

	if !slices.Equal(want, get) {
		t.Errorf("Want: %+v, Got: %+v", want, get)
	}
}

func TestLowestScoreWin(t *testing.T) {
	scoreBoard := NewScoreBoard(&[]Player{
		{name: "Kimberley", score: 20},
		{name: "Carman", score: 0},
	})

	scoreBoard.doesHighestScoreWin = false
	scoreBoard.state = StateGameover

	want := []Player{{"Carman", 0}}
	get := *scoreBoard.Winners()

	if !slices.Equal(want, get) {
		t.Errorf("Want: %+v, Got: %+v", want, get)
	}
}
