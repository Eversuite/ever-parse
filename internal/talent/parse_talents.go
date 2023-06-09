package talent

import (
	"encoding/json"
	"ever-parse/internal/character"
	"ever-parse/internal/reference"
	"ever-parse/internal/util"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gosimple/slug"
	"github.com/tidwall/gjson"
)

// MPDMapping represents the relevant "Properties" inside VVMetaPowerDefinition types.
type MPDMapping struct {
	MetaPowerCategory  reference.ObjectReference
	MetaPowerUIData    reference.ObjectReference
	MetaPowerTierIndex int
}

// MPUIMapping represents the relevant "Properties" inside VVMetaPowerUIData types.
type MPUIMapping struct {
	MetaPowerTitle       reference.PropertyReference
	MetaPowerDescription reference.PropertyReference
	MetaPowerIcon        reference.ObjectReference
}

func (m MPUIMapping) GetNameProperty() reference.PropertyReference {
	return m.MetaPowerTitle
}

func (m MPUIMapping) GetDescriptionProperty() reference.PropertyReference {
	return m.MetaPowerDescription
}

type Info struct {
	Id          string
	Name        string
	Description string
	Hero        string
	Category    string
	Tier        int
}

// ParseTalents Parses hero talents and writes to the talents.json file
func ParseTalents(root string, group *sync.WaitGroup) {
	talents := make([]Info, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		//Accept all MPD_* files and create mappings
		if strings.HasPrefix(info.Name(), "MPD_") {
			mpdMapping := createMdpMapping(path)
			mpuiFilePath := createMpuiFilePath(mpdMapping)
			mpuiMapping := createMpuiMapping(mpuiFilePath)
			hero := reference.Source(path)

			//Generate talent information and append to list
			talentInfo := Info{
				talentId(mpuiFilePath, hero),
				reference.GetName(mpuiMapping),
				util.ToValidHtml(reference.GetDescription(mpuiMapping)),
				hero,
				GenerateTalentCategoryId(hero, mpdMapping.MetaPowerCategory.ObjectPath),
				mpdMapping.MetaPowerTierIndex,
			}
			if character.IsBlacklisted(talentInfo.Hero) {
				talents = append(talents, talentInfo)
				//Copy the talent icon to the output folder
				reference.CopyImageFile(mpuiMapping.MetaPowerIcon, talentInfo.Id, group, "talent")
			}
		}
		return nil
	})
	//Write file containing all the talents
	util.Check(err)
	err = util.WriteInfo("talents.json", talents)
	util.Check(err, "talents.json", talents)
}

// createMdpMapping Parses the "Properties" field inside a VVMetaPowerDefinition type and converts it to an MPDMapping.
func createMdpMapping(path string) MPDMapping {
	mpdContent, err := os.ReadFile(path)
	util.Check(err, path)
	mpdJson := gjson.Get(string(mpdContent), "#(Type%\"VVMetaPowerDefinition*\")#|0.Properties").String()
	var mpdMapping MPDMapping
	err = json.Unmarshal([]byte(mpdJson), &mpdMapping)
	util.Check(err, path)
	return mpdMapping
}

// createMpuiFilePath Get the MPUI file path based on the MPDMapping.
func createMpuiFilePath(mpdMapping MPDMapping) string {
	mpuiFilePath :=
		reference.FixRoot(
			strings.ReplaceAll(mpdMapping.MetaPowerUIData.ObjectPath, ".0", ".json"))
	return mpuiFilePath
}

// createMpuiMapping Parses the "Properties" field from a VVMetaPowerUIData type and converts it to a MPUIMapping.
func createMpuiMapping(mpuiFilePath string) MPUIMapping {
	mpuiContent, err := os.ReadFile(mpuiFilePath)
	util.Check(err, mpuiFilePath)
	mpuiJson := gjson.Get(string(mpuiContent), "#(Type%\"VVMetaPowerUIData*\")#|0.Properties").String()
	var mpuiMapping MPUIMapping
	err = json.Unmarshal([]byte(mpuiJson), &mpuiMapping)
	util.Check(err, mpuiFilePath)
	return mpuiMapping
}

// GenerateTalentCategoryId Finds the talent category id by parsing the string character for character in reverse.
// The reason for going: it is easier to find the right substring due to multiple _
func GenerateTalentCategoryId(source string, objectPath string) (subString string) {
	reversedString := reverseString(objectPath)
	startOfSubstring := strings.Index(reversedString, ".")
	if startOfSubstring == -1 {
		return reverseString(subString)
	}
	tempString := reversedString[startOfSubstring+len("."):]
	endOfSubstring := strings.Index(tempString, "_")
	if endOfSubstring == -1 {
		return reverseString(subString)
	}
	subString = tempString[:endOfSubstring]
	talentCategory := reverseString(subString)
	return slug.Make(source + "-" + reference.AddSpace(talentCategory))
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func talentId(path string, source string) string {
	delimiter := "MPUI_"
	return source + "-" + reference.GenerateId(path, delimiter)
}
