package main

import (
	"encoding/json"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/tidwall/gjson"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	ParseAbilities(".")
}

type PropertyReference struct {
	TableId      string
	Key          string
	SourceString string
}

type IconReference struct {
	ObjectPath string
}

type AbilityMapping struct {
	AbilityIcon IconReference
	AbilityName PropertyReference
	AbilityDescription   PropertyReference
	NextLevelPreviewText PropertyReference
}

type AbilityInfo struct {
	Id          string
	Name        string
	Description string
}

func ParseAbilities(root string) {
	abilities := make([]AbilityInfo, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(info.Name(), "BP_UIAbility") {
			content, err := os.ReadFile(path)
			check(err)
			//Parse the ability mappings
			abilityRawJson := gjson.Get(string(content), "#(Type%\"BP_UIAbility*\")#|0.Properties").String()
			var abilityMapping AbilityMapping
			err = json.Unmarshal([]byte(abilityRawJson), &abilityMapping)
			check(err)
			//fmt.Println(path)
			//Check that it uses references for its data
			if abilityMapping.AbilityName.TableId != "" {
				name := getAbilityReferenceValue(abilityMapping.AbilityName)
				id := slug.Make(name)
				abilityInfo := AbilityInfo{
					id,
					name,
					getAbilityReferenceValue(abilityMapping.AbilityDescription)}
				abilities = append(abilities, abilityInfo)
				copyImageFile(abilityMapping.AbilityIcon, id)
			}
			//Get skills which have strings directly in the abilities
			if abilityMapping.AbilityName.TableId == "" && abilityMapping.AbilityName.SourceString != "" {
				name := abilityMapping.AbilityName.SourceString
				id := slug.Make(name)
				abilityInfo := AbilityInfo{
					id,
					name,
					abilityMapping.AbilityDescription.SourceString,
				}
				fmt.Println(path)
				fmt.Println(abilityInfo)
				abilities = append(abilities, abilityInfo)
				copyImageFile(abilityMapping.AbilityIcon, id)
			}
		}
		return nil
	})
	check(err)

	f, _ := os.Create("test.json")
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	err = enc.Encode(abilities)
	check(err)
	err = f.Close()
	check(err)
}

func fixRoot(path string) string {
	return strings.ReplaceAll(path, "/Game/", "Game/")
}

func getAbilityReferenceValue(abilityReference PropertyReference) string {
	correctRoot := fixRoot(abilityReference.TableId)
	regex := regexp.MustCompile("\\..*")
	cleanedPath := regex.ReplaceAllString(correctRoot, ".json")
	content, err := os.ReadFile(cleanedPath)
	check(err)
	return gjson.Get(string(content), "#.StringTable.KeysToMetaData."+abilityReference.Key+"|0").String()
}

func copyImageFile(abilityIcon IconReference, id string) {
	correctRoot := fixRoot(abilityIcon.ObjectPath)
	regex := regexp.MustCompile("\\..*")
	cleanedPath := regex.ReplaceAllString(correctRoot, ".png")
	content, err := os.ReadFile(cleanedPath)
	check(err)
	err = os.WriteFile("./icons/"+id+".png", content, 0644)
	check(err)
}

type CharacterMapping struct {
	CharacterKitName PropertyReference
	CharacterKitDescription PropertyReference
}
func parseCharacters(root string) {
	characters := make([]AbilityInfo, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(info.Name(), "BP_Player_") {
			content, err := os.ReadFile(path)
			check(err)
			characterRawJson := gjson.Get(string(content), "#(Type%\"BP_Player_*\")#|0.Properties").String()

		}
	}
}
