package main

import (
	"context"
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/bot"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/webhook"
	"github.com/DisgoOrg/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	dmWebhookID    = discord.Snowflake(os.Getenv("DM_WEBHOOK_ID"))
	dmWebhookToken = os.Getenv("DM_WEBHOOK_TOKEN")

	botToken = os.Getenv("BOT_TOKEN")
	//botGuildID   = discord.Snowflake(os.Getenv("BOT_GUILD_ID"))
	botChannelID = discord.Snowflake(os.Getenv("BOT_CHANNEL_ID"))
)

type Bot struct {
	bot             *core.Bot
	dmWebhookClient *webhook.Client

	// DMChannelID -> ThreadID
	dmThreads map[discord.Snowflake]discord.Snowflake
	// ThreadID -> DMChannelID
	threadDMs map[discord.Snowflake]discord.Snowflake

	// DMMessageID -> ThreadMessageID
	dmMessageIDs map[discord.Snowflake]discord.Snowflake
	// ThreadMessageID -> DMMessageID
	threadMessageIDs map[discord.Snowflake]discord.Snowflake
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelDebug)

	httpClient := &http.Client{Timeout: 10 * time.Second}
	disgo, err := bot.New(botToken,
		bot.WithRestClientOpts(
			rest.WithHTTPClient(httpClient),
		),
		bot.WithGatewayOpts(
			gateway.WithGatewayIntents(discord.GatewayIntentGuilds|discord.GatewayIntentGuildMessages|discord.GatewayIntentGuildMessageTyping|discord.GatewayIntentDirectMessages|discord.GatewayIntentDirectMessageTyping),
		),
	)
	if err != nil {
		log.Fatal("Error creating bot: ", err)
	}

	webhookClient := webhook.NewClient(dmWebhookID, dmWebhookToken,
		webhook.WithRestClientConfigOpts(
			rest.WithHTTPClient(httpClient),
		),
	)

	dmThreadBot := &Bot{
		bot:              disgo,
		dmWebhookClient:  webhookClient,
		dmThreads:        make(map[discord.Snowflake]discord.Snowflake),
		threadDMs:        make(map[discord.Snowflake]discord.Snowflake),
		dmMessageIDs:     make(map[discord.Snowflake]discord.Snowflake),
		threadMessageIDs: make(map[discord.Snowflake]discord.Snowflake),
	}

	disgo.AddEventListeners(&events.ListenerAdapter{
		OnDMMessageCreate:   dmMessageCreateListener(dmThreadBot),
		OnDMMessageUpdate:   dmMessageUpdateListener(dmThreadBot),
		OnDMMessageDelete:   dmMessageDeleteListener(dmThreadBot),
		OnDMUserTypingStart: dmUserTypingStartListener(dmThreadBot),

		OnGuildMessageCreate:     guildMessageCreateListener(dmThreadBot),
		OnGuildMessageUpdate:     guildMessageUpdateListener(dmThreadBot),
		OnGuildMessageDelete:     guildMessageDeleteListener(dmThreadBot),
		OnGuildMemberTypingStart: guildMemberTypingStartListener(dmThreadBot),
	})

	if err = disgo.ConnectGateway(context.Background()); err != nil {
		log.Fatal("Error connecting to gateway: ", err)
	}

	defer disgo.Close(context.Background())

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}
