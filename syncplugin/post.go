package syncplugin

import (
	"alaninnovates.com/trade-builder/common"
	"alaninnovates.com/trade-builder/database"
	"alaninnovates.com/trade-builder/tradeplugin/trade"
	"context"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/snowflake/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
)

func PostCommand(b *common.Bot) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name:        "post",
			Description: "Post a saved trade to the website",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "trade_id",
					Description: "The id of your saved trade",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "expire_time",
					Description: "How long before your trade post expires (in HH:MM or mm/dd UTC)",
					Required:    false,
				},
				discord.ApplicationCommandOptionBool{
					Name:        "server_sync",
					Description: "Post to servers subscribed with server sync",
					Required:    false,
				},
			},
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"": func(event *events.ApplicationCommandInteractionCreate) error {
				data := event.SlashCommandInteractionData()
				tradeID := data.String("trade_id")
				expireTime, expExists := data.OptString("expire_time")
				serverSync, ssExists := data.OptBool("server_sync")

				if !ssExists {
					serverSync = true
				}

				var duration time.Duration
				if expExists {
					var err error
					duration, err = common.ParseHHMM(expireTime)
					if err != nil {
						duration, err = common.ParseMMDD(expireTime)
					}
					if err != nil {
						return event.CreateMessage(discord.MessageCreate{
							Content: "Invalid duration.",
						})
					}
					if duration.Minutes() < 10 {
						return event.CreateMessage(discord.MessageCreate{
							Content: "Duration must be at least 10 minutes.",
						})
					}
				} else {
					duration = 3 * 24 * time.Hour
				}

				if tradeID == "" {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You need to provide a non-empty id for the trade.",
					})
				}
				oid, err := primitive.ObjectIDFromHex(tradeID)
				if err != nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "Invalid id.",
					})
				}
				var h bson.D
				err = b.Db.Collection("trades").FindOne(context.Background(), bson.M{
					"user_id": event.User().ID,
					"_id":     oid,
				}).Decode(&h)
				if err != nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have a trade with that id.",
					})
				}

				var t database.Trade
				for _, v := range h {
					if v.Key == "lookingFor" {
						var lookingFor map[string]interface{}
						lookingForBson, _ := bson.Marshal(v.Value)
						_ = bson.Unmarshal(lookingForBson, &lookingFor)
						t.LookingFor = lookingFor
					} else if v.Key == "offering" {
						var offering map[string]interface{}
						offeringBson, _ := bson.Marshal(v.Value)
						_ = bson.Unmarshal(offeringBson, &offering)
						t.Offering = offering
					}
				}

				var avatarUrl string
				avatarUrlQuery := event.User().AvatarURL()

				if avatarUrlQuery == nil {
					discrim, err := strconv.Atoi(event.User().Discriminator)
					if err != nil {
						return err
					}
					avatarUrl = "https://cdn.discordapp.com/embed/avatars/" + strconv.Itoa(discrim%5) + ".png"
				} else {
					avatarUrl = *avatarUrlQuery
				}

				var globalName string
				globalNameQuery := event.User().GlobalName

				if globalNameQuery == nil {
					globalName = event.User().Username
				} else {
					globalName = *globalNameQuery
				}

				post, err := b.Db.Collection("posts").InsertOne(context.TODO(), database.WebsitePost{
					UserId:         event.User().ID.String(),
					UserName:       event.User().Username,
					UserGlobalName: globalName,
					UserAvatar:     avatarUrl,
					ExpireTime:     primitive.NewDateTimeFromTime(time.Now().Add(duration)),
					ServerSync:     serverSync,
					Trade:          t,
					Locked:         false,
				})
				if err != nil {
					return err
				}
				postID := post.InsertedID.(primitive.ObjectID).Hex()

				desc := fmt.Sprintf("Your trade has been posted to the trade builder website. You can view it [here](https://tradebuilder.app/trade/%s)", postID)
				if serverSync {
					desc += "\nYour trade will be posted to servers subscribed with server sync."
				}
				err = event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title:       "Trade posted!",
							Description: desc,
							Color:       common.ColorSuccess,
						},
					},
				})
				if err != nil {
					return err
				}

				tr := trade.NewTrade()
				tr.FromBson(h)
				render := trade.RenderTrade(tr)

				if serverSync {
					cur, err := b.Db.Collection("serversync").Find(context.TODO(), bson.M{})
					if err != nil {
						return err
					}
					for cur.Next(context.TODO()) {
						var ss database.ServerSync
						err = cur.Decode(&ss)
						if err != nil {
							return err
						}
						_, err = b.Client.Rest().CreateWebhookMessage(snowflake.MustParse(ss.WebhookId), ss.WebhookToken, discord.WebhookMessageCreate{
							Embeds: []discord.Embed{
								{
									Title: "New trade post!",
									Description: fmt.Sprintf("A new trade has been posted by %s. You can view and reply online [here](https://tradebuilder.app/trade/%s)",
										event.User().Username, postID),
									Image: &discord.EmbedResource{
										URL: "attachment://trade.png",
									},
									Color: common.ColorPrimary,
								},
							},
							Files: []*discord.File{
								{
									Name:   "trade.png",
									Reader: render,
								},
							},
						}, false, 0)
						if err != nil {
							return err
						}
					}
				}
				return nil
			},
		},
	}
}
