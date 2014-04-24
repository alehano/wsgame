package main

import (
	"github.com/alehano/wsgame/game"
	"github.com/alehano/wsgame/utils"
	"log"
)

var allRooms = make(map[string]*room)
var freeRooms = make(map[string]*room)
var roomsCount int

type room struct {
	name string

	// Registered connections.
	playerConns map[*playerConn]bool

	// Update state for all conn.
	updateAll chan bool

	// Register requests from the connections.
	join chan *playerConn

	// Unregister requests from connections.
	leave chan *playerConn
}

// Run the room in goroutine
func (r *room) run() {
	for {
		select {
		case c := <-r.join:
			r.playerConns[c] = true
			r.updateAllPlayers()

			// if room is full - delete from freeRooms
			if len(r.playerConns) == 2 {
				delete(freeRooms, r.name)
				// pair players
				var p []*game.Player
				for k, _ := range r.playerConns {
					p = append(p, k.Player)
				}
				game.PairPlayers(p[0], p[1])
			}

		case c := <-r.leave:
			c.GiveUp()
			r.updateAllPlayers()
			delete(r.playerConns, c)
			if len(r.playerConns) == 0 {
				goto Exit
			}
		case <-r.updateAll:
			r.updateAllPlayers()
		}
	}

Exit:

	// delete room
	delete(allRooms, r.name)
	delete(freeRooms, r.name)
	roomsCount -= 1
	log.Print("Room closed:", r.name)
}

func (r *room) updateAllPlayers() {
	for c := range r.playerConns {
		c.sendState()
	}
}

func NewRoom(name string) *room {
	if name == "" {
		name = utils.RandString(16)
	}

	room := &room{
		name:        name,
		playerConns: make(map[*playerConn]bool),
		updateAll:   make(chan bool),
		join:        make(chan *playerConn),
		leave:       make(chan *playerConn),
	}

	allRooms[name] = room
	freeRooms[name] = room

	// run room
	go room.run()

	roomsCount += 1

	return room
}
