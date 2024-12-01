package syncplugin

import (
	"alaninnovates.com/trade-builder/common"
	"alaninnovates.com/trade-builder/common/loaders"
	"alaninnovates.com/trade-builder/database"
	"context"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

/*
/subscribe <type:lf/offering> <sticker_name>

premium only: recieve a dm whenver someone posts a trade with this item
*/
func SubscribeCommand(b *common.Bot) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name:        "subscribe",
			Description: "Subscribe to sticker trades",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "type",
					Description: "Type of trade",
					Required:    true,
					Choices: []discord.ApplicationCommandOptionChoiceString{
						{
							Name:  "Looking For",
							Value: "lf",
						},
						{
							Name:  "Offering",
							Value: "offering",
						},
					},
				},
				discord.ApplicationCommandOptionString{
					Name:         "sticker",
					Description:  "Sticker to subscribe to",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"subscribe": func(event *events.ApplicationCommandInteractionCreate) error {
				data := event.SlashCommandInteractionData()
				sticker := data.String("sticker")
				sticker = strings.ToLower(sticker)
				tradeType := data.String("type")
				_, err := b.Db.Collection("subscriptions").InsertOne(context.TODO(), database.StickerSubscription{
					UserId:    event.Member().User.ID.String(),
					Sticker:   sticker,
					TradeType: tradeType,
				})
				if err != nil {
					return err
				}
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title:       "Success!",
							Description: "You have subscribed to " + tradeType + " trades for " + sticker,
							Color:       common.ColorSuccess,
						},
					},
				})
			},
		},
		AutocompleteHandlers: map[string]handler.AutocompleteHandler{
			"sticker": makeAutocompleteHandler(loaders.GetAllStickers(), "sticker"),
		},
	}
}

func ServersyncCommand(b *common.Bot) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name:        "serversync",
			Description: "Server sync",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "setup",
					Description: "Set up server sync",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionChannel{
							Name:        "channel",
							Description: "The channel to set up server sync in",
							Required:    true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "remove",
					Description: "Remove server sync",
				},
			},
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"setup": func(event *events.ApplicationCommandInteractionCreate) error {
				data := event.SlashCommandInteractionData()
				channel := data.Channel("channel")
				webhook, err := b.Client.Rest().CreateWebhook(channel.ID, discord.WebhookCreate{
					Name:   b.BotInfo.Username,
					Avatar: b.BotInfo.AvatarIcon,
				})
				if err != nil {
					return err
				}
				_, err = b.Db.Collection("serversync").InsertOne(context.TODO(), database.ServerSync{
					ServerId:     event.GuildID().String(),
					ChannelId:    channel.ID.String(),
					WebhookId:    webhook.ID().String(),
					WebhookToken: webhook.Token,
				})
				if err != nil {
					return err
				}
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title:       "Success!",
							Description: "Server sync has been set up in <#" + channel.ID.String() + ">",
							Color:       common.ColorSuccess,
						},
					},
				})
			},
			"remove": func(event *events.ApplicationCommandInteractionCreate) error {
				_, err := b.Db.Collection("serversync").DeleteOne(context.TODO(), bson.D{{"server_id", event.GuildID().String()}})
				if err != nil {
					return err
				}
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title:       "Success!",
							Description: "Server sync has been removed.",
							Color:       common.ColorSuccess,
						},
					},
				})
			},
		},
	}
}

func makeAutocompleteHandler(options []string, fieldName string) func(*events.AutocompleteInteractionCreate) error {
	return func(event *events.AutocompleteInteractionCreate) error {
		//fmt.Printf("evt: %d now: %d", event.ID().Time().UnixMilli(), time.Now().UnixMilli())
		input := event.Data.String(fieldName)
		input = strings.ToLower(input)
		matches := make([]discord.AutocompleteChoice, 0)
		i := 0
		for _, opt := range options {
			if i >= 25 {
				break
			}
			if strings.Contains(strings.ToLower(opt), input) {
				matches = append(matches, discord.AutocompleteChoiceString{
					Name:  opt,
					Value: opt,
				})
				i++
			}
		}
		return event.AutocompleteResult(matches)
	}
}

func Initialize(h *handler.Handler, b *common.Bot) {
	h.AddCommands(SubscribeCommand(b), ServersyncCommand(b), PostCommand(b))
}
