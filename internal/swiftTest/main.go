package main

type ScoreBoard struct {
	players             *[]Player
	state               GameState
	doesHighestScoreWin bool
}

func NewScoreBoard(players *[]Player) *ScoreBoard {
	if players != nil {
		return &ScoreBoard{
			players: players,
		}
	} else {
		return &ScoreBoard{
			players: &[]Player{
				{name: "Kirby", score: -50},
				{name: "Bella", score: 5},
			},
		}
	}
}
func (b *ScoreBoard) ResetScore(to int) {
	for i := range *b.players {
		(*b.players)[i].score = to
	}
}

func (b *ScoreBoard) Winners() *[]Player {
	if b.state != StateGameover {
		return &[]Player{}
	} else {
		winningScore := 0
		for _, player := range *b.players {
			if b.doesHighestScoreWin {
				if player.score > winningScore {
					winningScore = player.score
				}
			} else {
				if player.score < winningScore {
					winningScore = player.score
				}
			}
		}
		result := []Player{}
		for _, player := range *b.players {
			if player.score == winningScore {
				result = append(result, player)
			}
		}
		return &result
	}
}

type GameState int

const (
	StateSetup GameState = iota
	StatePlaying
	StateGameover
)

type Player struct {
	name  string
	score int
}

func main() {

}
