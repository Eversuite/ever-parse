package cleanup

import (
	"encoding/json"
	"ever-parse/internal/character"
	"ever-parse/internal/util"
	"os"
)

func CleanupCharacters() {
	// Read the JSON file into a byte slice
	data, err := os.ReadFile("parsedData/characters.json")
	if err != nil {
		panic(err)
	}

	// Parse the JSON data into a slice of Person structs
	var characters []character.Info
	err = json.Unmarshal(data, &characters)
	if err != nil {
		panic(err)
	}

	for i, character := range characters {
		if character.Id == "spell-healer" {
			characters[i].Id = "ho-t-healer"
		}
		if character.Id == "spell-d-p-s" {
			characters[i].Id = "spellcasting-d-p-s"
		}

	}

	// Write the JSON data to a file
	err = util.WriteInfo("characters.json", characters)
	util.Check(err, "characters.json", characters)
}
