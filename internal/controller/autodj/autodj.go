package autodj

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/split"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/storage"
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
)

const (
	cmdBase ctr.Command = ""
	cmdBack ctr.Command = "back"

	cmdUpdate    ctr.Command = "update"
	cmdGetUpdate ctr.Command = "get-update"
	cmdReset     ctr.Command = "reset"
	cmdSubmit    ctr.Command = "submit"
	cmdCancel    ctr.Command = "cancel"

	cmdNoOp ctr.Command = "no-op"
)

type autodj struct {
	router   *ctr.Router
	dj       AutoDJ
	session  ctr.Session
	onCancel ctr.OnCancelHandler
	onError  bot.ErrorsHandler

	confStorage         storage.Storage[localModels.AutoDJConfig]
	targetUpdateStorage storage.Storage[string]
}

type AutoDJ interface {
	Config(ctx context.Context) (localModels.AutoDJConfig, error)
	SetConfig(ctx context.Context, config localModels.AutoDJConfig) error
}

func Register(
	router *ctr.Router,
	dj AutoDJ,
	session ctr.Session,
	onCancel ctr.OnCancelHandler,
	onError bot.ErrorsHandler,
) {
	a := &autodj{
		router:   router,
		dj:       dj,
		session:  session,
		onCancel: onCancel,
		onError:  onError,

		confStorage: storage.Storage[localModels.AutoDJConfig]{},
	}

	router.RegisterCallback(cmdBase, a.init)

	// settings
	router.RegisterCallbackPrefix(cmdUpdate, a.update)
	router.RegisterHandler(cmdGetUpdate, a.getUpdate)
	router.RegisterCallback(cmdReset, a.reset)

	// filler
	router.RegisterCallback(cmdNoOp, a.nullHandler)
}

func (a *autodj) init(ctx context.Context, b *bot.Bot, update *models.Update) {
	a.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	res, err := a.dj.Config(ctx)
	if err != nil {
		// handle errors
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			a.onError(err)
		}
		return
	}

	a.confStorage.Set(userId, res)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        a.configRepr(res),
		ReplyMarkup: a.mainMenuMarkup(),
		ParseMode:   models.ParseModeMarkdown,
	}); err != nil {
		a.onError(err)
	}
}

func (a *autodj) update(ctx context.Context, b *bot.Bot, update *models.Update) {
	a.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	a.targetUpdateStorage.Set(userId, update.CallbackQuery.Data)

	var msg string

	switch a.router.GetState(update.CallbackQuery.Data) {
	case "genre":
		msg = ctr.SchAutoDJAskGenre
	case "playlist":
		msg = ctr.SchAutoDJAskPlaylist
	case "language":
		msg = ctr.SchAutoDJAskLanguage
	case "mood":
		msg = ctr.SchAutoDJAskMood
	}

	a.session.Redirect(userId, a.router.Path(cmdGetUpdate))

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   msg,
	}); err != nil {
		a.onError(err)
	}
}

func (a *autodj) getUpdate(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID

	msg := update.Message.Text
	conf := a.confStorage.Get(userId)

	switch a.targetUpdateStorage.Get(userId) {
	case "genre":
		conf.Genres = split.SplitMsg(msg)
	case "playlist":
		conf.Playlists = split.SplitMsg(msg)
	case "language":
		conf.Languages = split.SplitMsg(msg)
	case "mood":
		conf.Moods = split.SplitMsg(msg)
	}

	a.targetUpdateStorage.Del(userId)
	a.confStorage.Set(userId, conf)

	a.session.Redirect(userId, ctr.NullStatus)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        a.configRepr(conf),
		ReplyMarkup: a.mainMenuMarkup(),
		ParseMode:   models.ParseModeMarkdown,
	}); err != nil {
		a.onError(err)
	}
}

func (a *autodj) reset(ctx context.Context, b *bot.Bot, update *models.Update) {
	a.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	currentConf, err := a.dj.Config(ctx)
	if err != nil {
		// handle errors
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			a.onError(err)
		}
		return
	}

	a.confStorage.Set(userId, currentConf)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        a.configRepr(currentConf),
		ReplyMarkup: a.mainMenuMarkup(),
		ParseMode:   models.ParseModeMarkdown,
	}); err != nil {
		a.onError(err)
	}
}

func (a *autodj) nullHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	a.callbackAnswer(ctx, b, update.CallbackQuery)
}

func (a *autodj) callbackAnswer(ctx context.Context, b *bot.Bot, callbackQuery *models.CallbackQuery) {
	ok, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
	})
	if err != nil {
		a.onError(err)
		return
	}
	if !ok {
		a.onError(fmt.Errorf("callback answer failed"))
	}
}

func (a *autodj) configRepr(conf localModels.AutoDJConfig) string {
	var b strings.Builder

	b.WriteString("*Настройки автодиджея:*\n")
	b.WriteString(fmt.Sprintf("*Жанры:* %s", strings.Join(conf.Genres, ", ")))
	b.WriteString(fmt.Sprintf("*Плейлисты:* %s", strings.Join(conf.Playlists, ", ")))
	b.WriteString(fmt.Sprintf("*Языки:* %s", strings.Join(conf.Languages, ", ")))
	b.WriteString(fmt.Sprintf("*Настроения:* %s", strings.Join(conf.Moods, ", ")))

	if conf.IsPlaying {
		b.WriteString("*Сейчас играет*")
	} else {
		b.WriteString("*Сейчас не играет*")
	}

	return b.String()
}
