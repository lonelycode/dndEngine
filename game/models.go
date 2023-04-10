package game

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Stat struct {
	Value    int
	Modifier int
}

type Speed struct {
	Max     int
	Current int
	History []*SpeedModifier
}

type SpeedModifier struct {
	TimeStamp time.Time
	Amount    int
	From      string
}

type HPModifier struct {
	TimeStamp time.Time
	Amount    int
	From      string
}

type HitPoints struct {
	Max     int
	Extra   int
	Current int
	History []*HPModifier
}

type Character struct {
	Name      string
	HitPoints *HitPoints
	Speed     *Speed
	Str       *Stat
	Dex       *Stat
	Con       *Stat
	Int       *Stat
	Wis       *Stat
	Cha       *Stat
	Senses    []*MetaModifier
	Skills    []*MetaModifier
	Languages []string
	Size      int // square feet
}

type DiceRoll rune
type DiceType int

const (
	Neutral      DiceRoll = 0
	Advantage    DiceRoll = 1
	Disadvantage DiceRoll = 2

	D4  DiceType = 4
	D6  DiceType = 6
	D8  DiceType = 8
	D10 DiceType = 10
	D12 DiceType = 12
	D20 DiceType = 20
)

type MetaModifier struct {
	Name      string
	Value     int
	Range     int
	RangeType string
	Targets   int
	RollWith  DiceRoll
}

type Save struct {
	Stat    string // "str", "dex", "con", "int", "wis", "cha"
	Value   int
	Outcome string
}

type Action struct {
	Name        string
	Description string
	AttackType  string
	OnHit       []Dice
	Save        *Save
	HitModifier *MetaModifier
}

type Dice struct {
	NumDice int
	Dice    DiceType
	Mod     int
}

func (d *Dice) Equals(other *Dice) bool {
	if d.NumDice != other.NumDice {
		return false
	}
	if d.Dice != other.Dice {
		return false
	}
	if d.Mod != other.Mod {
		return false
	}
	return true
}

func ParseDiceRoll(s string) (*Dice, error) {
	parts := strings.Split(s, "d")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid dice roll format")
	}

	numDice, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid dice count")
	}

	mod := 0
	diceParts := strings.Split(parts[1], "+")
	if len(diceParts) > 1 {
		mod, err = strconv.Atoi(diceParts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid modifier")
		}
	}

	diceInt, err := strconv.Atoi(diceParts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid dice type")
	}

	var dice DiceType
	switch diceInt {
	case int(D4):
		dice = D4
	case int(D6):
		dice = D6
	case int(D8):
		dice = D8
	case int(D10):
		dice = D10
	case int(D12):
		dice = D12
	case int(D20):
		dice = D20
	default:
		return nil, fmt.Errorf("invalid dice type: %d", diceInt)
	}

	return &Dice{NumDice: numDice, Dice: dice, Mod: mod}, nil
}
