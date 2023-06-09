package shard

func IsBlacklisted(name string) bool {
	var blacklist = [49]string{
		"armour",
		"attack-speed-star",
		"backwards-shot",
		"chaotic-swirl",
		"charge-shield",
		"clip-discharge",
		"collateral",
		"cool-headed",
		"crater-making",
		"damage-amp",
		"damage-aura",
		"damage-reduction",
		"doublenado",
		"echoing-fissure",
		"enveloping-flames",
		"extended-range",
		"healing-buff-petals",
		"healing-swap",
		"health-sprout-overflow",
		"heat-pressure",
		"ignite-bomb",
		"noxious-bloom",
		"protection",
		"protective-slash",
		"pummel-kinetic-actuators",
		"pummel-pressure-points",
		"rapid-fire",
		"rapid-reload",
		"ray-shot-reactive-shield",
		"recombined-ammo",
		"root-pummel",
		"self-heal",
		"shen-armour-melting",
		"slash-armour-debuff",
		"slide-auto-reload",
		"slide-evasive-manouvers",
		"speed-boost",
		"split-beam",
		"spread-shot",
		"sucker-punch",
		"supernova",
		"teleport-attack-dmg",
		"teleport-shield",
		"thorns",
		"tornado-daze",
		"untamed-gale",
		"vendor-token",
		"vireball-chain-reaction",
		"vireball-velocity-flashover",
	}

	for _, n := range blacklist {
		if n == name {
			return true
		}
	}
	return false
}
