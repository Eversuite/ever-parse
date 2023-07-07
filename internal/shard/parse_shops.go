package shard

import (
	"encoding/json"
	"ever-parse/internal/reference"
	"ever-parse/internal/util"
	"os"
	"path/filepath"
	"strings"
)

type Category struct {
	Key struct {
		TagName string `json:"TagName"`
	} `json:"Key"`
	Value struct {
		CategoryLabel   Label      `json:"CategoryLabel"`
		ShardCollection []Template `json:"ShardCollection"`
	} `json:"Value"`
}

type Template struct {
	ObjectName string `json:"ObjectName"`
	ObjectPath string `json:"ObjectPath"`
}

type Label struct {
	TableId string `json:"TableId"`
	Key     string `json:"Key"`
}

type Collection struct {
	Type       string `json:"Type"`
	Name       string `json:"Name"`
	Class      string `json:"Class"`
	Properties struct {
		ShardCategories []Category `json:"ShardCategories"`
	} `json:"Properties"`
}

type ShopInfo struct {
	Id     string
	Type   string
	Source string
}

func ParseShops(root string) []ShopInfo {
	shopShards := make([]ShopInfo, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(info.Name(), "ShardCollection") {
			content, err := os.ReadFile(path)
			util.Check(err, path)
			rawJSON := string(content)
			var shardCollections []Collection
			err = json.Unmarshal([]byte(rawJSON), &shardCollections)
			util.Check(err, path)
			for _, shardCollection := range shardCollections {
				for _, shardCategory := range shardCollection.Properties.ShardCategories {
					for _, shardTemplate := range shardCategory.Value.ShardCollection {
						if strings.HasPrefix(info.Name(), "DA_TownVendorShardCollection") {
							shardShopInfo := ShopInfo{
								generateId(shardTemplate.ObjectPath),
								getShardType(shardCategory.Key.TagName),
								util.KebabToTitle(generateId(path)),
							}
							shopShards = append(shopShards, shardShopInfo)
						} else if strings.HasPrefix(info.Name(), "DA_ShadowVendorShardCollection") {
							shardShopInfo := ShopInfo{
								generateId(shardTemplate.ObjectPath),
								getShardType(shardCategory.Key.TagName),
								"Secret Shop",
							}
							shopShards = append(shopShards, shardShopInfo)
						}

					}
				}
			}
		}
		return nil
	})
	util.Check(err, root)
	return shopShards
}

func generateId(s string) string {
	start := strings.Index(s, "_") + 1
	end := strings.LastIndex(s, ".")
	if start > 0 && end > start {
		return reference.GenerateId(s[start:end], "_")
	}
	return ""
}

func getShardType(s string) string {
	parts := strings.Split(s, ".")
	if len(parts) > 0 {
		return strings.ToLower(parts[len(parts)-1])
	}
	return ""
}
