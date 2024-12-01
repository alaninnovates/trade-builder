package marketplugin

import (
	"alaninnovates.com/trade-builder/common"
	"alaninnovates.com/trade-builder/common/loaders"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"strconv"
)

var categories = []discord.ApplicationCommandOptionChoiceString{
	{
		Name:  "All",
		Value: "all stickers",
	},
	{
		Name:  "Cubs",
		Value: "cubs",
	},
	{
		Name:  "Hives",
		Value: "hives",
	},
	{
		Name:  "Bees",
		Value: "bees",
	},
	{
		Name:  "Bears",
		Value: "bears",
	},
	{
		Name:  "Mobs",
		Value: "mobs",
	},
	{
		Name:  "Critters",
		Value: "critters",
	},
	{
		Name:  "Nectars",
		Value: "nectars",
	},
	{
		Name:  "Flowers",
		Value: "flowers",
	},
	{
		Name:  "Puffshrooms",
		Value: "puffshrooms",
	},
	{
		Name:  "Leaves",
		Value: "leaves",
	},
	{
		Name:  "Tools",
		Value: "tools",
	},
	{
		Name:  "Star Signs",
		Value: "star signs",
	},
	{
		Name:  "Beesmas Lights",
		Value: "beesmas lights",
	},
	{
		Name:  "Field Stamps",
		Value: "field stamps",
	},
}

var stickerCategories = map[string][]string{
	"all stickers":   loaders.GetAllStickers(),
	"cubs":           {"Bee Cub Skin", "Brown Cub Skin", "Doodle Cub Skin", "Gingerbread Cub Skin", "Noob Cub Skin", "Peppermint Robo Cub Skin", "Robo Cub Skin", "Snow Cub Skin", "Star Cub Skin", "Stick Cub Skin"},
	"hives":          {"Basic Black Hive Skin", "Basic Blue Hive Skin", "Basic Green Hive Skin", "Basic Pink Hive Skin", "Basic Red Hive Skin", "Basic White Hive Skin", "Wavy Cyan Hive Skin", "Wavy Doodle Hive Skin", "Wavy Festive Hive Skin", "Wavy Purple Hive Skin", "Wavy Yellow Hive Skin"},
	"bees":           {"4-Pronged Vector Bee", "Blob Bumble Bee", "Bomber Bee Bear", "Diamond Diamond Bee", "Drooping Stubborn Bee", "Flying Brave Bee", "Flying Festive Bee", "Flying Ninja Bee", "Flying Photon Bee", "Flying Rad Bee", "Round Basic Bee", "Round Rascal Bee", "Wobbly Looker Bee"},
	"bears":          {"Chef Hat Polar Bear", "Flying Bee Bear", "Glowering Gummy Bear", "Honey Bee Bear", "Panicked Science Bear", "Party Robo Bear", "Royal Bear", "Shy Brown Bear", "Sideways Spirit Bear", "Sitting Green Shirt Bear", "Sitting Mother Bear", "Squashed Head Bear", "Stretched Head Bear", "Uplooking Bear"},
	"mobs":           {"Forward Facing Aphid", "Forward Facing Spider", "Left Facing Ant", "Menacing Mantis", "Right Facing Stump Snail", "Standing Bean Bug", "Standing Caterpillar", "Walking Stick Nymph"},
	"critters":       {"Blue Triangle Critter", "Critter In A Stocking", "Flying Magenta Critter", "Grey Shape Companion", "Orange Leg Critter", "Purple Pointed Critter", "Round Green Critter"},
	"nectars":        {"Comforting Nectar Icon", "Invigorating Nectar Icon", "Motivating Nectar Icon", "Refreshing Nectar Icon", "Satisfying Nectar Icon"},
	"flowers":        {"Blue Flower Field Stamp", "Clover Field Stamp", "Dandelion Field Stamp", "Purple 4-Point Flower", "Rose Field Stamp", "Small Dandelion", "Small Pink Tulip", "Small Tickseed", "Small White Daisy", "Sunflower Field Stamp"},
	"puffshrooms":    {"Spore Covered Puffshroom", "Black Truffle Mushroom", "Chanterelle Mushroom", "Fly Agaric Mushroom", "Morel Mushroom", "Oiler Mushroom", "Porcini Mushroom", "Prismatic Mushroom", "White Button Mushroom"},
	"leaves":         {"Blowing Leaf", "Cordate Leaf", "Cunate Leaf", "Elliptic Leaf", "Hastate Leaf", "Lanceolate Leaf", "Lyrate Leaf", "Oblique Leaf", "Reniform Leaf", "Rhomboid Leaf", "Spatulate Leaf"},
	"tools":          {"Bubble Wand", "Clippers", "Dark Scythe", "Golden Rake", "Honey Dipper", "Magnet", "Petal Wand", "Porcelain Dipper", "Rake", "Scissors", "Scooper", "Scythe", "Spark Staff", "Super-Scooper", "Tide Popper", "Vacuum"},
	"star signs":     {"Aquarius Star Sign", "Aries Star Sign", "Cancer Star Sign", "Capricorn Star Sign", "Gemini Star Sign", "Leo Star Sign", "Libra Star Sign", "Pisces Star Sign", "Sagittarius Star Sign", "Scorpio Star Sign", "Taurus Star Sign", "Virgo Star Sign"},
	"beesmas lights": {"Blue Beesmas Light", "Green Beesmas Light", "Red Beesmas Light", "Yellow Beesmas Light"},
	"field stamps":   {"Ant Field Stamp", "Bamboo Field Stamp", "Blue Flower Field Stamp", "Cactus Field Stamp", "Clover Field Stamp", "Coconut Field Stamp", "Dandelion Field Stamp", "Hub Field Stamp", "Mountain Top Field Stamp", "Mushroom Field Stamp", "Pepper Patch Stamp", "Pine Tree Forest Stamp", "Pineapple Patch Stamp", "Pumpkin Patch Stamp", "Rose Field Stamp", "Spider Field Stamp", "Strawberry Field Stamp", "Stump Field Stamp", "Sunflower Field Stamp"},
	"vouchers":       {"Bear Bee Voucher", "Cub Buddy Voucher", "x2 Bee Gather Voucher", "x2 Convert Speed"},
}

