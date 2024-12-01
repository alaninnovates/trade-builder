package statsplugin

import (
	"alaninnovates.com/trade-builder/common"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"time"
)

func Initialize(h *handler.Handler, b *common.Bot, devMode bool) {
	for b.Client == nil {
		time.Sleep(1)
	}
	b.Client.AddEventListeners(&events.ListenerAdapter{
		OnApplicationCommandInteraction: func(event *events.ApplicationCommandInteractionCreate) {
			if event.Data.Type() != discord.ApplicationCommandTypeSlash {
				return
			}
			data := event.SlashCommandInteractionData()
			guildName, ok := event.Guild()
			guild := ""
			if ok {
				guild = guildName.Name
			} else {
				guild = "Dms"
			}
			b.Logger.Info(fmt.Sprintf("%s used %s in %s", event.User().Tag(), data.CommandPath(), guild))
		},
	})
}
