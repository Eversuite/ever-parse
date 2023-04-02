package ability

import (
	"encoding/json"
	"ever-parse/internal/reference"
	"ever-parse/internal/util"
	"os"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
	"github.com/tidwall/gjson"
)

// BPUIAbilityMapping represents the relevant "Properties" inside a BP_UIAbility_* type.
type BPUIAbilityMapping struct {
	AbilityIcon                      reference.ObjectReference
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
	Slot        string
	Properties  string
}

func (m BPUIAbilityMapping) GetNameProperty() reference.PropertyReference {
	return m.AbilityName
}

func (m BPUIAbilityMapping) GetDescriptionProperty() reference.PropertyReference {
	return m.AbilityDescription
}

func (m BPUIAbilityMapping) GetCurveProperty() reference.CurveTableReference {
	return m.DescriptionValuesFromCurveTables
}

// ParseAbilities Parses hero abilities and writes to the abilities.json file
func ParseAbilities(root string) {
	abilities := make([]Info, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		//Shards are also a BP_UIAbility type/file , just stored in a folder called "Shards". Skip them
		if info.IsDir() && info.Name() == "Shards" {
			return filepath.SkipDir
		}
		//Accept all the BP_UIAbility_* files not located inside the "Shards" and create mappings
		if strings.HasPrefix(info.Name(), "BP_UIAbility") {
			err, abilityMapping := createBPUIAbilityMapping(path)
			if err != nil {
				println("Failed to parse: " + path)
				println("Error:" + err.Error())
				return nil
			}

			id := abilityId(path)
			abilityInfo := Info{
				id,
				reference.GetName(abilityMapping),
				util.ToValidHtml(reference.GetDescription(abilityMapping)),
				reference.Source(path),
				parseAbilitySlot(abilityMapping.AbilityName),
				reference.GetCurveProperties(abilityMapping),
			}
			//check if ability.info is inside array
			if util.IsHeroWhitelisted(abilityInfo.Source) {
				abilities = append(abilities, abilityInfo)
				//Copy the ability icon to the output folder
				reference.CopyImageFile(abilityMapping.AbilityIcon, id, "abilities")
			}
		}
		return nil
	})
	//Write file containing all teh talents
	util.Check(err)
	err = util.WriteInfo("abilities.json", abilities)
	util.Check(err, "abilities.json", abilities)
}

// createBPUIAbilityMapping CParses hte "Properties" field inside a BP_UIAbility_* type and creates a mapping
func createBPUIAbilityMapping(path string) (error, BPUIAbilityMapping) {
	content, err := os.ReadFile(path)
	util.Check(err, path)
	abilityRawJson := gjson.Get(string(content), "#(Type%\"BP_UIAbility*\")#|0.Properties").String()
	var abilityMapping BPUIAbilityMapping
	err = json.Unmarshal([]byte(abilityRawJson), &abilityMapping)
	return err, abilityMapping
}

func abilityId(path string) string {
	delimiter := "BP_UIAbility_"
	return slug.Make(reference.GenerateId(path, delimiter))
}

func parseAbilitySlot(abilityReference reference.PropertyReference) string {
	if abilityReference.Key == "" {
		return ""
	}
	return abilityReference.Key[0:1]
}
