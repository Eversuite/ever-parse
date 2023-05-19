package specials

import (
	"os"
	"strings"
)

// Stance A stance is an extra ability set of a hero. A hero can contain a certain amount of abilities
// per set. An ability may belong to several different stances. A stance is represented by the active bit
// in the unsigned integer.
type Stance uint8

// AllStances In this case every bit is 1 aka 'active' hence the ability belongs to all stances.
const AllStances = Stance(^uint8(0))

// Constants for the mini-tank stances
const (
	mtFighter Stance = 1 << iota
	mtTank
)

// mtStances Are the two possible stances of mini-tank
var mtStances = map[string]Stance{
	"FighterStance": mtFighter,
	"TankStance":    mtTank,
}

// StanceParser A StanceParser is a common struct that can be used to parse any
// stance specific data. It has a public stance map that maps name->Stance and a private
// function to determine which file belongs to which stance. (See GetStance)
type StanceParser struct {
	StanceMap           map[string]Stance
	determineStanceFunc func(path string) Stance
}

// GetStance Resolves the stance for the current file.
// It is expected that the returned value is > 0 and <= AllStances
func (m *StanceParser) GetStance(path string) Stance {
	return m.determineStanceFunc(path)
}

// StanceResolver Returns a map[Stance]string whereas the key is the number of the stance and the associated value
// is the name of the stance. This is used to retrieve the name of the stance a certain ability belongs to.
func (m *StanceParser) StanceResolver() map[Stance]string {
	res := make(map[Stance]string, 0)
	for k, v := range m.StanceMap {
		res[v] = k
	}
	return res
}

func newMiniTankParser() *StanceParser {
	return &StanceParser{
		StanceMap: mtStances,
		determineStanceFunc: func(path string) Stance {
			directories := strings.Split(path, string(os.PathSeparator))
			stanceName := directories[4]
			selected := mtStances[stanceName]

			if selected == 0 {
				return AllStances
			}

			return selected
		},
	}
}

const (
	Water Stance = 1 << iota
	Lightning
)

var cmStances = map[string]Stance{
	"WaterForm":     Water,
	"LightningForm": Lightning,
}

func newControlMageParser() *StanceParser {
	return &StanceParser{
		StanceMap: cmStances,
		determineStanceFunc: func(path string) Stance {
			directories := strings.Split(path, string(os.PathSeparator))
			stanceName := directories[4]
			selected := cmStances[stanceName]
			if selected == 0 {
				return AllStances
			}

			return selected
		},
	}
}

const (
	RhDamage Stance = 1 << iota
	RhHealing
)

var rhStances = map[string]Stance{
	"Damage":  RhDamage,
	"Healing": RhHealing,
}

func newRangeHealerParser() *StanceParser {
	return &StanceParser{
		StanceMap: rhStances,
		determineStanceFunc: func(path string) Stance {
			index := strings.LastIndex(path, "_")
			fileSuffix := strings.LastIndex(path, ".json")
			if index < 0 || index >= len(path) {
				return AllStances
			}
			stanceName := path[index+1 : fileSuffix]
			selected := rhStances[stanceName]
			if selected == 0 {
				return AllStances
			}
			return selected
		},
	}
}

// Parsers Global variable to retrieve a set of specific parsers for certain characters.
// Using pointers to ensure only a single instance is in circulation.
// Since no other module should ever need to create their own instance of a character-specific parser
// the factory functions have been made private to this module.
var Parsers = map[string]*StanceParser{
	"mini-tank":     newMiniTankParser(),
	"control-mage":  newControlMageParser(),
	"ranged-healer": newRangeHealerParser(),
}