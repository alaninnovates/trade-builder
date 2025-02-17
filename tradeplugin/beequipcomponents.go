package tradeplugin

import (
	"alaninnovates.com/trade-builder/common"
	"alaninnovates.com/trade-builder/common/bssdata"
	"alaninnovates.com/trade-builder/common/loaders"
	"alaninnovates.com/trade-builder/tradeplugin/trade"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/json"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strconv"
	"strings"
)

var MissingMessage = discord.MessageUpdate{
	Content:    json.Ptr("Your trade seems to have gone missing... Create a new one with `/trade create`"),
	Embeds:     &[]discord.Embed{},
	Components: &[]discord.ContainerComponent{},
}

func GenerateBeequipMessage(userId string, stepTypePlural string, options []string) discord.MessageUpdate {
	stepTypeCapitalized := cases.Title(language.English).String(stepTypePlural)
	var optSelectMenuOptions []discord.StringSelectMenuOption
	for _, opt := range options {
		optSelectMenuOptions = append(optSelectMenuOptions, discord.StringSelectMenuOption{
			Label: opt,
			Value: opt,
		})
	}
	return discord.MessageUpdate{
		Embeds: &[]discord.Embed{
			{
				Title:       "Select " + stepTypeCapitalized,
				Description: fmt.Sprintf("Select the %s for your beequip.\n\nCurrently selected %s: None", stepTypePlural, stepTypePlural),
				Color:       common.ColorPrimary,
			},
		},
		Components: &[]discord.ContainerComponent{
			discord.ActionRowComponent{
				discord.StringSelectMenuComponent{
					CustomID:    fmt.Sprintf("handler:beequip-info-number:%s:%s", userId, stepTypePlural),
					Placeholder: "Select " + stepTypeCapitalized,
					Options:     optSelectMenuOptions,
				},
			},
			discord.ActionRowComponent{
				discord.ButtonComponent{
					CustomID: fmt.Sprintf("handler:beequip-confirm:%s", userId),
					Label:    "Confirm " + stepTypeCapitalized,
					Style:    discord.ButtonStylePrimary,
				},
			},
		},
	}
}

func GenerateAbilityMessage(event *events.ComponentInteractionCreate, beequip trade.Beequip) discord.MessageUpdate {
	abilities := loaders.GetBeequipAbility(beequip.Name)
	var optSelectMenuOptions []discord.StringSelectMenuOption
	for _, opt := range abilities {
		optSelectMenuOptions = append(optSelectMenuOptions, discord.StringSelectMenuOption{
			Label: opt,
			Value: opt,
		})
	}
	return discord.MessageUpdate{
		Embeds: &[]discord.Embed{
			{
				Title:       "Select Abilities",
				Description: "Select the abilities for your beequip.\n\nCurrently selected abilities: None",
				Color:       common.ColorPrimary,
			},
		},
		Components: &[]discord.ContainerComponent{
			discord.ActionRowComponent{
				discord.StringSelectMenuComponent{
					CustomID:    fmt.Sprintf("handler:beequip-info-ability:%s", event.User().ID.String()),
					Placeholder: "Select abilities",
					Options:     optSelectMenuOptions,
					MinValues:   json.Ptr(0),
					MaxValues:   len(optSelectMenuOptions),
				},
			},
			discord.ActionRowComponent{
				discord.ButtonComponent{
					CustomID: fmt.Sprintf("handler:beequip-confirm:%s", event.User().ID.String()),
					Label:    "Confirm Abilities",
					Style:    discord.ButtonStylePrimary,
				},
			},
		},
	}
}

func GenerateWaxMessage(event *events.ComponentInteractionCreate) discord.MessageUpdate {
	return discord.MessageUpdate{
		Embeds: &[]discord.Embed{
			{
				Title:       "Select Waxes",
				Description: "Select the waxes for your beequip.\n\nCurrently selected wax order: None",
				Color:       common.ColorPrimary,
			},
		},
		Components: &[]discord.ContainerComponent{
			discord.ActionRowComponent{
				discord.StringSelectMenuComponent{
					CustomID:    fmt.Sprintf("handler:beequip-info-wax:%s", event.User().ID.String()),
					Placeholder: "Select wax",
					Options:     bssdata.WaxSelectMenuOptions,
				},
			},
			discord.ActionRowComponent{
				// todo: UNDO wax button
				discord.ButtonComponent{
					CustomID: fmt.Sprintf("handler:beequip-confirm:%s", event.User().ID.String()),
					Label:    "Confirm Waxes",
					Style:    discord.ButtonStylePrimary,
				},
			},
		},
	}
}

