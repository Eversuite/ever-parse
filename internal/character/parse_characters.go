package character

import (
	"encoding/json"
	"ever-parse/internal/reference"
	"ever-parse/internal/util"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/tidwall/gjson"
	"os"
	"path/filepath"
	"strings"
)

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
			util.Check(err, path)
			characterRawJson := gjson.Get(string(content), "#(Type%\"BP_Player_*\")#|0.Properties").String()
			var characterMapping Mapping
			err = json.Unmarshal([]byte(characterRawJson), &characterMapping)
			util.Check(err, path)
			if characterMapping.CharacterKitDescription.TableId != "" {
				name := reference.GetReferenceValue(characterMapping.CharacterKitName)
				id := slug.Make(reference.CharacterId(path))
				characterInfo := Info{
					id,
					name,
					reference.GetReferenceValue(characterMapping.CharacterKitDescription),
					reference.GetReferenceValue(characterMapping.CharacterKitRole),
				}
				fmt.Println(path)
				fmt.Println(characterInfo)
				characters = append(characters, characterInfo)
				reference.CopyImageFile(characterMapping.CharacterPreviewImage, id+"-preview")
				reference.CopyImageFile(characterMapping.CharacterDefaultSkinImage, id+"-default")
				reference.CopyImageFile(characterMapping.CharacterPortrait, id+"-portrait")
			}
			if characterMapping.CharacterKitDescription.TableId == "" && characterMapping.CharacterKitDescription.SourceString != "" {
				name := characterMapping.CharacterKitName.SourceString
				id := slug.Make(reference.CharacterId(path))
				characterInfo := Info{
					id,
					name,
					characterMapping.CharacterKitDescription.SourceString,
					characterMapping.CharacterKitRole.SourceString,
				}
				fmt.Println(path)
				fmt.Println(characterInfo)
				characters = append(characters, characterInfo)
				reference.CopyImageFile(characterMapping.CharacterPreviewImage, id+"-preview")
				reference.CopyImageFile(characterMapping.CharacterDefaultSkinImage, id+"-default")
				reference.CopyImageFile(characterMapping.CharacterPortrait, id+"-portrait")
			}
		}
		return nil
	})
	util.Check(err)
	f, _ := os.Create("characters.json")
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	err = enc.Encode(characters)
	util.Check(err, characters)
	err = f.Close()
	util.Check(err, "Unable to close file")
}