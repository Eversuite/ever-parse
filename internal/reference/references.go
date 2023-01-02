package reference

import (
	"bytes"
	"encoding/json"
	"ever-parse/internal/util"
	"github.com/gosimple/slug"
	"github.com/tidwall/gjson"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

type PropertyReference struct {
	TableId      string
	Key          string
	SourceString string
}

type ObjectReference struct {
	ObjectPath string
}

// CurveTableReference type alias to work with reference entries.
// The CurveTableReference contains an arbitrary amount of properties with their corresponding CurveTable references.
type CurveTableReference map[string]CurveTableReferenceEntry

// CurveTableReferenceEntry contains a reference to the CurveTable and its corresponding row associated with the property
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

// DataMapping is an interface for mapping types which grant access to a "name" property and "description" property.
// The properties must be PropertyReference s
type DataMapping interface {
	GetNameProperty() PropertyReference
	GetDescriptionProperty() PropertyReference
}

type CurvePropertiesMapping interface {
	GetCurveProperty() CurveTableReference
}

const noneName string = "None"

const ProjectVImagePath = "icons"

var whitespaceRegex = regexp.MustCompile("\\s")
var jsonRegex = regexp.MustCompile("\\..*")

// GetName resolves the actual name for the "name"-property of the DataMapping.
// This is either done by resolving the corresponding table entry or the property's own `SourceString`-entry.
func GetName(m DataMapping) string {
	if m.GetNameProperty().TableId != "" {
		return GetReferenceValue(m.GetNameProperty())
	}

	if m.GetNameProperty().SourceString != "" {
		return m.GetNameProperty().SourceString
	}

	return "UnknownNameProperty"
}

// GetDescription resolves the actual description for the "Description"-property of the DataMapping.
func GetDescription(m DataMapping) string {
	if m.GetNameProperty().TableId != "" {
		return GetReferenceValue(m.GetDescriptionProperty())
	}

	if m.GetNameProperty().SourceString != "" {
		return m.GetDescriptionProperty().SourceString
	}

	return "UnknownDescriptionProperty"
}

func GetCurveProperties(m CurvePropertiesMapping) string {
	abilityProps := ""
	if m.GetCurveProperty() != nil {
		jsonBytes, err := json.Marshal(m.GetCurveProperty().GetValues())
		util.Check(err, m)
		abilityProps = string(jsonBytes)
	}
	return abilityProps
}

func (c CurveTableReference) GetValues() map[string][]CurvePoint {
	result := make(map[string][]CurvePoint, len(c))
	for key, entry := range c {
		points := entry.getValue()
		result[key] = points
	}
	return result
}

func (ce CurveTableReferenceEntry) getValue() (curvePoints []CurvePoint) {
	if ce.RowName == noneName || len(ce.CurveTable.ObjectPath) == 0 {
		return
	}
	correctRoot := FixRoot(ce.CurveTable.ObjectPath)
	cleanedPath := jsonRegex.ReplaceAllString(correctRoot, ".json")
	content, err := os.ReadFile(cleanedPath)
	util.Check(err, ce, correctRoot, cleanedPath)
	curvePointsJson := gjson.Get(string(content), "#.Rows."+ce.RowName+".Keys|0").String()
	curvePointsJson = whitespaceRegex.ReplaceAllString(curvePointsJson, "")
	if len(curvePointsJson) == 0 {
		return
	}
	err = json.Unmarshal([]byte(curvePointsJson), &curvePoints)
	util.Check(err, ce.CurveTable.ObjectPath, ce.RowName, curvePointsJson)
	return
}

func GetReferenceValue(propertyReference PropertyReference) string {
	correctRoot := FixRoot(propertyReference.TableId)
	cleanedPath := jsonRegex.ReplaceAllString(correctRoot, ".json")
	content, err := os.ReadFile(cleanedPath)
	util.Check(err, propertyReference, correctRoot, cleanedPath)
	return gjson.Get(string(content), "#.StringTable.KeysToMetaData."+propertyReference.Key+"|0").String()
}

func CopyImageFile(abilityIcon ObjectReference, id string, paths ...string) {
	if abilityIcon.ObjectPath == "" {
		return
	}
	correctRoot := FixRoot(abilityIcon.ObjectPath)
	cleanedPath := jsonRegex.ReplaceAllString(correctRoot, ".png")
	content, err := os.ReadFile(cleanedPath)
	util.Check(err, cleanedPath)

	// Build the image path.
	paths = append([]string{".", ProjectVImagePath}, paths...)
	dir, err := util.CreateDir(paths...)
	util.Check(err)

	file := filepath.Join(dir, id+".png")
	err = os.WriteFile(file, content, 0644)
	util.Check(err, content)
}

func GenerateId(path string, delimiter string) string {
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
	changedToDash := strings.ReplaceAll(removedFileEnding, "_", " ")
	return slug.Make(AddSpace(changedToDash))
}

func Source(path string) string {
	folders := strings.Split(path, string(os.PathSeparator))
	return AddSpace(folders[3])
}

// AddSpace Adds a space before any uppercase character
// Example: MiniTank would be Mini Tank
func AddSpace(s string) string {
	buf := &bytes.Buffer{}
	for i, character := range s {
		if unicode.IsUpper(character) && i > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteRune(character)
	}
	return buf.String()
}

func FixRoot(path string) string {
	return strings.ReplaceAll(path, "/Game/", "Game/")
}
