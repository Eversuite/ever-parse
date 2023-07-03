package shard

import (
	"encoding/json"
	"ever-parse/internal/reference"
	"ever-parse/internal/util"
	"github.com/gosimple/slug"
	"github.com/tidwall/gjson"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Mapping struct {
	ShardIcon            reference.ObjectReference   `json:"AbilityIcon"`
	ShardName            reference.PropertyReference `json:"AbilityName"`
	ShardDescription     reference.PropertyReference `json:"AbilityDescription"`
	ShardPreviewText     reference.PropertyReference `json:"NextLevelPreviewText"`
	shardCurveProperties reference.CurveTableReference
}

func (m Mapping) GetNameProperty() reference.PropertyReference {
	return m.ShardName
}

func (m Mapping) GetDescriptionProperty() reference.PropertyReference {
	return m.ShardDescription
}

func (m Mapping) GetCurveProperty() reference.CurveTableReference {
	return m.shardCurveProperties
}

type ShardInfo struct {
	Id          string
	Name        string
	Description string
	Type        string
	Source      string
	Properties  string
}

func ParseShards(root string, group *sync.WaitGroup) {
	shards := make([]ShardInfo, 0)

	err := filepath.WalkDir(root, dirWalker(&shards, group))
	util.Check(err, root, shards)

	shopShards := ParseShops(".")

	shards = populateTypeAndSource(&shards, shopShards)

	err = util.WriteInfo("shards.json", shards)
	util.Check(err, "shards.json", shards)

}

func dirWalker(shards *[]ShardInfo, group *sync.WaitGroup) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && strings.ToLower(d.Name()) == "shards" {
			return filepath.Walk(path, fileWalker(shards, group))
		}
		return err
	}
}

func fileWalker(shards *[]ShardInfo, group *sync.WaitGroup) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, err error) error {
		if !strings.HasPrefix(info.Name(), "BP_UIAbility") {
			return err
		}

		content, err := os.ReadFile(path)
		util.Check(err, path)
		//Parse the ability mappings
		abilityRawJson := gjson.Get(string(content), "#(Type%\"BP_UIAbility*\")#|0.Properties").String()
		var abilityMapping Mapping
		err = json.Unmarshal([]byte(abilityRawJson), &abilityMapping)

		ctDescription := gjson.Get(abilityRawJson, "DescriptionValuesFromCurveTables")
		err, tableReference := reference.FixCurveTableValues(ctDescription)
		abilityMapping.shardCurveProperties = tableReference

		if err != nil {
			println("Failed to parse: " + path)
			return nil
		}
		util.Check(err, path)
		id := slug.Make(shardId(path))
		shardInfo := ShardInfo{
			id,
			reference.GetName(abilityMapping),
			util.ToValidHtml(reference.GetDescription(abilityMapping)),
			"",
			"Random drop",
			reference.GetCurveProperties(abilityMapping),
		}

		if !IsBlacklisted(shardInfo.Id) {
			*shards = append(*shards, shardInfo)
			reference.CopyImageFile(abilityMapping.ShardIcon, shardInfo.Id, group, "shard")
		}
		return err
	}
}

func shardId(path string) string {
	delimiter := "BP_UIAbility_"
	return reference.GenerateId(path, delimiter)
}

func populateTypeAndSource(shardInfo *[]ShardInfo, ShardShopInfo []ShopInfo) []ShardInfo {
	for i, shard := range *shardInfo {
		for _, shopShard := range ShardShopInfo {
			if shard.Id == shopShard.Id {
				(*shardInfo)[i].Type = shopShard.Type
				(*shardInfo)[i].Source = shopShard.Source
			}
		}
	}
	return *shardInfo
}
