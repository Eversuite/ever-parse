package talent

import (
	"encoding/json"
	"ever-parse/internal/reference"
	"ever-parse/internal/util"
	"os"
	"path/filepath"
	"strings"

	"github.com/tidwall/gjson"
)

// VVMetaPowerCategoryMapping represents the relevant "Properties" inside a VVMetaPowerCategoryMapping type.
type VVMetaPowerCategoryMapping struct {
	MetaPowerCategoryName        reference.PropertyReference
	MetaPowerCategoryDescription reference.PropertyReference
}

func (m VVMetaPowerCategoryMapping) GetNameProperty() reference.PropertyReference {
	return m.MetaPowerCategoryName
}

func (m VVMetaPowerCategoryMapping) GetDescriptionProperty() reference.PropertyReference {
	return m.MetaPowerCategoryDescription
}

type TreeInfo struct {
	Id          string
	Name        string
	Description string
	Source      string
}

// ParseTalentTrees Parses hero talent trees and writes to the talent-trees.json file
func ParseTalentTrees(root string) {
	talentTrees := make([]TreeInfo, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		//Accept all the MPA_* files and create mappings
		if strings.HasPrefix(info.Name(), "MPA_") {
			err, talentTreeMapping := createVVMetaPowerCategoryMapping(path)
			if err != nil {
				println("Failed to parse: " + path)
				println("Error:" + err.Error())
				return nil
			}
			source := reference.Source(path)
			id := GenerateTalentCategoryId(source, path)
			if util.IsHeroWhitelisted(source) {
				talentTrees = append(talentTrees, TreeInfo{
					Id:          id,
					Name:        reference.GetName(talentTreeMapping),
					Description: reference.GetDescription(talentTreeMapping),
					Source:      source,
				})
			}

		}
		return nil
	})
	util.Check(err)
	err = util.WriteInfo("talent-tree.json", talentTrees)
	util.Check(err, "talent-tree.json", talentTrees)
}

func createVVMetaPowerCategoryMapping(path string) (error, VVMetaPowerCategoryMapping) {
	content, err := os.ReadFile(path)
	util.Check(err, path)
	talentTreeRawJson := gjson.Get(string(content), "#(Type%\"VVMetaPowerCategory\")#|0.Properties").String()
	var talentTreeMapping VVMetaPowerCategoryMapping
	err = json.Unmarshal([]byte(talentTreeRawJson), &talentTreeMapping)
	return err, talentTreeMapping
}
