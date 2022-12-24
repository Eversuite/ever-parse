package character

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
	CharacterKitName          reference.PropertyReference
	CharacterKitDescription   reference.PropertyReference
	CharacterKitRole          reference.PropertyReference
	CharacterDefaultSkinImage reference.ImageReference
	CharacterPreviewImage     reference.ImageReference
	CharacterPortrait         reference.ImageReference
}

func (m Mapping) GetNameProperty() reference.PropertyReference {
	return m.CharacterKitName
}

func (m Mapping) GetDescriptionProperty() reference.PropertyReference {
	return m.CharacterKitDescription
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
			id := slug.Make(reference.CharacterId(path))
			characterInfo := Info{
				id,
				reference.GetName(characterMapping),
				reference.GetDescription(characterMapping),
				characterMapping.getRole(),
			}
			characters = append(characters, characterInfo)
			reference.CopyImageFile(characterMapping.CharacterPreviewImage, id+"-preview", "characters", "preview")
			reference.CopyImageFile(characterMapping.CharacterDefaultSkinImage, id+"-default", "characters", "default-skin")
			reference.CopyImageFile(characterMapping.CharacterPortrait, id+"-portrait", "characters", "portraits")
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

func (m Mapping) getRole() string {
	if m.CharacterKitRole.TableId != "" {
		return reference.GetReferenceValue(m.CharacterKitRole)
	}

	if m.CharacterKitRole.SourceString != "" {
		return m.CharacterKitRole.SourceString
	}

	return "UnknownRoleProperty"
}