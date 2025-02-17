package tradeplugin

import (
	"alaninnovates.com/trade-builder/common"
	"alaninnovates.com/trade-builder/tradeplugin/trade"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/json"
	"strconv"
	"strings"
)

/*
	flow:
		send embed: select buffs with context menu. button to confirm buffs
			for each item selected, send modal to input value
			update embed with buffs
		same for debuffs/ability/bonuses
		send embed: select abilities with multi select context menu
			for each item selected, update embed with abilities
		send embed: select waxes with 5 context menus (one for each slot)
*/

var MissingMessage = discord.MessageUpdate{
	Content:    json.Ptr("Your trade seems to have gone missing... Create a new one with `/trade create`"),
	Embeds:     &[]discord.Embed{},
	Components: &[]discord.ContainerComponent{},
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
			switch step {
			case trade.BeequipInProgressStepBuffs:
				t.SetBeequipInProgressStep(trade.BeequipInProgressStepDebuffs)
				return event.UpdateMessage(discord.MessageUpdate{
					Content: json.Ptr("Select debuffs"),
					Embeds:  &[]discord.Embed{},
				})
			case trade.BeequipInProgressStepDebuffs:
				t.SetBeequipInProgressStep(trade.BeequipInProgressStepAbility)
				return event.UpdateMessage(discord.MessageUpdate{
					Content: json.Ptr("Select abilities"),
					Embeds:  &[]discord.Embed{},
				})
			case trade.BeequipInProgressStepAbility:
				t.SetBeequipInProgressStep(trade.BeequipInProgressStepBonuses)
				return event.UpdateMessage(discord.MessageUpdate{
					Content: json.Ptr("Select bonuses"),
					Embeds:  &[]discord.Embed{},
				})
			case trade.BeequipInProgressStepBonuses:
				t.SetBeequipInProgressStep(trade.BeequipInProgressStepWaxes)
				return event.UpdateMessage(discord.MessageUpdate{
					Content: json.Ptr("Select waxes"),
					Embeds:  &[]discord.Embed{},
				})
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
					Content:    json.Ptr("Beequip added!"),
					Embeds:     &[]discord.Embed{},
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
			}
			t.SetBeequipInProgressData(beequip)
			currSelectedBuffString := ""
			for k, v := range beequip.Buffs {
				currSelectedBuffString += k + ": " + strconv.Itoa(v) + "\n"
			}
			return event.UpdateMessage(discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					{
						Title:       "Select Buffs",
						Description: "Select the buffs for your beequip.\n\nCurrently selected buffs:\n" + currSelectedBuffString,
						Color:       common.ColorPrimary,
					},
				},
			})
		},
	}
}
