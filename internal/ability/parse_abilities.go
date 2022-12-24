package ability

import (
	"encoding/json"
	"ever-parse/internal/reference"
	"ever-parse/internal/util"
	"github.com/gosimple/slug"
	"github.com/tidwall/gjson"
	"os"
	"path/filepath"
	"strings"
)

type Mapping struct {
	AbilityIcon                      reference.ImageReference
	AbilityName                      reference.PropertyReference
	AbilityDescription               reference.PropertyReference
	NextLevelPreviewText             reference.PropertyReference
	DescriptionValuesFromCurveTables reference.CurveTableReference
}

type Info struct {
	Id          string
	Name        string
	Description string
	Source      string
	Properties  string
}

func (m Mapping) GetNameProperty() reference.PropertyReference {
	return m.AbilityName
}

func (m Mapping) GetDescriptionProperty() reference.PropertyReference {
	return m.AbilityDescription
}

func (m Mapping) GetCurveProperty() reference.CurveTableReference {
	return m.DescriptionValuesFromCurveTables
}

func ParseAbilities(root string) {
	abilities := make([]Info, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == "Shards" {
			return filepath.SkipDir
		}
		if strings.HasPrefix(info.Name(), "BP_UIAbility") {
			content, err := os.ReadFile(path)
			util.Check(err, path)
			//Parse the ability mappings
			abilityRawJson := gjson.Get(string(content), "#(Type%\"BP_UIAbility*\")#|0.Properties").String()
			var abilityMapping Mapping
			err = json.Unmarshal([]byte(abilityRawJson), &abilityMapping)
			if err != nil {
				println("Failed to parse: " + path)
				return nil
			}
			util.Check(err, path)
			id := slug.Make(reference.AbilityId(path))
			abilityInfo := Info{
				id,
				reference.GetName(abilityMapping),
				reference.GetDescription(abilityMapping),
				slug.Make(reference.AbilitySource(path)),
				reference.GetCurveProperties(abilityMapping),
			}

			abilities = append(abilities, abilityInfo)

			reference.CopyImageFile(abilityMapping.AbilityIcon, id, "abilities")
		}
		return nil
	})
	util.Check(err)

	err = util.WriteInfo("abilities.json", abilities)
	util.Check(err, "abilities.json", abilities)

}