package tradeplugin

import (
	"alaninnovates.com/trade-builder/common"
	"alaninnovates.com/trade-builder/common/loaders"
	"alaninnovates.com/trade-builder/tradeplugin/trade"
	"context"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
)

func TradeCommand(b *common.Bot, tradeService *State) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name:        "trade",
			Description: "Trade items",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "create",
					Description: "Start building your trade",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "lookingfor",
					Description: "Add a sticker you are looking for",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:         "sticker_name",
							Description:  "The name of the sticker",
							Required:     true,
							Autocomplete: true,
						},
						discord.ApplicationCommandOptionInt{
							Name:        "quantity",
							Description: "The quantity of the sticker",
							Required:    true,
							MinValue:    json.Ptr(1),
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "offering",
					Description: "Add a sticker you are offering",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:         "sticker_name",
							Description:  "The name of the sticker",
							Required:     true,
							Autocomplete: true,
						},
						discord.ApplicationCommandOptionInt{
							Name:        "quantity",
							Description: "The quantity of the sticker",
							Required:    true,
							MinValue:    json.Ptr(1),
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "remove",
					Description: "Remove a sticker from your trade",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "type",
							Description: "The type of sticker to remove",
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
							Name:         "sticker_name",
							Description:  "The name of the sticker",
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "view",
					Description: "View your trade",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "info",
					Description: "Get a copy-pasteable message of your trade",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "save",
					Description: "Save your trade",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "save_name",
							Description: "The name of the save",
							Required:    true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommandGroup{
					Name:        "saves",
					Description: "Manage your saves",
					Options: []discord.ApplicationCommandOptionSubCommand{
						{
							Name:        "list",
							Description: "List your saved trades",
						},
						{
							Name:        "load",
							Description: "Load a saved trade",
							Options: []discord.ApplicationCommandOption{
								discord.ApplicationCommandOptionString{
									Name:        "id",
									Description: "The id of the trade",
									Required:    true,
								},
							},
						},
						{
							Name:        "delete",
							Description: "Delete a saved trade",
							Options: []discord.ApplicationCommandOption{
								discord.ApplicationCommandOptionString{
									Name:        "id",
									Description: "The id of the trade",
									Required:    true,
								},
							},
						},
					},
				},
			},
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"create": func(event *events.ApplicationCommandInteractionCreate) error {
				tradeService.CreateTrade(event.User().ID)
				return event.CreateMessage(discord.MessageCreate{
					Content: "Created new hive. You can now add offers with the `/trade offering` command or looking for request with the `/trade lookingfor` command.",
				})
			},
			"lookingfor": func(event *events.ApplicationCommandInteractionCreate) error {
				t := tradeService.GetTrade(event.User().ID)
				if t == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have an active trade. Create one with the `/trade create` command.",
					})
				}
				data := event.SlashCommandInteractionData()
				stickerName, _ := data.OptString("sticker_name")
				quantity, _ := data.OptInt("quantity")
				t.AddLookingFor(stickerName, quantity)
				b.Redis.Client().Incr(context.Background(), "daily:lookingfor:"+stickerName)
				b.Redis.Client().Incr(context.Background(), "weekly:lookingfor:"+stickerName)
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title:       "Added Sticker",
							Description: "Added " + stickerName + " x" + strconv.Itoa(quantity) + " to your looking for list.",
							Color:       common.ColorSuccess,
						},
					},
				})
			},
			"offering": func(event *events.ApplicationCommandInteractionCreate) error {
				t := tradeService.GetTrade(event.User().ID)
				if t == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have an active trade. Create one with the `/trade create` command.",
					})
				}
				data := event.SlashCommandInteractionData()
				stickerName, _ := data.OptString("sticker_name")
				quantity, _ := data.OptInt("quantity")
				t.AddOffering(stickerName, quantity)
				b.Redis.Client().Incr(context.Background(), "daily:offering:"+stickerName)
				b.Redis.Client().Incr(context.Background(), "weekly:offering:"+stickerName)
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title:       "Added Sticker",
							Description: "Added " + stickerName + " x" + strconv.Itoa(quantity) + " to your offering list.",
							Color:       common.ColorSuccess,
						},
					},
				})
			},
			"remove": func(event *events.ApplicationCommandInteractionCreate) error {
				t := tradeService.GetTrade(event.User().ID)
				if t == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have an active trade. Create one with the `/trade create` command.",
					})
				}
				data := event.SlashCommandInteractionData()
				stickerName, _ := data.OptString("sticker_name")
				stickerType, _ := data.OptString("type")
				t.Remove(stickerType, stickerName)
				if stickerType == "lf" {
					b.Redis.Client().Decr(context.Background(), "daily:lookingfor:"+stickerName)
					b.Redis.Client().Decr(context.Background(), "weekly:lookingfor:"+stickerName)
				} else {
					b.Redis.Client().Decr(context.Background(), "daily:offering:"+stickerName)
					b.Redis.Client().Decr(context.Background(), "weekly:offering:"+stickerName)
				}
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title:       "Removed Sticker",
							Description: "Removed " + stickerName + " from your " + stickerType + " list.",
							Color:       common.ColorSuccess,
						},
					},
				})
			},
			"view": func(event *events.ApplicationCommandInteractionCreate) error {
				t := tradeService.GetTrade(event.User().ID)
				if t == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have an active trade. Create one with the `/trade create` command.",
					})
				}
				err := event.DeferCreateMessage(false)
				if err != nil {
					return err
				}
				r := trade.RenderTrade(t)
				_, err = b.Client.Rest().UpdateInteractionResponse(b.Client.ApplicationID(), event.Token(), discord.MessageUpdate{
					Embeds: &[]discord.Embed{
						{
							Title: event.User().Username + "'s Trade",
							Image: &discord.EmbedResource{
								URL: "attachment://trade.png",
							},
							Color: common.ColorPrimary,
						},
					},
					Files: []*discord.File{
						{
							Name:   "trade.png",
							Reader: r,
						},
					},
					Components: &[]discord.ContainerComponent{
						discord.ActionRowComponent{
							discord.ButtonComponent{
								Label:    "Add Looking For",
								CustomID: "handler:addlf:" + event.User().ID.String(),
								Style:    discord.ButtonStylePrimary,
							},
							discord.ButtonComponent{
								Label:    "Add Offer",
								CustomID: "handler:addoffer:" + event.User().ID.String(),
								Style:    discord.ButtonStylePrimary,
							},
							discord.ButtonComponent{
								Label:    "Rerender",
								CustomID: "handler:rerender:" + event.User().ID.String(),
								Style:    discord.ButtonStylePrimary,
							},
						},
					},
				})
				return err
			},
			"info": func(event *events.ApplicationCommandInteractionCreate) error {
				t := tradeService.GetTrade(event.User().ID)
				if t == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have an active trade. Create one with the `/trade create` command.",
					})
				}
				lfText := ""
				for _, lf := range t.GetLookingFor() {
					lfText += "- " + lf.Name + " x" + strconv.Itoa(lf.Quantity) + "\n"
				}
				offeringText := ""
				for _, offering := range t.GetOffering() {
					offeringText += "- " + offering.Name + " x" + strconv.Itoa(offering.Quantity) + "\n"
				}
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title:       "Trade Info",
							Description: "Looking For:\n" + lfText + "\nOffering:\n" + offeringText,
							Color:       common.ColorPrimary,
						},
					},
				})
			},
			"save": func(event *events.ApplicationCommandInteractionCreate) error {
				t := tradeService.GetTrade(event.User().ID)
				if t == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have an active trade. Create one with the `/trade create` command.",
					})
				}
				data := event.SlashCommandInteractionData()
				name, _ := data.OptString("save_name")
				if name == "" {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You need to provide a non-empty name for the trade.",
					})
				}
				userSaveCount, _ := b.Db.Collection("trades").CountDocuments(context.Background(), bson.M{"user_id": event.User().ID})
				if int(userSaveCount) >= common.MaxFreeSaves {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You have reached the maximum number of free saves.",
					})
				}
				res, err := b.Db.Collection("trades").UpdateOne(context.Background(), bson.M{
					"user_id": event.User().ID,
					"name":    name,
				}, bson.D{{
					"$set",
					t.ToBson(),
				}}, options.Update().SetUpsert(true))
				if err != nil {
					b.Logger.Error("Error saving trade: ", err)
				}
				id := res.UpsertedID
				if id == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "Updated save.",
					})
				}
				oid, _ := id.(primitive.ObjectID)
				hiveId, _ := oid.MarshalText()
				return event.CreateMessage(discord.MessageCreate{
					Content: "Saved trade. ID: `" + string(hiveId) + "`",
				})
			},
			"saves/load": func(event *events.ApplicationCommandInteractionCreate) error {
				data := event.SlashCommandInteractionData()
				id, _ := data.OptString("id")
				if id == "" {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You need to provide a non-empty id for the trade.",
					})
				}
				oid, err := primitive.ObjectIDFromHex(id)
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
				userTrade := tradeService.CreateTrade(event.User().ID)
				userTrade.FromBson(h)
				return event.CreateMessage(discord.MessageCreate{
					Content: "Loaded trade.",
				})
			},
			"saves/list": func(event *events.ApplicationCommandInteractionCreate) error {
				var results []bson.D
				cur, _ := b.Db.Collection("trades").Find(context.Background(), bson.M{"user_id": event.User().ID})
				err := cur.All(context.Background(), &results)
				if err != nil {
					b.Logger.Error("Failed to list trade saves for user: %v", err)
				}
				if len(results) == 0 {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have any saves.",
					})
				}
				var saves []string
				var rows []discord.ContainerComponent
				row := discord.ActionRowComponent{}
				for i, result := range results {
					id, _ := result.Map()["_id"].(primitive.ObjectID).MarshalText()
					name := result.Map()["name"].(string)
					saves = append(saves, fmt.Sprintf("%d. %s (`%s`)", i+1, name, id))
					row = row.AddComponents(discord.NewPrimaryButton(name, fmt.Sprintf("handler:save-id:%s:%s", event.User().ID.String(), id)))
					if (i+1)%5 == 0 {
						rows = append(rows, row)
						row = discord.ActionRowComponent{}
					}
				}
				if len(row.Components()) > 0 {
					rows = append(rows, row)
				}
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title:       "Your Saves",
							Description: strings.Join(saves, "\n"),
							Footer: &discord.EmbedFooter{
								Text: "Press the buttons below to get mobile friendly ids",
							},
						},
					},
					Components: rows,
				})
			},
			"saves/delete": func(event *events.ApplicationCommandInteractionCreate) error {
				data := event.SlashCommandInteractionData()
				id, _ := data.OptString("id")
				if id == "" {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You need to provide the ID of the save you want to delete.",
					})
				}
				oid, err := primitive.ObjectIDFromHex(id)
				if err != nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "Invalid ID.",
					})
				}
				res, err := b.Db.Collection("trades").DeleteOne(context.Background(), bson.M{
					"_id":     oid,
					"user_id": event.User().ID,
				})
				if err != nil {
					b.Logger.Error("Error deleting trade save: ", err)
				}
				if res.DeletedCount == 0 {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have a save with that ID.",
					})
				}
				return event.CreateMessage(discord.MessageCreate{
					Content: "Deleted save.",
				})
			},
		},
		AutocompleteHandlers: map[string]handler.AutocompleteHandler{
			"lookingfor": makeAutocompleteHandler(loaders.GetAllStickers(), "sticker_name"),
			"offering":   makeAutocompleteHandler(loaders.GetAllStickers(), "sticker_name"),
			"remove":     makeAutocompleteHandler(loaders.GetAllStickers(), "sticker_name"),
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

func AddLookingForButton() handler.Component {
	return handler.Component{
		Name:  "addlf",
		Check: common.UserIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			return event.Modal(discord.ModalCreate{
				Title:    "Add Looking For",
				CustomID: "handler:addlfmodal",
				Components: []discord.ContainerComponent{
					discord.ActionRowComponent{
						discord.TextInputComponent{
							CustomID: "name",
							Style:    discord.TextInputStyleShort,
							Label:    "Name",
						},
					},
					discord.ActionRowComponent{
						discord.TextInputComponent{
							CustomID:  "quantity",
							Style:     discord.TextInputStyleShort,
							Label:     "Quantity",
							MinLength: json.Ptr(1),
							MaxLength: 2,
						},
					},
				},
			})
		},
	}
}

