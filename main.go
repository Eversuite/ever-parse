package main

import (
	"ever-parse/internal/ability"
	"ever-parse/internal/character"
	"ever-parse/internal/reference"
	"ever-parse/internal/shards"
	"ever-parse/internal/util"
	"os"
	"path/filepath"
)

func main() {
	cleanupPreviousRun()
	ability.ParseAbilities(".")
	character.ParseCharacters(".")
	shards.ParseShards(".")
}

func cleanupPreviousRun() {
	allPaths := [][]string{
		{".", util.ParsedDataDir},
		{".", reference.ProjectVImagePath},
	}

	for _, paths := range allPaths {
		dir := filepath.Join(paths...)
		err := os.RemoveAll(dir)
		util.Check(err, "Unable to remove ", dir)
	}
}