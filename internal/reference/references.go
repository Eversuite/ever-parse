package reference

import (
	"bytes"
	"encoding/json"
	"ever-parse/internal/util"
	"github.com/tidwall/gjson"
	"io/fs"
	"os"
	"regexp"
	"strings"
	"unicode"
)

type PropertyReference struct {
	TableId      string
	Key          string
	SourceString string
}

type ImageReference struct {
	ObjectPath string
}

type CurveTableReference map[string]CurveTableReferenceEntry

type CurveTableReferenceEntry struct {
	CurveTable struct {
		ObjectPath string
	}
	RowName string
}

type CurvePoint struct {
	Time  float64
	Value float64
}

const NONE_NAME string = "None"

func (c CurveTableReference) GetValues() (result map[string][]CurvePoint) {
	result = make(map[string][]CurvePoint, len(c))
	for key, entry := range c {
		points := entry.getValue()
		result[key] = points
	}
	return
}

func (ce CurveTableReferenceEntry) getValue() (curvePoints []CurvePoint) {
	if ce.RowName == NONE_NAME || len(ce.CurveTable.ObjectPath) == 0 {
		return
	}
	correctRoot := fixRoot(ce.CurveTable.ObjectPath)
	regex := regexp.MustCompile("\\..*")
	cleanedPath := regex.ReplaceAllString(correctRoot, ".json")
	content, err := os.ReadFile(cleanedPath)
	util.Check(err, ce, correctRoot, cleanedPath)
	curvePointsJson := gjson.Get(string(content), "#.Rows."+ce.RowName+".Keys|0").String()
	if len(curvePointsJson) == 0 {
		return
	}
	err = json.Unmarshal([]byte(curvePointsJson), &curvePoints)
	util.Check(err, ce.CurveTable.ObjectPath, ce.RowName, curvePointsJson)
	return
}

func GetReferenceValue(propertyReference PropertyReference) string {
	correctRoot := fixRoot(propertyReference.TableId)
	regex := regexp.MustCompile("\\..*")
	cleanedPath := regex.ReplaceAllString(correctRoot, ".json")
	content, err := os.ReadFile(cleanedPath)
	util.Check(err, cleanedPath)
	return gjson.Get(string(content), "#.StringTable.KeysToMetaData."+propertyReference.Key+"|0").String()
}

func CopyImageFile(abilityIcon ImageReference, id string) {
	if abilityIcon.ObjectPath == "" {
		return
	}
	correctRoot := fixRoot(abilityIcon.ObjectPath)
	regex := regexp.MustCompile("\\..*")
	cleanedPath := regex.ReplaceAllString(correctRoot, ".png")
	content, err := os.ReadFile(cleanedPath)
	util.Check(err, cleanedPath)
	err = os.MkdirAll("./icons", fs.ModeDir)
	util.Check(err)
	err = os.WriteFile("./icons/"+id+".png", content, 0644)
	util.Check(err, content)
}

func AbilityId(path string) string {
	delimiter := "BP_UIAbility_"
	pos := strings.LastIndex(path, delimiter)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(delimiter)
	if adjustedPos >= len(path) {
		return ""
	}
	rawAbilityName := path[adjustedPos:]
	removedFileEnding := strings.ReplaceAll(rawAbilityName, ".json", "")
	return addSpace(strings.ReplaceAll(removedFileEnding, "_", "-"))
}

func CharacterId(path string) string {
	folders := strings.Split(path, string(os.PathSeparator))
	return addSpace(folders[4])
}

func AbilitySource(path string) string {
	folders := strings.Split(path, string(os.PathSeparator))
	return addSpace(folders[3])
}

func addSpace(s string) string {
	buf := &bytes.Buffer{}
	for i, character := range s {
		if unicode.IsUpper(character) && i > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteRune(character)
	}
	return buf.String()
}

func fixRoot(path string) string {
	return strings.ReplaceAll(path, "/Game/", "Game/")
}