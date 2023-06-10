package ability

var slotMapping = map[string]string{
	"twirl":         "Q",
	"take-flight":   "W",
	"graceful-dash": "E",
	"finale":        "R",
}

func FixAbilityData(abilities []Info) []Info {
	for i, ability := range abilities {
		if slotMapping[ability.Id] != "" {
			abilities[i].Slot = slotMapping[ability.Id]
		}
	}
	return abilities
}
