package character

import (
	"encoding/json"
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

// BPPlayerMapping represents the relevant "Properties" inside BP_Player_* types.
type BPPlayerMapping struct {
	CharacterKitName          reference.PropertyReference
	CharacterKitDescription   reference.PropertyReference
	CharacterKitRole          reference.PropertyReference
	CharacterDefaultSkinImage reference.ObjectReference
	CharacterPreviewImage     reference.ObjectReference
	CharacterPortrait         reference.ObjectReference
}

func (m BPPlayerMapping) GetNameProperty() reference.PropertyReference {
	return m.CharacterKitName
}

func (m BPPlayerMapping) GetDescriptionProperty() reference.PropertyReference {
	return m.CharacterKitDescription
}

type Info struct {
	Id            string
	Name          string
	Description   string
	Role          string
	StanceMapJson *string `json:"stance_map"`
}

// ParseCharacters Parses heroes and writes to the heroes.json file
func ParseCharacters(root string, group *sync.WaitGroup) {
	characters := make([]Info, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// Find BP_PLayer_* files and create a mapping
		if strings.HasPrefix(info.Name(), "BP_Player_") {
			characterMapping := createBPPlayerMapping(path)
			id := characterId(path)
			// This is only loosely coupled with the ability parsing... But it's the best our current structure allows for.
			specialParser := specials.Parsers[id]

			// Generate character information and append to list
			characterInfo := Info{
				id,
				reference.GetName(characterMapping),
				reference.GetDescription(characterMapping),
				characterMapping.getRole(),
				getStanceNameMapAsJson(specialParser),
			}

			if util.IsHeroWhitelisted(characterInfo.Id) {
				characters = append(characters, characterInfo)
				//Copy character images to output folder
				reference.CopyImageFile(characterMapping.CharacterPreviewImage, id+"-preview", group, "characters", "preview")
				reference.CopyImageFile(characterMapping.CharacterDefaultSkinImage, id+"-default", group, "characters", "default-skin")
				reference.CopyImageFile(characterMapping.CharacterPortrait, id+"-portrait", group, "characters", "portraits")
			}
		}
		return nil
	})
	//Write file containing all the characters
	util.Check(err)
	err = util.WriteInfo("characters.json", characters)
	util.Check(err, "Unable to write characters", characters)
}

// getStanceNameMapAsJson Retrieves the specials.Stance->StanceName map from the specific parser, if any exists.
// The result is parsed into a JSON string.
// This leads to an entry like { "StanceMapJson" : { "1":"FighterStance", "2":"DefenderStance" } }
// This information is solely included to convert the numerical stance values into names on the website.
func getStanceNameMapAsJson(parser *specials.StanceParser) *string {
	if parser == nil {
		return nil
	}

	jsonString, err := json.Marshal(parser.StanceResolver())
	util.Check(err, parser.StanceResolver())
	stringified := string(jsonString)
	return &stringified
}

// createBPPlayerMapping creates a BPPlayerMapping from a BP_Player_* file
func createBPPlayerMapping(path string) BPPlayerMapping {
	content, err := os.ReadFile(path)
	util.Check(err, path)
	characterRawJson := gjson.Get(string(content), "#(Type%\"BP_Player_*\")#|0.Properties").String()
	var characterMapping BPPlayerMapping
	err = json.Unmarshal([]byte(characterRawJson), &characterMapping)
	util.Check(err, path)
	return characterMapping
}

func (m BPPlayerMapping) getRole() string {
	if m.CharacterKitRole.TableId != "" {
		return reference.GetReferenceValue(m.CharacterKitRole)
	}

	if m.CharacterKitRole.SourceString != "" {
		return m.CharacterKitRole.SourceString
	}

	return "UnknownRoleProperty"
}

// characterId returns the id of a character based on the filepath
func characterId(path string) string {
	folders := strings.Split(path, string(os.PathSeparator))
	return slug.Make(reference.AddSpace(folders[4]))
}
