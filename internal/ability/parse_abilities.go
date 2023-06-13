package ability

import (
	"encoding/json"
	"ever-parse/internal/character"
	"ever-parse/internal/reference"
	"ever-parse/internal/specials"
	"ever-parse/internal/util"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gosimple/slug"
	"github.com/tidwall/gjson"
)

// BPUIAbilityMapping represents the relevant "Properties" inside a BP_UIAbility_* type.
type BPUIAbilityMapping struct {
	id                               string
	tecName                          string
	super                            *string
	AbilityIcon                      reference.ObjectReference
	AbilityName                      reference.PropertyReference
	AbilityDescription               reference.PropertyReference
	NextLevelPreviewText             reference.PropertyReference
	descriptionValuesFromCurveTables reference.CurveTableReference
	stance                           specials.Stance
	source                           string
	metapower                        bool
}

type Info struct {
	Id          string
	Name        string
	Description string
	Source      string
	Slot        string
	MetaPower   bool
	Stance      specials.Stance `json:"stance"`
	Properties  string
}

func (m BPUIAbilityMapping) GetNameProperty() reference.PropertyReference {
	return m.AbilityName
}

func (m BPUIAbilityMapping) GetDescriptionProperty() reference.PropertyReference {
	return m.AbilityDescription
}

func (m BPUIAbilityMapping) GetCurveProperty() reference.CurveTableReference {
	return m.descriptionValuesFromCurveTables
}

// ParseAbilities Parses hero abilities and writes to the abilities.json file
func ParseAbilities(root string, group *sync.WaitGroup) {
	abilities := make([]Info, 0)
	mappings := make([]BPUIAbilityMapping, 0)
	walkError := filepath.Walk(root, func(path string, info os.FileInfo, walkFnError error) error {
		//Shards are also a BP_UIAbility type/file , just stored in a folder called "Shards". Skip them
		if info.IsDir() && info.Name() == "Shards" {
			return filepath.SkipDir
		}

		//Accept all the BP_UIAbility_* files not located inside the "Shards" and create mappings
		if shouldParseFile(info.Name()) {
			err, abilityMapping := createBPUIAbilityMapping(path)
			if err != nil {
				println("Failed to parse: " + path)
				println("Error:" + err.Error())
				return nil
			}
			id := abilityId(path)
			abilityMapping.id = id
			abilityMapping.stance = GetStance(path)
			abilityMapping.source = reference.Source(path)
			abilityMapping.metapower = strings.Contains(strings.ToLower(path), "metapower")
			mappings = append(mappings, abilityMapping)
		}

		return nil
	})

	//Write file containing all the abilities
	util.Check(walkError)

	for _, abilityMapping := range mappings {
		var superMapping *BPUIAbilityMapping = nil
		if abilityMapping.super != nil {
			for _, other := range mappings {
				if other.tecName == *abilityMapping.super {
					superMapping = &other
					break
				}
			}
		}

		abilityInfo := Info{
			abilityMapping.id,
			GetAbilityName(abilityMapping, superMapping),
			util.ToValidHtml(GetAbilityDescription(abilityMapping, superMapping)),
			abilityMapping.source,
			parseAbilitySlot(abilityMapping, superMapping),
			abilityMapping.metapower,
			abilityMapping.stance,
			GetAbilityCurveValues(abilityMapping, superMapping),
		}
		//check if ability.info is inside array
		if !character.IsBlacklisted(abilityInfo.Source) {
			abilities = append(abilities, abilityInfo)
			//Copy the ability icon to the output folder
			reference.CopyImageFile(abilityMapping.AbilityIcon, abilityMapping.id, group, "abilities")
		}

		abilities = FixAbilityData(abilities)
	}

	walkError = util.WriteInfo("abilities.json", abilities)
	util.Check(walkError, "abilities.json", abilities)
}

// GetStance Retrieves the Stance for a certain ability identified by the path argument.
// The source (aka the hero it belongs to) is determined by reference.Source.
// If a special parser could be determined the special parser is used to create the specials.Stance value.
func GetStance(path string) specials.Stance {
	specialParser := specials.Parsers[slug.Make(reference.Source(path))]
	if specialParser != nil {
		return specialParser.GetStance(path)
	}
	return specials.AllStances
}

func shouldParseFile(name string) bool {
	if strings.HasPrefix(name, "BP_UIAbility") {
		return true
	}
	if strings.HasPrefix(name, "BP_MageTank_UIAbility_Mist.json") {
		return true
	}
	return false
}

// createBPUIAbilityMapping CParses hte "Properties" field inside a BP_UIAbility_* type and creates a mapping
func createBPUIAbilityMapping(path string) (error, BPUIAbilityMapping) {
	content, err := os.ReadFile(path)
	util.Check(err, path)
	var abilityRawJson = ""
	rawJSON := string(content)
	if strings.Contains(path, "BP_MageTank_UIAbility_Mist") {
		abilityRawJson = gjson.Get(rawJSON, "#(Type%\"BP_MageTank_UIAbility_Mist*\")#|0.Properties").String()
	} else {
		abilityRawJson = gjson.Get(rawJSON, "#(Type%\"BP_UIAbility*\")#|0.Properties").String()
	}

	var abilityMapping BPUIAbilityMapping
	err = json.Unmarshal([]byte(abilityRawJson), &abilityMapping)

	tecName := gjson.Get(rawJSON, "#(Name%BP_*UIAbility*_C)#|0.Name")
	if tecName.Type != gjson.Null {
		abilityMapping.tecName = tecName.String()
	}

	super := gjson.Get(rawJSON, "#(Super%BP_*UIAbility*_C)#|0.Super")
	if super.Type != gjson.Null {
		local := super.String()
		abilityMapping.super = &local
	}

	// Evaluate the field with GJson so we can determine if it's an array or not
	ctDescription := gjson.Get(abilityRawJson, "DescriptionValuesFromCurveTables")
	// But sometimes everything else already fails or no curve table values exist
	if err != nil || ctDescription.Type == gjson.Null {
		// In such a case we want to return early and let the calling function handle the error
		return err, abilityMapping
	}

	err, fixedCurveValues := reference.FixCurveTableValues(ctDescription)
	abilityMapping.descriptionValuesFromCurveTables = fixedCurveValues

	return err, abilityMapping
}

func abilityId(path string) string {
	var delimiter = ""
	if strings.Contains(path, "UIAbility_Mist.json") {
		delimiter = "UIAbility_"
	} else {
		delimiter = "BP_UIAbility_"
	}
	return slug.Make(reference.GenerateId(path, delimiter))
}

func parseAbilitySlot(m BPUIAbilityMapping, superMapping *BPUIAbilityMapping) string {
	if m.AbilityName.Key != "" {
		return m.AbilityName.Key[0:1]
	}
	if superMapping != nil && superMapping.AbilityName.Key != "" {
		return superMapping.AbilityName.Key[0:1]
	}

	return ""
}