func TopCommand(b *common.Bot) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name:        "top",
			Description: "Trade market",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "demand",
					Description: "Top demand for looking for",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "duration",
							Description: "day/week",
							Required:    true,
							Choices: []discord.ApplicationCommandOptionChoiceString{
								{
									Name:  "Day",
									Value: "daily",
								},
								{
									Name:  "Week",
									Value: "weekly",
								},
							},
						},
						discord.ApplicationCommandOptionString{
							Name:        "category",
							Description: "Category of the sticker",
							Required:    true,
							Choices:     categories,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "offer",
					Description: "Top demand for offers",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "duration",
							Description: "day/week",
							Required:    true,
							Choices: []discord.ApplicationCommandOptionChoiceString{
								{
									Name:  "Day",
									Value: "daily",
								},
								{
									Name:  "Week",
									Value: "weekly",
								},
							},
						},
						discord.ApplicationCommandOptionString{
							Name:        "category",
							Description: "Category of the sticker",
							Required:    true,
							Choices:     categories,
						},
					},
				},
			},
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"demand": func(event *events.ApplicationCommandInteractionCreate) error {
				data := event.SlashCommandInteractionData()
				duration := data.String("duration")
				category := data.String("category")

				redisKeys := make([]string, 0)
				for _, sticker := range stickerCategories[category] {
					redisKeys = append(redisKeys, duration+":lookingfor:"+sticker)
				}

				redisValues := b.Redis.Client().MGet(b.Redis.Context(), redisKeys...).Val()

				message := "```"

				idx := 0
				for i, sticker := range stickerCategories[category] {
					if redisValues[i] == nil {
						continue
					}
					message += strconv.Itoa(idx+1) + ". " + sticker + ": " + redisValues[i].(string) + "\n"
					idx++
				}

				message += "```"

				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title:       "Top demand for " + category + " in the last " + duration,
							Description: message,
							Color:       common.ColorPrimary,
						},
					},
				})
			},
			"offer": func(event *events.ApplicationCommandInteractionCreate) error {
				data := event.SlashCommandInteractionData()
				duration := data.String("duration")
				category := data.String("category")

				redisKeys := make([]string, 0)
				for _, sticker := range stickerCategories[category] {
					redisKeys = append(redisKeys, duration+":offering:"+sticker)
				}

				redisValues := b.Redis.Client().MGet(b.Redis.Context(), redisKeys...).Val()

				message := "```"

				idx := 0
				for i, sticker := range stickerCategories[category] {
					if redisValues[i] == nil {
						continue
					}
					message += strconv.Itoa(idx+1) + ". " + sticker + ": " + redisValues[i].(string) + "\n"
					idx++
				}

				message += "```"

				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title:       "Top offers for " + category + " in the last " + duration,
							Description: message,
							Color:       common.ColorPrimary,
						},
					},
				})
			},
		},
	}
}

func Initialize(h *handler.Handler, b *common.Bot) {
	h.AddCommands(TopCommand(b))
}
