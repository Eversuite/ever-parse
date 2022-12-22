package main

import (
	"ever-parse/internal/ability"
	"ever-parse/internal/character"
	"ever-parse/internal/shards"
)

func main() {
	ability.ParseAbilities(".")
	character.ParseCharacters(".")
	shards.ParseShards(".")
}