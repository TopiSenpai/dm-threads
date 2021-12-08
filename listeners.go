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
			}
		}
	}
}

func guildMessageCreateListener(bot *Bot) func(event *events.GuildMessageCreateEvent) {
	return func(event *events.GuildMessageCreateEvent) {

	}
}
