package character

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

func check(e error, path string) {
	if e != nil {
		println(path)
		panic(e)
	}
}

type Mapping struct {
	CharacterKitName          reference.PropertyReference
	CharacterKitDescription   reference.PropertyReference
	CharacterKitRole          reference.PropertyReference
	CharacterDefaultSkinImage reference.ImageReference
	CharacterPreviewImage     reference.ImageReference
	CharacterPortrait         reference.ImageReference
}

type Info struct {
	Id          string
	Name        string
	Description string
	Role        string
}

func ParseCharacters(root string) {
	characters := make([]Info, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(info.Name(), "BP_Player_") {
			content, err := os.ReadFile(path)
			check(err, path)
			characterRawJson := gjson.Get(string(content), "#(Type%\"BP_Player_*\")#|0.Properties").String()
			var characterMapping Mapping
			err = json.Unmarshal([]byte(characterRawJson), &characterMapping)
			check(err, path)
			if characterMapping.CharacterKitDescription.TableId != "" {
				name := reference.GetReferenceValue(characterMapping.CharacterKitName)
				id := slug.Make(name)
				characterInfo := Info{
					id,
					name,
					reference.GetReferenceValue(characterMapping.CharacterKitDescription),
					reference.GetReferenceValue(characterMapping.CharacterKitRole)}
				fmt.Println(path)
				fmt.Println(characterInfo)
				characters = append(characters, characterInfo)
				reference.CopyImageFile(characterMapping.CharacterPreviewImage, id+"_preview")
				reference.CopyImageFile(characterMapping.CharacterDefaultSkinImage, id+"_default")
				reference.CopyImageFile(characterMapping.CharacterPortrait, id+"_portrait")
			}
		}
		return nil
	})
	check(err, "")
	f, _ := os.Create("characters.json")
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	err = enc.Encode(characters)
	check(err, "")
	err = f.Close()
	check(err, "")
}
