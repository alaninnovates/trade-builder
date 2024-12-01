package common

import (
	"alaninnovates.com/trade-builder/database"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"log/slog"
)

type BotInfo struct {
	Username   string
	AvatarIcon *discord.Icon
}

type Bot struct {
	Logger  *slog.Logger
	Client  bot.Client
	BotInfo BotInfo
	Db      database.Database
	Redis   database.Redis
}
