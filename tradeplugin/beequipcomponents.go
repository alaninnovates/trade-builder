package tradeplugin

import (
	"alaninnovates.com/trade-builder/common"
	"alaninnovates.com/trade-builder/tradeplugin/trade"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/json"
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

func ConfirmButton(tradeService *State) handler.Component {
	return handler.Component{
		Name:  "beequipconfirm",
		Check: common.UserIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			t := tradeService.GetTrade(event.User().ID)
			if t == nil {
				return event.UpdateMessage(discord.MessageUpdate{
					Content:    json.Ptr("Your trade seems to have gone missing... Create a new one with `/trade create`"),
					Embeds:     &[]discord.Embed{},
					Components: &[]discord.ContainerComponent{},
				})
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
