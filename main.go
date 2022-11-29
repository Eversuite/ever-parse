package main

import (
	"encoding/json"
	"github.com/gosimple/slug"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	FilePathWalkDir(".")
}

type AbilityReference struct {
	TableId string
	Key     string
}

type AbilityIcon struct {
	ObjectPath string
}

type AbilityMapping struct {
	AbilityIcon          AbilityIcon
	AbilityName          AbilityReference
	AbilityDescription   AbilityReference
	NextLevelPreviewText AbilityReference
}

type AbilityInfo struct {
	Id          string
	Name        string
	Description string
}

func FilePathWalkDir(root string) {
	abilities := make([]AbilityInfo, 0)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(info.Name(), "BP_UIAbility") {
			content, err := os.ReadFile(path)
			if err != nil {
				log.Fatal("Error when opening file: ", err)
			}
			//Parse the ability mappings
			abilityRawJson := gjson.Get(string(content), "#(Type%\"BP_UIAbility*\")#|0.Properties").String()
			var abilityMapping AbilityMapping
			json.Unmarshal([]byte(abilityRawJson), &abilityMapping)
			//fmt.Println(path)
			//Check that it uses references for it's data
			if abilityMapping.AbilityName.TableId != "" {
				name := getAbilityReferenceValue(abilityMapping.AbilityName)
				id := slug.Make(name)
				abilityInfo := AbilityInfo{
					id,
					name,
					getAbilityReferenceValue(abilityMapping.AbilityDescription)}
				abilities = append(abilities, abilityInfo)
			}
		}
		return nil
	})
	f, _ := os.Create("test.json")
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	enc.Encode(abilities)
	f.Close()
}

func getAbilityReferenceValue(abilityReference AbilityReference) string {
	correctRoot := strings.ReplaceAll(abilityReference.TableId, "/Game", "Game")
	regex := regexp.MustCompile("\\..*")
	cleanedPath := regex.ReplaceAllString(correctRoot, ".json")
	content, err := os.ReadFile(cleanedPath)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	return gjson.Get(string(content), "#.StringTable.KeysToMetaData."+abilityReference.Key+"|0").String()
}
