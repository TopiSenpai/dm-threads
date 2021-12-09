package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

func generateEmbeds(message *core.Message) []discord.Embed {
	embeds := make([]discord.Embed, len(message.Embeds)+1)
	embeds[0] = discord.Embed{
		Author: &discord.EmbedAuthor{
			Name:    message.Author.Tag(),
			IconURL: message.Author.EffectiveAvatarURL(1024),
		},
		Description: message.Content,
	}

	for i := range message.Embeds {
		if len(embeds) == 10 {
			break
		}
		embeds[i+1] = message.Embeds[i]
	}
	return embeds
}

func guildMessageCreateListener(bot *Bot) func(event *events.GuildMessageCreateEvent) {
	return func(event *events.GuildMessageCreateEvent) {
		if event.Message.IsWebhookMessage() {
			return
		}
		dmID, ok := bot.threadDMs[event.ChannelID]
		if !ok {
			return
		}
		messageCreate := discord.MessageCreate{
			Embeds: generateEmbeds(event.Message),
			Files:  FilesFromAttachments(event.Bot(), event.Message.Attachments),
		}

		message, err := event.Bot().RestServices.ChannelService().CreateMessage(dmID, messageCreate)
		if err != nil {
			event.Bot().Logger.Error("failed to create dm message: ", err)
			return
		}
		bot.dmMessageIDs[event.Message.ID] = message.ID
	}
}

func guildMessageUpdateListener(bot *Bot) func(event *events.GuildMessageUpdateEvent) {
	return func(event *events.GuildMessageUpdateEvent) {
		dmMessageID, ok := bot.dmMessageIDs[event.Message.ID]
		if !ok {
			return
		}
		embeds := generateEmbeds(event.Message)
		messageUpdate := discord.MessageUpdate{
			Embeds: &embeds,
			Files:  FilesFromAttachments(event.Bot(), event.Message.Attachments),
		}
		dmChannelID := bot.threadDMs[event.ChannelID]
		_, err := event.Bot().RestServices.ChannelService().UpdateMessage(dmChannelID, dmMessageID, messageUpdate)
		if err != nil {
			event.Bot().Logger.Error("failed to update dm message: ", err)
			return
		}
	}
}

func guildMessageDeleteListener(bot *Bot) func(event *events.GuildMessageDeleteEvent) {
	return func(event *events.GuildMessageDeleteEvent) {
		dmMessageID, ok := bot.dmMessageIDs[event.MessageID]
		if !ok {
			return
		}
		delete(bot.threadMessageIDs, event.Message.ID)
		dmChannelID := bot.threadDMs[event.ChannelID]
		if err := event.Bot().RestServices.ChannelService().DeleteMessage(dmChannelID, dmMessageID); err != nil {
			event.Bot().Logger.Error("failed to delete dm message: ", err)
			return
		}
	}
}

func guildMemberTypingStartListener(bot *Bot) func(event *events.GuildMemberTypingStartEvent) {
	return func(event *events.GuildMemberTypingStartEvent) {
		println("thread typing: ", event.ChannelID)
		dmChannelID, ok := bot.threadDMs[event.ChannelID]
		if !ok {
			return
		}
		if err := event.Bot().RestServices.ChannelService().SendTyping(dmChannelID); err != nil {
			event.Bot().Logger.Error("failed to send dm typing: ", err)
			return
		}
	}
}
