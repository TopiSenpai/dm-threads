package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"sync"
)

func FilesFromAttachments(bot *core.Bot, attachments []discord.Attachment) []*discord.File {
	var wg sync.WaitGroup
	files := make([]*discord.File, len(attachments))
	for ii := range attachments {
		wg.Add(1)
		i := ii
		go func() {
			defer wg.Done()
			rs, err := bot.RestServices.HTTPClient().Get(attachments[i].URL)
			if err != nil {
				bot.Logger.Error("failed to get attachment: ", err)
				return
			}
			files[i] = discord.NewFile(attachments[i].Filename, rs.Body)
		}()
	}
	wg.Wait()
	return files
}
