package shard

import (
	"ever-parse/internal/ability"
	"ever-parse/internal/reference"
	"ever-parse/internal/util"
	"github.com/gosimple/slug"
	"io/fs"
	"path/filepath"
	"strings"
)

func ParseShards(root string) {
	shards := make([]ability.Info, 0)

	err := filepath.WalkDir(root, dirWalker(&shards))
	util.Check(err, root, shards)

	err = util.WriteInfo("shard.json", shards)
	util.Check(err, "shard.json", shards)
}

func dirWalker(shards *[]ability.Info) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && strings.ToLower(d.Name()) == "shard" {
			return filepath.Walk(path, fileWalker(shards))
		}
		return err
	}
}

func fileWalker(shards *[]ability.Info) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, err error) error {
		if !strings.HasPrefix(info.Name(), "BP_UIAbility") {
			return err
		}

		err, abilityMapping := ability.CreateBPUIAbilityMapping(path)
		if err != nil {
			println("Failed to parse: " + path)
			println("Error:" + err.Error())
			return nil
		}
		id := ability.CreateId(path)
		shardInfo := ability.Info{
			Id:          id,
			Name:        reference.GetName(abilityMapping),
			Description: reference.GetDescription(abilityMapping),
			Source:      slug.Make(reference.Source(path)),
			Properties:  reference.GetCurveProperties(abilityMapping),
		}
		*shards = append(*shards, shardInfo)
		reference.CopyImageFile(abilityMapping.AbilityIcon, shardInfo.Id, "shard")
		return err
	}
}
