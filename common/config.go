package common

import "github.com/disgoorg/disgo/discord"

var (
	MaxFreeSaves    = 5
	MaxPremiumSaves = 25
	LinksActionRow  = discord.ActionRowComponent{
		discord.NewLinkButton("Documentation", "https://trade.meta-bee.com/"),
		discord.NewLinkButton("Support server", "https://discord.gg/meta-bee-995988457136603147"),
	}
)