func AddLookingForModal(b *common.Bot, tradeService *State) handler.Modal {
	return handler.Modal{
		Name: "addlfmodal",
		Handler: func(event *events.ModalSubmitInteractionCreate) error {
			quantityStr := event.Data.Text("quantity")
			quantityInt, err := strconv.Atoi(quantityStr)
			if err != nil || quantityInt < 1 {
				return event.CreateMessage(discord.MessageCreate{
					Content: "Quantity must be an integer",
					Flags:   discord.MessageFlagEphemeral,
				})
			}
			queryStickerName := event.Data.Text("name")
			sticker := ""
			for _, s := range loaders.GetAllStickers() {
				if strings.Contains(strings.ToLower(s), strings.ToLower(queryStickerName)) {
					sticker = s
					break
				}
			}
			if sticker == "" {
				return event.CreateMessage(discord.MessageCreate{
					Content: "Sticker `" + queryStickerName + "` not found",
					Flags:   discord.MessageFlagEphemeral,
				})
			}
			t := tradeService.GetTrade(event.User().ID)
			if t == nil {
				return event.CreateMessage(discord.MessageCreate{
					Content: "Your trade seems to have gone missing... Create a new one with `/trade create`",
					Flags:   discord.MessageFlagEphemeral,
				})
			}
			t.AddLookingFor(sticker, quantityInt)
			b.Redis.Client().Incr(context.Background(), "daily:lookingfor:"+sticker)
			b.Redis.Client().Incr(context.Background(), "weekly:lookingfor:"+sticker)
			return event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{
					{
						Title:       "Added Sticker",
						Description: "Added " + sticker + " x" + strconv.Itoa(quantityInt) + " to your looking for list.",
						Color:       common.ColorSuccess,
					},
				},
				Flags: discord.MessageFlagEphemeral,
			})
		},
	}
}