func ConfirmButton(tradeService *State) handler.Component {
	return handler.Component{
		Name:  "beequip-confirm",
		Check: common.UserIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			t := tradeService.GetTrade(event.User().ID)
			if t == nil {
				return event.UpdateMessage(MissingMessage)
			}
			step := t.GetBeequipInProgressStep()
			beequip := t.GetBeequipInProgressData()
			switch step {
			case trade.BeequipInProgressStepBuffs:
				t.SetBeequipInProgressStep(trade.BeequipInProgressStepDebuffs)
				debuffs := loaders.GetBeequipDebuffs(beequip.Name)
				if len(debuffs) == 0 {
					if len(loaders.GetBeequipAbility(beequip.Name)) == 0 {
						t.SetBeequipInProgressStep(trade.BeequipInProgressStepBonuses)
						return event.UpdateMessage(GenerateBeequipMessage(event.User().ID.String(), "bonuses", loaders.GetBeequipBonuses(beequip.Name)))
					}
					t.SetBeequipInProgressStep(trade.BeequipInProgressStepAbility)
					return event.UpdateMessage(GenerateAbilityMessage(event, beequip))
				}
				return event.UpdateMessage(GenerateBeequipMessage(event.User().ID.String(), "debuffs", debuffs))
			case trade.BeequipInProgressStepDebuffs:
				t.SetBeequipInProgressStep(trade.BeequipInProgressStepAbility)
				return event.UpdateMessage(GenerateAbilityMessage(event, beequip))
			case trade.BeequipInProgressStepAbility:
				t.SetBeequipInProgressStep(trade.BeequipInProgressStepBonuses)
				return event.UpdateMessage(GenerateBeequipMessage(event.User().ID.String(), "bonuses", loaders.GetBeequipBonuses(beequip.Name)))
			case trade.BeequipInProgressStepBonuses:
				t.SetBeequipInProgressStep(trade.BeequipInProgressStepWaxes)
				return event.UpdateMessage(GenerateWaxMessage(event))
			case trade.BeequipInProgressStepWaxes:
				t.SetBeequipInProgressStep(trade.BeequipInProgressStepNone)
				switch t.GetBeequipInProgressType() {
				case "lookingFor":
					t.AddLookingForBeequip(t.GetBeequipInProgressData())
				case "offering":
					t.AddOfferingBeequip(t.GetBeequipInProgressData())
				}
				t.SetBeequipInProgressData(trade.Beequip{})
				t.SetBeequipInProgressType("")
				return event.UpdateMessage(discord.MessageUpdate{
					Embeds: &[]discord.Embed{
						{
							Title:       "Beequip Added",
							Description: "Your beequip has been added! View your trade with `/trade view`",
							Color:       common.ColorSuccess,
						},
					},
					Components: &[]discord.ContainerComponent{},
				})
			}
			return nil
		},
	}
}

func AddNumberInfoButton(tradeService *State) handler.Component {
	return handler.Component{
		Name: "beequip-info-number",
		Check: func(event *events.ComponentInteractionCreate) bool {
			allow := event.User().ID.String() == strings.Split(event.StringSelectMenuInteractionData().CustomID(), ":")[2]
			if !allow {
				_ = event.CreateMessage(discord.NewMessageCreateBuilder().
					SetContent("This is not your trade!").
					SetEphemeral(true).
					Build())
			}
			return allow
		},
		Handler: func(event *events.ComponentInteractionCreate) error {
			t := tradeService.GetTrade(event.User().ID)
			if t == nil {
				return event.UpdateMessage(MissingMessage)
			}
			interactionData := event.StringSelectMenuInteractionData()
			infoGroup := strings.Split(interactionData.CustomID(), ":")[3] // i.e. "buffs"
			infoType := interactionData.Values[0]                          // i.e. "+% Energy"
			beequip := t.GetBeequipInProgressData()
			switch infoGroup {
			case "buffs":
				if beequip.Buffs == nil {
					beequip.Buffs = make(map[string]int)
				}
				beequip.Buffs[infoType] = 0
			case "debuffs":
				if beequip.Debuffs == nil {
					beequip.Debuffs = make(map[string]int)
				}
				beequip.Debuffs[infoType] = 0
			case "bonuses":
				if beequip.Bonuses == nil {
					beequip.Bonuses = make(map[string]int)
				}
				beequip.Bonuses[infoType] = 0
			}
			t.SetBeequipInProgressData(beequip)
			return event.Modal(discord.ModalCreate{
				Title:    "Add Value for " + infoType,
				CustomID: "handler:beequip-info-modal:" + infoGroup,
				Components: []discord.ContainerComponent{
					discord.ActionRowComponent{
						discord.TextInputComponent{
							CustomID:  "value",
							Style:     discord.TextInputStyleShort,
							Label:     "Value",
							MinLength: json.Ptr(1),
							MaxLength: 2,
						},
					},
				},
			})
		},
	}
}

