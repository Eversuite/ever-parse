package talent

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
	MetaPowerIcon        reference.ImageReference
	MetaPowerTitle       reference.PropertyReference
	MetaPowerDescription reference.PropertyReference
}

type Info struct {
	Id          string
	Name        string
	Description string
	Source      string
}

func (m Mapping) GetNameProperty() reference.PropertyReference {
	return m.MetaPowerTitle
}

func (m Mapping) GetDescriptionProperty() reference.PropertyReference {
	return m.MetaPowerDescription
}

func ParseTalents(root string) {
	talents := make([]Info, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == "Shards" {
			return filepath.SkipDir
		}
		if strings.HasPrefix(info.Name(), "MPUI_") {
			content, err := os.ReadFile(path)
			util.Check(err, path)
			talentRawJson := gjson.Get(string(content), "#(Type%\"VVMetaPowerUIData*\")#|0.Properties").String()
			var talentMapping Mapping
			err = json.Unmarshal([]byte(talentRawJson), &talentMapping)
			if err != nil {
				println("Failed to parse: " + path)
				return nil
			}
			util.Check(err, path)
			id := slug.Make(reference.TalentId(path))
			talentInfo := Info{
				id,
				reference.GetName(talentMapping),
				reference.GetDescription(talentMapping),
				slug.Make(reference.Source(path)),
			}
			talents = append(talents, talentInfo)
			reference.CopyImageFile(talentMapping.MetaPowerIcon, id, "talent")
		}
		return nil
	})
	util.Check(err)
	err = util.WriteInfo("talent.json", talents)
	util.Check(err, "talent.json", talents)
}
