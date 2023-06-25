package shard

func IsBlacklisted(name string) bool {
	var blacklist = []string{
		"armour",
		"ashen-armour-melting",
		"attack-speed-star",
		"backwards-shot",
		"boss-shard-enrage",
		"boss-shard-flametouch",
		"boss-shard-max-health",
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
		"ghost-pepper",
		"healing-buff-petals",
		"healing-swap",
		"health-sprout-overflow",
		"heat-pressure",
		"ignite-bomb",
		"luum-fissure",
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
		"sorcerers-gloves",
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