func AddOfferButton() handler.Component {
	return handler.Component{
		Name:  "addoffer",
		Check: common.UserIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			return event.Modal(discord.ModalCreate{
				Title:    "Add Offer",
				CustomID: "handler:addoffermodal",
				Components: []discord.ContainerComponent{
					discord.ActionRowComponent{
						discord.TextInputComponent{
							CustomID: "name",
							Style:    discord.TextInputStyleShort,
							Label:    "Name",
						},
					},
					discord.ActionRowComponent{
						discord.TextInputComponent{
							CustomID:  "quantity",
							Style:     discord.TextInputStyleShort,
							Label:     "Quantity",
							MinLength: json.Ptr(1),
							MaxLength: 2,
						},
					},
				},
			})
		},
	}
}

func AddOfferModal(b *common.Bot, tradeService *State) handler.Modal {
	return handler.Modal{
		Name: "addoffermodal",
		Handler: func(event *events.ModalSubmitInteractionCreate) error {
			quantityStr := event.Data.Text("quantity")
			quantityInt, err := strconv.Atoi(quantityStr)
			if err != nil || quantityInt < 1 {
				return event.CreateMessage(discord.MessageCreate{
					Content: "Quantity must be an integer",
					Flags:   discord.MessageFlagEphemeral,
				})
			}
			queryStickerName := event.Data.Text("name")
			sticker := ""
			for _, s := range loaders.GetAllStickers() {
				if strings.Contains(strings.ToLower(s), strings.ToLower(queryStickerName)) {
					sticker = s
					break
				}
			}
			if sticker == "" {
				return event.CreateMessage(discord.MessageCreate{
					Content: "Sticker `" + queryStickerName + "` not found",
					Flags:   discord.MessageFlagEphemeral,
				})
			}
			t := tradeService.GetTrade(event.User().ID)
			if t == nil {
				return event.CreateMessage(discord.MessageCreate{
					Content: "Your trade seems to have gone missing... Create a new one with `/trade create`",
					Flags:   discord.MessageFlagEphemeral,
				})
			}
			t.AddOffering(sticker, quantityInt)
			b.Redis.Client().Incr(context.Background(), "daily:offering:"+sticker)
			b.Redis.Client().Incr(context.Background(), "weekly:offering:"+sticker)
			return event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{
					{
						Title:       "Added Sticker",
						Description: "Added " + sticker + " x" + strconv.Itoa(quantityInt) + " to your looking for list.",
						Color:       common.ColorSuccess,
					},
				},
				Flags: discord.MessageFlagEphemeral,
			})
		},
	}
}

