package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

func dmMessageCreateListener(bot *Bot) func(event *events.DMMessageCreateEvent) {
	return func(event *events.DMMessageCreateEvent) {
		threadID, ok := bot.userThreads[event.ChannelID]
		if !ok {
			if ch, ok := event.Bot().Caches.Channels().Get(botChannelID).(core.GuildMessageChannel); ok {
				thread, err := ch.CreateThread(discord.GuildPublicThreadCreate{
					Name: event.Message.Author.Tag(),
				})
				if err != nil {
					event.Bot().Logger.Error("failed to create new thread: ", err)
					return
				}
				threadID = thread.ID()
			}
		}
		webhookMessageCreate := discord.WebhookMessageCreate{
			Content:   event.Message.Content,
			Username:  event.Message.Author.Username,
			AvatarURL: event.Message.Author.EffectiveAvatarURL(1024),
			Embeds:    event.Message.Embeds,
			Files:     make([]*discord.File, len(event.Message.Attachments)),
		}
		for i := range event.Message.Attachments {
			rs, err := event.Bot().RestServices.HTTPClient().Get(event.Message.Attachments[i].URL)
			if err != nil {
				event.Bot().Logger.Error("failed to get attachment: ", err)
				continue
			}

			webhookMessageCreate.Files[i] = &discord.File{
				Name:   event.Message.Attachments[i].Filename,
				Reader: rs.Body,
			}
		}
		message, err := bot.dmWebhookClient.CreateMessageInThread(webhookMessageCreate, threadID)
		if err != nil {
			event.Bot().Logger.Error("failed to create message: ", err)
			return
		}
		bot.userMessageIDs[event.Message.ID] = message.ID
	}
}

func dmMessageUpdateListener(bot *Bot) func(event *events.DMMessageUpdateEvent) {
	return func(event *events.DMMessageUpdateEvent) {
	}
}

func dmMessageDeleteListener(bot *Bot) func(event *events.DMMessageDeleteEvent) {
	return func(event *events.DMMessageDeleteEvent) {
	}
}

func guildMessageCreateListener(bot *Bot) func(event *events.GuildMessageCreateEvent) {
	return func(event *events.GuildMessageCreateEvent) {

	}
}

func guildMessageUpdateListener(bot *Bot) func(event *events.GuildMessageUpdateEvent) {
	return func(event *events.GuildMessageUpdateEvent) {

	}
}

func guildMessageDeleteListener(bot *Bot) func(event *events.GuildMessageDeleteEvent) {
	return func(event *events.GuildMessageDeleteEvent) {

	}
}
