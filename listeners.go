package main

import "github.com/DisgoOrg/disgo/core/events"

func dmMessageCreateListener(bot *Bot) func(event *events.DMMessageCreateEvent) {
	return func(event *events.DMMessageCreateEvent) {
		threadID, ok := bot.userThreads[event.ChannelID]
		if !ok {
			bot
		}
	}
}

func guildMessageCreateListener(bot *Bot) func(event *events.GuildMessageCreateEvent) {
	return func(event *events.GuildMessageCreateEvent) {

	}
}