func RerenderButton(b *common.Bot, tradeService *State) handler.Component {
	return handler.Component{
		Name: "rerender",
		Handler: func(event *events.ComponentInteractionCreate) error {
			data := strings.Split(event.ButtonInteractionData().CustomID(), ":")
			uid := data[2]
			userId, _ := snowflake.Parse(uid)
			t := tradeService.GetTrade(userId)
			if t == nil {
				return event.UpdateMessage(discord.MessageUpdate{
					Content:     json.Ptr("Your trade seems to have gone missing... Create a new one with `/trade create`"),
					Embeds:      &[]discord.Embed{},
					Components:  &[]discord.ContainerComponent{},
					Attachments: &[]discord.AttachmentUpdate{},
				})
			}
			err := event.DeferCreateMessage(true)
			if err != nil {
				return err
			}
			message := event.Message
			r := trade.RenderTrade(t)
			_, err = b.Client.Rest().UpdateMessage(message.ChannelID, message.ID, discord.MessageUpdate{
				Files: []*discord.File{
					{
						Name:   "trade.png",
						Reader: r,
					},
				},
			})
			if err != nil {
				b.Logger.Error("Failed to re-render trade: ", err)
				cause := ""
				if strings.Contains(err.Error(), "403") || strings.Contains(err.Error(), "50001") {
					cause = "I didn't have permission to edit that message!"
				} else {
					cause = "Something went wrong!"
				}
				_, err = b.Client.Rest().UpdateInteractionResponse(b.Client.ApplicationID(), event.Token(), discord.MessageUpdate{
					Content: json.Ptr("Failed to re-render trade: " + cause),
				})
				return err
			}
			_, err = b.Client.Rest().UpdateInteractionResponse(b.Client.ApplicationID(), event.Token(), discord.MessageUpdate{
				Content: json.Ptr("Rerendered trade!"),
			})
			return err
		},
		Check: common.UserIDCheck(),
	}
}

func SaveIdButton() handler.Component {
	return handler.Component{
		Name:  "save-id",
		Check: common.UserIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			id := strings.Split(event.ButtonInteractionData().CustomID(), ":")[3]
			return event.CreateMessage(discord.MessageCreate{Content: id})
		},
	}
}

func Initialize(h *handler.Handler, b *common.Bot) {
	tradeService := NewTradeService()
	h.AddCommands(TradeCommand(b, tradeService))
	h.AddComponents(AddLookingForButton(), AddOfferButton(), RerenderButton(b, tradeService), SaveIdButton())
	h.AddModals(AddLookingForModal(b, tradeService), AddOfferModal(b, tradeService))
}
