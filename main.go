package main

import (
	"ever-parse/internal/ability"
	"ever-parse/internal/character"
	"ever-parse/internal/reference"
	"ever-parse/internal/shard"
	"ever-parse/internal/talent"
	"ever-parse/internal/util"
	"os"
	"path/filepath"
)

func main() {
	cleanupPreviousRun()
	ability.ParseAbilities(".")
	character.ParseCharacters(".")
	shard.ParseShards(".")
	talent.ParseTalents(".")

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
