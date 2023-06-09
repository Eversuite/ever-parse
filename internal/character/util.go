package character

func IsBlacklisted(name string) bool {
	var blacklist = [2]string{
		"chef-support",
		"gun-d-p-s",
	}

	for _, n := range blacklist {
		if n == name {
			return true
		}
	}
	return false
}

func FixCharacterData(characters []Info) []Info {
	for i, character := range characters {
		if character.Id == "spell-healer" {
			characters[i].Id = "ho-t-healer"
		}
		if character.Id == "spell-d-p-s" {
			characters[i].Id = "spellcasting-d-p-s"
		}
	}
	return characters
}
