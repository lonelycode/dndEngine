package game

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

type PlayerInSession struct {
	Character
	Initiative int
	IsNPC      bool
	LocationX  int
	LocationY  int
}

func (p *PlayerInSession) GetLocationX() int {
	return p.LocationX
}

func (p *PlayerInSession) GetLocationY() int {
	return p.LocationY
}

func (p *PlayerInSession) GetSize() int {
	return p.Size
}

func (p *PlayerInSession) NPC() bool {
	return p.IsNPC
}

type PhysicalObject interface {
	GetLocationX() int
	GetLocationY() int
	GetSize() int
	NPC() bool
}

type Players []*PlayerInSession

// Sort Sorts the players by initiative from highest to lowest, if initiative is equal, then by IsNPC
func (p Players) Sort() {
	sort.SliceStable(p, func(i, j int) bool {
		if p[i].Initiative == p[j].Initiative {
			return !p[i].IsNPC && p[j].IsNPC
		}
		return p[i].Initiative > p[j].Initiative
	})
}

type CombatSession struct {
	Players            Players
	CurrentPlayerIndex int
	TimeIndex          int
	Area               *CombatArea
}

type Structure struct {
	Name             string
	Description      string
	DifficultTerrain bool
	X                int
	Y                int
	Width            int
	Height           int
}

func (p *Structure) GetLocationX() int {
	return p.X
}

func (p *Structure) GetLocationY() int {
	return p.Y
}

func (p *Structure) GetSize() int {
	return p.Width * p.Height
}

func (p *Structure) NPC() bool {
	return false
}

type CombatArea struct {
	Length     int
	Width      int
	Structures []*Structure
}

// CheckCollision Returns true and the structure or player name if the given player is colliding with any structure or
// player in the area with their movement which is given by the parameters moveToX and MoveToY
func (cs *CombatSession) CheckCollision(playerIndex int, moveToX, moveToY int) (bool, string) {
	// Check if the player is colliding with any structure in the area
	for _, s := range cs.Area.Structures {
		if collide(moveToX, moveToY, cs.Players[playerIndex], s) {
			return true, s.Name
		}
	}

	// Check if the player is colliding with any other player in the area
	for i, p := range cs.Players {
		if i != playerIndex && collide(moveToX, moveToY, cs.Players[playerIndex], p) {
			return true, p.Name
		}
	}

	return false, ""
}

func collide(moveToX, moveToY int, p1, p2 PhysicalObject) bool {
	if moveToX == p2.GetLocationX() && moveToY == p2.GetLocationY() {
		return true
	}
	if p2.NPC() {
		return false // NPCs don't block movement of other players
	}
	dx := math.Abs(float64(moveToY - p2.GetLocationX()))
	dy := math.Abs(float64(moveToY - p2.GetLocationY()))
	if dx <= float64(p1.GetSize()+p2.GetSize())/2.0 && dy <= float64(p1.GetSize()+p2.GetSize())/2.0 {
		// The players are colliding
		return true
	}
	return false
}

// CheckDistance returns false if the player is moving to a position that is further away from the
// player than their speed property
func (cs *CombatSession) CheckDistance(playerIndex int, moveToX, moveToY int) bool {
	player := cs.Players[playerIndex]
	distance := math.Sqrt(math.Pow(float64(moveToX-player.GetLocationX()), 2) + math.Pow(float64(moveToY-player.GetLocationY()), 2))
	return distance <= float64(player.Speed.Current)
}

func (cs *CombatSession) NextPlayer() {
	cs.CurrentPlayerIndex++
	if cs.CurrentPlayerIndex >= len(cs.Players) {
		cs.CurrentPlayerIndex = 0
		cs.TimeIndex++
	}
}

func (cs *CombatSession) CurrentPlayer() *PlayerInSession {
	return cs.Players[cs.CurrentPlayerIndex]
}

func (cs *CombatSession) AddPlayer(p *PlayerInSession) {
	cs.Players = append(cs.Players, p)
	cs.Players.Sort()
}

// Move will move a player to the given location if the movement is valid, movement is valid if the player is not
// exceeding their Speed and is not colliding with another player or structure.
func (cs *CombatSession) Move(playerIndex, moveToX, moveToY int) error {
	player := cs.Players[playerIndex]

	// Check if the player is trying to exceed their speed
	if !cs.CheckDistance(playerIndex, moveToX, moveToY) {
		return errors.New("player is exceeding their speed")
	}

	// Check if the player is colliding with any other structure/player in the area
	if collision, name := cs.CheckCollision(playerIndex, moveToX, moveToY); collision {
		return fmt.Errorf("collision detected with %s", name)
	}

	// Calculate the distance travelled by the player
	distance := math.Sqrt(math.Pow(float64(moveToX-player.GetLocationX()), 2) + math.Pow(float64(moveToY-player.GetLocationY()), 2))

	// Update the player's current speed to subtract the distance travelled
	player.Speed.Current -= int(distance)

	// Move the player to the new position
	player.LocationX = moveToX
	player.LocationY = moveToY

	return nil
}
