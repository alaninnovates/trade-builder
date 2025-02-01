package main

import (
	"alaninnovates.com/trade-builder/common"
	"alaninnovates.com/trade-builder/database"
	"alaninnovates.com/trade-builder/marketplugin"
	"alaninnovates.com/trade-builder/miscplugin"
	"alaninnovates.com/trade-builder/statsplugin"
	"alaninnovates.com/trade-builder/syncplugin"
	"alaninnovates.com/trade-builder/tradeplugin"
	"context"
	"flag"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	useEnvFilePtr := flag.Bool("env", false, "a bool")
	devPtr := flag.Bool("dev", false, "a bool")
	syncCommandsPtr := flag.Bool("sync", false, "a bool")
	flag.Parse()
	if *useEnvFilePtr {
		if *devPtr {
			err := godotenv.Load(".env.dev")
			if err != nil {
				logger.Error("Failed to load .env.dev")
				panic(err)
			}
		} else {
			err := godotenv.Load(".env")
			if err != nil {
				logger.Error("Failed to load .env")
				panic(err)
			}
		}
	}
	devMode := *devPtr
	syncCommands := *syncCommandsPtr

	var (
		token    = os.Getenv("TOKEN")
		dbUri    = os.Getenv("MONGODB_URI")
		redisUri = os.Getenv("REDIS_URI")
	)

	tradeBuilder := &common.Bot{
		Logger: logger,
		Db:     *database.NewDatabase(),
		Redis:  *database.NewRedis(),
	}

	client, err := tradeBuilder.Db.Connect(dbUri)
	if err != nil {
		logger.Error("Failed to connect to database")
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	redis, err := tradeBuilder.Redis.Connect(redisUri)
	if err != nil {
		logger.Error("Failed to connect to redis")
		panic(err)
	}

	defer func() {
		if err := redis.Close(); err != nil {
			panic(err)
		}
	}()

	h := handler.New(log.New(log.LstdFlags | log.Lshortfile))
	tradeplugin.Initialize(h, tradeBuilder)
	marketplugin.Initialize(h, tradeBuilder)
	syncplugin.Initialize(h, tradeBuilder)
	miscplugin.Initialize(h, tradeBuilder)

	if tradeBuilder.Client, err = disgo.New(token,
		bot.WithShardManagerConfigOpts(
			sharding.WithGatewayConfigOpts(
				gateway.WithIntents(gateway.IntentGuilds),
				gateway.WithLogger(logger),
			),
			sharding.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds),
		),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			logger.Info("Bot is ready")
			user, _ := tradeBuilder.Client.Caches().SelfUser()
			data, _ := http.Get(user.EffectiveAvatarURL())
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					logger.Error("Failed to close body")
				}
			}(data.Body)
			bytes, _ := io.ReadAll(data.Body)
			tradeBuilder.BotInfo = common.BotInfo{
				Username: user.Username,
				AvatarIcon: &discord.Icon{
					Type: discord.IconTypeWEBP,
					Data: bytes,
				},
			}
			logger.Info("Fetched bot avatar")
		}),
		bot.WithEventListeners(h),
	); err != nil {
		logger.Error("Failed to create disgo client")
		panic(err)
	}

	statsplugin.Initialize(h, tradeBuilder, devMode)

	if !devMode && syncCommands {
		h.SyncCommands(tradeBuilder.Client)
	}
	if devMode && syncCommands {
		h.SyncCommands(tradeBuilder.Client, snowflake.GetEnv("GUILD_ID"))
		//_, _ = tradeBuilder.Client.Rest().SetGlobalCommands(tradeBuilder.Client.ApplicationID(), []discord.ApplicationCommandCreate{})
	}

	if err = tradeBuilder.Client.OpenShardManager(context.Background()); err != nil {
		logger.Error("Failed to open shard manager")
		panic(err)
	}

	logger.Info("Trade Builder is running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
