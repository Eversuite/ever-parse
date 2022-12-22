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

func ParseAbilities(root string) {
	abilities := make([]Info, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
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
				getAbilityProps(abilityMapping),
			}
			abilities = append(abilities, abilityInfo)
			reference.CopyImageFile(abilityMapping.AbilityIcon, id)
		}
		return nil
	})
	util.Check(err)

	f, _ := os.Create("abilities.json")
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	err = enc.Encode(abilities)
	util.Check(err, abilities)
	err = f.Close()
	util.Check(err, f)
}

func getAbilityProps(abilityMapping Mapping) string {
	abilityProps := ""
	if abilityMapping.DescriptionValuesFromCurveTables != nil {
		bytes, err := json.Marshal(abilityMapping.DescriptionValuesFromCurveTables.GetValues())
		util.Check(err, abilityMapping)
		abilityProps = string(bytes)
	}
	return abilityProps
}