func AddAbilityInfoButton(tradeService *State) handler.Component {
	return handler.Component{
		Name: "beequip-info-ability",
		Check: func(event *events.ComponentInteractionCreate) bool {
			allow := event.User().ID.String() == strings.Split(event.StringSelectMenuInteractionData().CustomID(), ":")[2]
			if !allow {
				_ = event.CreateMessage(discord.NewMessageCreateBuilder().
					SetContent("This is not your trade!").
					SetEphemeral(true).
					Build())
			}
			return allow
		},
		Handler: func(event *events.ComponentInteractionCreate) error {
			t := tradeService.GetTrade(event.User().ID)
			if t == nil {
				return event.UpdateMessage(MissingMessage)
			}
			interactionData := event.StringSelectMenuInteractionData()
			beequip := t.GetBeequipInProgressData()
			if beequip.Ability == nil {
				beequip.Ability = make(map[string]bool)
			}
			abilityValuesStr := ""
			for _, ability := range interactionData.Values {
				beequip.Ability[ability] = true
				abilityValuesStr += ":white_check_mark: " + ability + "\n"
			}
			for k := range beequip.Ability {
				for _, ability := range loaders.GetBeequipAbility(beequip.Name) {
					if k != ability {
						abilityValuesStr += ":x: " + ability + "\n"
					}
				}
			}
			t.SetBeequipInProgressData(beequip)
			return event.UpdateMessage(discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					{
						Title:       "Select Abilities",
						Description: fmt.Sprintf("Select the abilities for your beequip.\n\nCurrently selected abilities:\n%s", abilityValuesStr),
						Color:       common.ColorPrimary,
					},
				},
			})
		},
	}
}

func AddWaxInfoButton(tradeService *State) handler.Component {
	return handler.Component{
		Name: "beequip-info-wax",
		Check: func(event *events.ComponentInteractionCreate) bool {
			allow := event.User().ID.String() == strings.Split(event.StringSelectMenuInteractionData().CustomID(), ":")[2]
			if !allow {
				_ = event.CreateMessage(discord.NewMessageCreateBuilder().
					SetContent("This is not your trade!").
					SetEphemeral(true).
					Build())
			}
			return allow
		},
		Handler: func(event *events.ComponentInteractionCreate) error {
			t := tradeService.GetTrade(event.User().ID)
			if t == nil {
				return event.UpdateMessage(MissingMessage)
			}
			interactionData := event.StringSelectMenuInteractionData()
			beequip := t.GetBeequipInProgressData()
			if beequip.Waxes == nil {
				beequip.Waxes = []string{}
			}
			selectedWax := interactionData.Values[0]
			beequip.Waxes = append(beequip.Waxes, selectedWax)
			t.SetBeequipInProgressData(beequip)
			waxOrderStr := ""
			for _, wax := range beequip.Waxes {
				switch wax {
				case "Soft Wax":
					waxOrderStr += bssdata.SoftWaxEmoji + " "
				case "Hard Wax":
					waxOrderStr += bssdata.HardWaxEmoji + " "
				case "Caustic Wax":
					waxOrderStr += bssdata.CausticWaxEmoji + " "
				case "Swirled Wax":
					waxOrderStr += bssdata.SwirledWaxEmoji + " "
				}
			}
			return event.UpdateMessage(discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					{
						Title:       "Select Waxes",
						Description: "Select the waxes for your beequip.\n\nCurrently selected wax order:\n" + waxOrderStr,
						Color:       common.ColorPrimary,
					},
				},
			})
		},
	}
}

func AddBeequipInfoModal(tradeService *State) handler.Modal {
	return handler.Modal{
		Name: "beequip-info-modal",
		Handler: func(event *events.ModalSubmitInteractionCreate) error {
			valueStr := event.Data.Text("value")
			valueInt, err := strconv.Atoi(valueStr)
			if err != nil || valueInt < 1 {
				return event.CreateMessage(discord.MessageCreate{
					Content: "Value must be an integer",
					Flags:   discord.MessageFlagEphemeral,
				})
			}
			t := tradeService.GetTrade(event.User().ID)
			if t == nil {
				return event.UpdateMessage(MissingMessage)
			}
			infoGroup := strings.Split(event.Data.CustomID, ":")[2]
			beequip := t.GetBeequipInProgressData()
			switch infoGroup {
			case "buffs":
				for k, v := range beequip.Buffs {
					if v == 0 {
						beequip.Buffs[k] = valueInt
						break
					}
				}
			case "debuffs":
				for k, v := range beequip.Debuffs {
					if v == 0 {
						beequip.Debuffs[k] = valueInt
						break
					}
				}
			case "bonuses":
				for k, v := range beequip.Bonuses {
					if v == 0 {
						beequip.Bonuses[k] = valueInt
						break
					}
				}
			}
			t.SetBeequipInProgressData(beequip)
			currSelectedBuffString := ""
			var dataArr map[string]int
			switch infoGroup {
			case "buffs":
				dataArr = beequip.Buffs
			case "debuffs":
				dataArr = beequip.Debuffs
			case "bonuses":
				dataArr = beequip.Bonuses
			}
			for k, v := range dataArr {
				currSelectedBuffString += k + ": " + strconv.Itoa(v) + "\n"
			}
			return event.UpdateMessage(discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					{
						Title:       "Select " + cases.Title(language.English).String(infoGroup),
						Description: fmt.Sprintf("Select the %s for your beequip.\n\nCurrently selected %s:\n%s", infoGroup, infoGroup, currSelectedBuffString),
						Color:       common.ColorPrimary,
					},
				},
			})
		},
	}
}
