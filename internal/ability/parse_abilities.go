package ability

import (
	"encoding/json"
	"ever-parse/internal/reference"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/tidwall/gjson"
	"os"
	"path/filepath"
	"strings"
)

type Mapping struct {
	AbilityIcon          reference.ImageReference
	AbilityName          reference.PropertyReference
	AbilityDescription   reference.PropertyReference
	NextLevelPreviewText reference.PropertyReference
}

type Info struct {
	Id          string
	Name        string
	Description string
	Source      string
}

func ParseAbilities(root string) {
	abilities := make([]Info, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(info.Name(), "BP_UIAbility") {
			content, err := os.ReadFile(path)
			check(err, path)
			//Parse the ability mappings
			abilityRawJson := gjson.Get(string(content), "#(Type%\"BP_UIAbility*\")#|0.Properties").String()
			var abilityMapping Mapping
			err = json.Unmarshal([]byte(abilityRawJson), &abilityMapping)
			if err != nil {
				println("Failed to parse: " + path)
				return nil
			}
			check(err, path)
			//fmt.Println(path)
			//Check that it uses references for its data
			if abilityMapping.AbilityName.TableId != "" {
				name := reference.GetReferenceValue(abilityMapping.AbilityName)
				id := slug.Make(reference.AbilityId(path))
				abilityInfo := Info{
					id,
					name,
					reference.GetReferenceValue(abilityMapping.AbilityDescription),
					slug.Make(reference.AbilitySource(path)),
				}
				abilities = append(abilities, abilityInfo)
				reference.CopyImageFile(abilityMapping.AbilityIcon, id)
			}
			//Get skills which have strings directly in the abilities
			if abilityMapping.AbilityName.TableId == "" && abilityMapping.AbilityName.SourceString != "" {
				name := abilityMapping.AbilityName.SourceString
				id := slug.Make(reference.AbilityId(path))
				abilityInfo := Info{
					id,
					name,
					abilityMapping.AbilityDescription.SourceString,
					slug.Make(reference.AbilitySource(path)),
				}
				fmt.Println(path)
				fmt.Println(abilityInfo)
				abilities = append(abilities, abilityInfo)
				reference.CopyImageFile(abilityMapping.AbilityIcon, id)
			}
		}
		return nil
	})
	check(err, "")

	f, _ := os.Create("abilities.json")
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	err = enc.Encode(abilities)
	check(err, "")
	err = f.Close()
	check(err, "")
}

func check(e error, path string) {
	if e != nil {
		println(path)
		panic(e)
	}
}
