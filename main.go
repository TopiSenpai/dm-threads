package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/bot"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/webhook"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	dmWebhookID    = discord.Snowflake(os.Getenv("DM_WEBHOOK_ID"))
	dmWebhookToken = os.Getenv("DM_WEBHOOK_TOKEN")

	botToken     = os.Getenv("BOT_TOKEN")
	botGuildID   = discord.Snowflake(os.Getenv("BOT_GUILD_ID"))
	botChannelID = discord.Snowflake(os.Getenv("BOT_CHANNEL_ID"))
)

type Bot struct {
	bot             *core.Bot
	dmWebhookClient *webhook.Client

	// DMChannelID -> ThreadID
	userThreads map[discord.Snowflake]discord.Snowflake
	// DMMessageID -> ThreadMessageID
	userMessageIDs map[discord.Snowflake]discord.Snowflake
}

func main() {
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
		bot:             disgo,
		dmWebhookClient: webhookClient,
		userThreads:     make(map[discord.Snowflake]discord.Snowflake),
	}

	disgo.AddEventListeners(&events.ListenerAdapter{
		OnDMMessageCreate: dmMessageCreateListener(dmThreadBot),
		OnDMMessageUpdate: dmMessageUpdateListener(dmThreadBot),
		OnDMMessageDelete: dmMessageDeleteListener(dmThreadBot),

		OnGuildMessageCreate: guildMessageCreateListener(dmThreadBot),
		OnGuildMessageUpdate: guildMessageUpdateListener(dmThreadBot),
		OnGuildMessageDelete: guildMessageDeleteListener(dmThreadBot),
	})
}
