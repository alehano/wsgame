package game

import (
	"log"
)

type Player struct {
	Name  string
	Enemy *Player
}

func NewPlayer(name string) *Player {
	player := &Player{Name: name}
	return player
}

func PairPlayers(p1 *Player, p2 *Player) {
	p1.Enemy, p2.Enemy = p2, p1
}

func (p *Player) Command(command string) {

	log.Print("Command: '", command, "' received by player: ", p.Name)
}

func (p *Player) GetState() string {
	return "Game state for Player: " + p.Name
}

func (p *Player) GiveUp() {
	log.Print("Player gave up: ", p.Name)
}
