package bssdata

import "github.com/disgoorg/disgo/discord"

var SoftWaxEmoji = "<:SoftWax:1340564900581347360:>"
var HardWaxEmoji = "<:HardWax:1340564899037839530:>"
var CausticWaxEmoji = ":<CausticWax:1340564895522881606:>"
var SwirledWaxEmoji = ":<SwirledWax:1340564904117141554:>"

var WaxSelectMenuOptions = []discord.StringSelectMenuOption{
	{
		Label: "Soft Wax",
		Value: "Soft Wax",
		Emoji: &discord.ComponentEmoji{
			ID:   1340564900581347360,
			Name: "SoftWax",
		},
	},
	{
		Label: "Hard Wax",
		Value: "Hard Wax",
		Emoji: &discord.ComponentEmoji{
			ID:   1340564899037839530,
			Name: "HardWax",
		},
	},
	{
		Label: "Caustic Wax",
		Value: "Caustic Wax",
		Emoji: &discord.ComponentEmoji{
			ID:   1340564895522881606,
			Name: "CausticWax",
		},
	},
	{
		Label: "Swirled Wax",
		Value: "Swirled Wax",
		Emoji: &discord.ComponentEmoji{
			ID:   1340564904117141554,
			Name: "SwirledWax",
		},
	},
}
