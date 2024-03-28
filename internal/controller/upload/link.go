package upload

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
	"github.com/GintGld/fizteh-radio-bot/internal/service"
)

func (u *upload) linkUpload(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.linkUpload"

	u.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	u.session.Redirect(chatId, u.router.Path(cmdGetLink))

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   u.msgIdStorage.Get(chatId),
		Text:        ctr.LibUploadAskLink,
		ReplyMarkup: u.cancelMarkup(),
	}); err != nil {
		u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (u *upload) getLink(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.getLink"

	chatId := update.Message.Chat.ID

	msg := update.Message.Text

	inProgressMsg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   ctr.InProgress,
	})
	if err != nil {
		u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
	defer func() {
		if _, err := b.DeleteMessages(ctx, &bot.DeleteMessagesParams{
			ChatID:     chatId,
			MessageIDs: []int{update.Message.ID, inProgressMsg.ID},
		}); err != nil {
			u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
	}()

	res, err := u.mediaUpload.LinkDownload(ctx, chatId, msg)
	if err != nil {
		// Handle more errors.
		if errors.Is(err, service.ErrInvalidLink) {
			msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatId,
				Text:   ctr.LibUploadErrInvalidLink,
			})
			if err != nil {
				u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
			}
			u.msgIdStorage.Set(chatId, msg.ID)
			return
		}
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	u.session.Redirect(chatId, ctr.NullStatus)
	u.linkDownloadResStorage.Set(chatId, res)

	switch res.Type {
	case localModels.ResSong:
		u.linkTypeStorage.Set(chatId, localModels.ResSong)
		u.mediaConfigStorage.Set(chatId, res.MediaConf)

		if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatId,
			MessageID:   u.msgIdStorage.Get(chatId),
			Text:        res.MediaConf.String(),
			ReplyMarkup: u.mediaConfMarkup(res.MediaConf),
			ParseMode:   models.ParseModeHTML,
		}); err != nil {
			u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
	case localModels.ResAlbum:
		u.linkTypeStorage.Set(chatId, localModels.ResAlbum)

		if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatId,
			MessageID:   u.msgIdStorage.Get(chatId),
			Text:        u.albumRepr(res),
			ReplyMarkup: u.askUploadMarkup(),
			ParseMode:   models.ParseModeHTML,
		}); err != nil {
			u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
	case localModels.ResPlaylist:
		u.linkTypeStorage.Set(chatId, localModels.ResPlaylist)

		if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatId,
			MessageID:   u.msgIdStorage.Get(chatId),
			Text:        u.playlistRepr(res),
			ReplyMarkup: u.askUploadMarkup(),
			ParseMode:   models.ParseModeHTML,
		}); err != nil {
			u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
	}
}

func (u *upload) albumRepr(res localModels.LinkDownloadResult) string {
	var b strings.Builder

	totalDur := time.Duration(0)

	b.WriteString(fmt.Sprintf("<b>Альбом:</b> %s\n", res.Album.Name))

	for _, m := range res.Album.Values {
		totalDur += m.Duration
	}

	b.WriteString(fmt.Sprintf("<b>Количество песен:</b> %d\n", len(res.Album.Values)))
	b.WriteString(fmt.Sprintf("<b>Общая длительность:</b> %s", totalDur.String()))

	return b.String()
}

func (u *upload) playlistRepr(res localModels.LinkDownloadResult) string {
	var b strings.Builder

	totalDur := time.Duration(0)

	b.WriteString(fmt.Sprintf("<b>Плейлист:</b> %s\n", res.Playlist.Name))

	for _, m := range res.Playlist.Values {
		totalDur += m.Duration
	}

	b.WriteString(fmt.Sprintf("<b>Количество песен:</b> %d\n", len(res.Playlist.Values)))
	b.WriteString(fmt.Sprintf("<b>Общая длительность:</b> %s", totalDur.String()))

	return b.String()
}
