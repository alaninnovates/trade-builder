package miscplugin

import (
	"alaninnovates.com/trade-builder/common"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/json"
	"strconv"
)

var helpMenus = map[string]discord.Embed{
	"home": {
		Title:       "Trade Builder Help",
		Description: "Visit the documentation below for a list of all commands. Join the support server if you have any more questions.",
		Footer: &discord.EmbedFooter{
			Text: "Made by alaninnovates#0123",
		},
		Color: 0xffffff,
	},
	"trade": {
		Title: ":hammer: Trade building",
		Description: `
			▸ </trade create:1269548240538435738>
			▸ </trade lookingfor:1269548240538435738> <sticker_name> <quantity>
			▸ </trade offering:1269548240538435738> <sticker_name> <quantity>
			▸ </trade remove:1269548240538435738> <type:lf/offering> <sticker_name>
			▸ </trade view:1269548240538435738>
			▸ </trade info:1269548240538435738>
			▸ </trade save:1269548240538435738> <save_name>
			▸ </trade saves list:1269548240538435738>
			▸ </trade saves load:1269548240538435738> <id>
			▸ </trade saves delete:1269548240538435738> <id>

			[] = Optional | <> = Required`,
		Color: 0xfcba03,
	},
	"website": {
		Title: ":desktop: Website",
		Description: `
			▸ </post:1295106825955577968> <trade_id> [expire_time] [server_sync]

			[] = Optional | <> = Required`,
		Color: 0x03fc73,
	},
	"market": {
		Title: ":chart_with_upwards_trend: Market Demand",
		Description: `
			▸ </top demand:1271998287964143739> <duration> [category]
			▸ </top offer:1271998287964143739> <duration> [category]

			[] = Optional | <> = Required`,
		Color: 0x03b1fc,
	},
	// todo: diff color
	"sync": {
		Title: ":arrows_counterclockwise: Sync",
		Description: `
			▸ </serversync setup:1271998287964143741> <channel>
			▸ </serversync remove:1271998287964143741>

			[] = Optional | <> = Required`,
		Color: 0x03fc73,
	},
}

func HelpCommand(b *common.Bot) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name:        "help",
			Description: "Get help with the bot",
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"": func(event *events.ApplicationCommandInteractionCreate) error {
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						helpMenus["home"],
					},
					Components: []discord.ContainerComponent{
						common.LinksActionRow,
						discord.ActionRowComponent{
							discord.NewStringSelectMenu(
								"handler:help",
								"Select a category",
								discord.StringSelectMenuOption{
									Label: "Home",
									Value: "home",
									Emoji: &discord.ComponentEmoji{Name: "🏠"},
								},
								discord.StringSelectMenuOption{
									Label: "Trade Building",
									Value: "trade",
									Emoji: &discord.ComponentEmoji{Name: "🔨"},
								},
								discord.StringSelectMenuOption{
									Label: "Website",
									Value: "website",
									Emoji: &discord.ComponentEmoji{Name: "🖥️"},
								},
								discord.StringSelectMenuOption{
									Label: "Market Demand",
									Value: "market",
									Emoji: &discord.ComponentEmoji{Name: "📈"},
								},
								discord.StringSelectMenuOption{
									Label: "Sync",
									Value: "sync",
									Emoji: &discord.ComponentEmoji{Name: "🔄"},
								},
							),
						},
					},
				})
			},
		},
	}
}

func HelpComponent(b *common.Bot) handler.Component {
	return handler.Component{
		Name: "help",
		Handler: func(event *events.ComponentInteractionCreate) error {
			sectionName := event.StringSelectMenuInteractionData().Values[0]
			return event.UpdateMessage(discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					helpMenus[sectionName],
				},
			})
		},
	}
}

func StatsCommand(b *common.Bot) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name:        "stats",
			Description: "Get statistics about the bot",
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"": func(event *events.ApplicationCommandInteractionCreate) error {
				members := 0
				b.Client.Caches().GuildsForEach(func(e discord.Guild) {
					members += e.MemberCount
				})
				guildId, _ := strconv.Atoi(event.GuildID().String())
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Fields: []discord.EmbedField{
								{
									Name:   "Guilds on this shard (not total)",
									Value:  strconv.Itoa(b.Client.Caches().GuildsLen()),
									Inline: json.Ptr(true),
								},
								{
									Name:   "Members",
									Value:  strconv.Itoa(members),
									Inline: json.Ptr(true),
								},
								{
									Name:  "Shard ID",
									Value: strconv.Itoa(guildId >> 22 % len(b.Client.ShardManager().Shards())),
								},
							},
							Footer: &discord.EmbedFooter{
								Text: "Made by alaninnovates#0123",
							},
						},
					},
					Components: []discord.ContainerComponent{
						common.LinksActionRow,
					},
				})
			},
		},
	}
}

func Initialize(h *handler.Handler, b *common.Bot) {
	h.AddCommands(HelpCommand(b), StatsCommand(b))
	h.AddComponents(HelpComponent(b))
}
