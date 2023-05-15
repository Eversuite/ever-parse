package reference

import (
	"bytes"
	"ever-parse/internal/util"
	"github.com/gosimple/slug"
	"github.com/tidwall/gjson"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
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

// DataMapping is an interface for mapping types which grant access to a "name" property and "description" property.
// The properties must be PropertyReference s
type DataMapping interface {
	GetNameProperty() PropertyReference
	GetDescriptionProperty() PropertyReference
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

func GetReferenceValue(propertyReference PropertyReference) string {
	correctRoot := FixRoot(propertyReference.TableId)
	cleanedPath := jsonRegex.ReplaceAllString(correctRoot, ".json")
	content, err := os.ReadFile(cleanedPath)
	util.Check(err, propertyReference, correctRoot, cleanedPath)
	return gjson.Get(string(content), "#.StringTable.KeysToMetaData."+propertyReference.Key+"|0").String()
}

func CopyImageFile(abilityIcon ObjectReference, id string, group *sync.WaitGroup, paths ...string) {
	group.Add(1)

	go func() {
		defer group.Done()

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

		croppedImage, didCrop := cropImage(cleanedPath)

		if err == nil && didCrop {
			croppedFileName := filepath.Join(dir, id+"-cropped.png")
			croppedFile, err := os.OpenFile(croppedFileName, os.O_CREATE|os.O_RDWR, 0644)
			util.Check(err, "unable to create/write cropped file: "+id+"-cropped.png")
			err = png.Encode(croppedFile, *croppedImage)
			util.Check(err, "Could not safe file as PNG", croppedFile)
			return
		}

		file := filepath.Join(dir, id+".png")
		err = os.WriteFile(file, content, 0644)
		util.Check(err, content)
	}()
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
	return slug.Make(AddSpace(folders[3]))
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