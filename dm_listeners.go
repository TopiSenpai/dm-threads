package main

import (
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

func dmMessageCreateListener(bot *Bot) func(event *events.DMMessageCreateEvent) {
	return func(event *events.DMMessageCreateEvent) {
		if event.Message.Author.ID == event.Bot().ClientID {
			return
		}
		threadID, ok := bot.dmThreads[event.ChannelID]
		if !ok {
			thread, err := event.Bot().RestServices.ThreadService().CreateThread(botChannelID, discord.GuildPublicThreadCreate{
				Name:                event.Message.Author.Tag(),
				AutoArchiveDuration: discord.AutoArchiveDuration24h,
			})
			if err != nil {
				event.Bot().Logger.Error("failed to create new thread: ", err)
				return
			}
			threadID = thread.ID()
			bot.dmThreads[event.ChannelID] = thread.ID()
			bot.threadDMs[thread.ID()] = event.ChannelID
		}
		webhookMessageCreate := discord.WebhookMessageCreate{
			Content:   event.Message.Content,
			Username:  event.Message.Author.Username,
			AvatarURL: event.Message.Author.EffectiveAvatarURL(1024),
			Embeds:    event.Message.Embeds,
			Files:     FilesFromAttachments(event.Bot(), event.Message.Attachments),
		}

		message, err := bot.dmWebhookClient.CreateMessageInThread(webhookMessageCreate, threadID)
		if err != nil {
			event.Bot().Logger.Error("failed to create thread message: ", err)
			return
		}
		bot.threadMessageIDs[event.Message.ID] = message.ID
	}
}

func dmMessageUpdateListener(bot *Bot) func(event *events.DMMessageUpdateEvent) {
	return func(event *events.DMMessageUpdateEvent) {
		webhookMessageID, ok := bot.threadMessageIDs[event.Message.ID]
		if !ok {
			return
		}
		webhookMessageUpdate := discord.WebhookMessageUpdate{
			Content: &event.Message.Content,
			Embeds:  &event.Message.Embeds,
			Files:   FilesFromAttachments(event.Bot(), event.Message.Attachments),
		}
		threadID := bot.dmThreads[event.ChannelID]
		_, err := bot.dmWebhookClient.UpdateMessageInThread(webhookMessageID, webhookMessageUpdate, threadID)
		if err != nil {
			event.Bot().Logger.Error("failed to update thread message: ", err)
			return
		}
	}
}

func dmMessageDeleteListener(bot *Bot) func(event *events.DMMessageDeleteEvent) {
	return func(event *events.DMMessageDeleteEvent) {
		webhookMessageID, ok := bot.threadMessageIDs[event.MessageID]
		if !ok {
			return
		}
		delete(bot.threadMessageIDs, event.Message.ID)
		if err := bot.dmWebhookClient.DeleteMessage(webhookMessageID); err != nil {
			event.Bot().Logger.Error("failed to delete thread message: ", err)
			return
		}
	}
}

func dmUserTypingStartListener(bot *Bot) func(event *events.DMUserTypingStartEvent) {
	return func(event *events.DMUserTypingStartEvent) {
		threadID, ok := bot.dmThreads[event.ChannelID]
		if !ok {
			return
		}
		if err := event.Bot().RestServices.ChannelService().SendTyping(threadID); err != nil {
			event.Bot().Logger.Error("failed to send thread typing: ", err)
			return
		}
	}
}
