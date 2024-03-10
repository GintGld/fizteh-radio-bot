package upload

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
)

func (u *upload) linkUpload(ctx context.Context, b *bot.Bot, update *models.Update) {
	u.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	u.session.Redirect(userId, u.router.Path(cmdGetLink))

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        ctr.LibUploadAskLink,
		ReplyMarkup: u.getSettingDataMarkup(),
	}); err != nil {
		u.onError(err)
	}
}

func (u *upload) getLink(ctx context.Context, b *bot.Bot, update *models.Update) {
	u.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID

	msg := update.Message.Text

	res, err := u.mediaUpload.LinkDownload(msg)
	if err != nil {
		// Handle more errors.
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			u.onError(err)
		}
		return
	}

	u.linkDownloadResStorage.Set(userId, res)

	switch res.Type {
	case localModels.ResSong:
		conf := localModels.MediaConfig{
			Name:     res.Media.Name,
			Author:   res.Media.Author,
			Duration: res.Media.Duration,
		}

		u.mediaConfigStorage.Set(userId, conf)

		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      chatId,
			Text:        u.mediaConfRepr(conf),
			ReplyMarkup: u.mediaConfMarkup(conf),
		}); err != nil {
			u.onError(err)
		}
	case localModels.ResPlaylist:
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      chatId,
			Text:        u.playlistRepr(res),
			ReplyMarkup: u.playlistMarkup(),
		}); err != nil {
			u.onError(err)
		}
	}
}

func (u *upload) playlistRepr(res localModels.LinkDownloadResult) string {
	var b strings.Builder

	totalDur := time.Duration(0)

	b.WriteString(fmt.Sprintf("*Плейлист:* %s", res.Playlist.Name))

	for _, m := range res.Playlist.Values {
		totalDur += m.Duration
	}

	b.WriteString(fmt.Sprintf("*Количество песен:* %d", len(res.Playlist.Values)))
	b.WriteString(fmt.Sprintf("*Общая длительность:* %s", totalDur.String()))

	return b.String()
}
