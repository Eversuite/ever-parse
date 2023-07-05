package main

import (
	"ever-parse/internal/ability"
	"ever-parse/internal/character"
	"ever-parse/internal/consumable"
	"ever-parse/internal/reference"
	"ever-parse/internal/shard"
	"ever-parse/internal/talent"
	"ever-parse/internal/util"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	cleanupPreviousRun()
	group := &sync.WaitGroup{}

	parallelize(group, func() { ability.ParseAbilities(".", group) })
	parallelize(group, func() { character.ParseCharacters(".", group) })
	parallelize(group, func() { shard.ParseShards(".", group) })
	parallelize(group, func() { talent.ParseTalents(".", group) })
	parallelize(group, func() { talent.ParseTalentTrees(".") })
	parallelize(group, func() { consumable.ParseConsumables(".", group) })

	group.Wait()

	fmt.Printf("Run took: %+v\n", time.Now().Sub(start))

}

func parallelize(group *sync.WaitGroup, f func()) {
	group.Add(1)
	go func() {
		defer group.Done()
		f()
	}()

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
