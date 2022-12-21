package reference

import (
	"bytes"
	"github.com/tidwall/gjson"
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

func GetReferenceValue(propertyReference PropertyReference) string {
	correctRoot := fixRoot(propertyReference.TableId)
	regex := regexp.MustCompile("\\..*")
	cleanedPath := regex.ReplaceAllString(correctRoot, ".json")
	content, err := os.ReadFile(cleanedPath)
	check(err, "")
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
	check(err, "")
	err = os.WriteFile("./icons/"+id+".png", content, 0644)
	check(err, "")
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

func check(e error, path string) {
	if e != nil {
		println(path)
		panic(e)
	}
}
