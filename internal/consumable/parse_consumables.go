package consumable

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
	ConsumableIcon        reference.ObjectReference   `json:"AbilityIcon"`
	ConsumableName        reference.PropertyReference `json:"AbilityName"`
	ConsumableDescription reference.PropertyReference `json:"AbilityDescription"`
}

func (m Mapping) GetNameProperty() reference.PropertyReference {
	return m.ConsumableName
}

func (m Mapping) GetDescriptionProperty() reference.PropertyReference {
	return m.ConsumableDescription
}

type ConsumableInfo struct {
	Id          string
	Name        string
	Description string
}

func ParseConsumables(root string, group *sync.WaitGroup) {
	consumables := make([]ConsumableInfo, 0)

	err := filepath.WalkDir(root, dirWalker(&consumables, group))
	util.Check(err, root, consumables)

	err = util.WriteInfo("consumables.json", consumables)
	util.Check(err, "consumables.json", consumables)
}

func dirWalker(consumables *[]ConsumableInfo, group *sync.WaitGroup) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && strings.ToLower(d.Name()) == "consumables" {
			return filepath.Walk(path, fileWalker(consumables, group))
		}
		return err
	}
}

func fileWalker(consumables *[]ConsumableInfo, group *sync.WaitGroup) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, err error) error {
		if !strings.HasPrefix(info.Name(), "BP_UIAbility") {
			return err
		}

		content, err := os.ReadFile(path)
		util.Check(err, path)
		//Parse the ability mappings
		consumableRawJson := gjson.Get(string(content), "#(Type%\"BP_UIAbility*\")#|0.Properties").String()
		var consumableMapping Mapping
		err = json.Unmarshal([]byte(consumableRawJson), &consumableMapping)
		if err != nil {
			println("Failed to parse: " + path)
			return nil
		}
		util.Check(err, path)
		id := slug.Make(consumableId(path))
		consumableInfo := ConsumableInfo{
			id,
			reference.GetName(consumableMapping),
			util.ToValidHtml(reference.GetDescription(consumableMapping)),
		}
		*consumables = append(*consumables, consumableInfo)
		reference.CopyImageFile(consumableMapping.ConsumableIcon, consumableInfo.Id, group, "shard")
		return err
	}
}

func consumableId(path string) string {
	delimiter := "BP_UIAbility_"
	return reference.GenerateId(path, delimiter)
}
