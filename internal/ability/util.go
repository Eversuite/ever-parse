package ability

import (
	"ever-parse/internal/reference"
)

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

func GetAbilityName(m BPUIAbilityMapping, superMapping *BPUIAbilityMapping) string {
	this := reference.GetName(m)
	if this != reference.UnknownNameProperty {
		return this
	}

	if superMapping != nil {
		return reference.GetName(*superMapping)
	}

	return reference.UnknownNameProperty
}

func GetAbilityCurveValues(m BPUIAbilityMapping, superMapping *BPUIAbilityMapping) string {
	this := reference.GetCurveProperties(m)
	if this != "" {
		return this
	}

	if superMapping != nil {

		return reference.GetCurveProperties(*superMapping)
	}
	return this
}

func GetAbilityDescription(m BPUIAbilityMapping, superMapping *BPUIAbilityMapping) string {
	this := reference.GetDescription(m)
	if this != reference.UnknownDescriptionProperty {
		return this
	}

	if superMapping != nil {
		return reference.GetDescription(*superMapping)
	}

	return reference.UnknownDescriptionProperty
}
