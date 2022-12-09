package main

import (
	"ever-parse/internal/ability"
	"ever-parse/internal/character"
)

func main() {
	ability.ParseAbilities(".")
	character.ParseCharacters(".")
}
