package autodj

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/slice"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/split"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/storage"
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
)

const (
	cmdUpdate       ctr.Command = "update"
	cmdGetUpdate    ctr.Command = "get-update"
	cmdReset        ctr.Command = "reset"
	cmdGetCurrConf  ctr.Command = "submit"
	cmdCancel       ctr.Command = "cancel"
	cmdSend         ctr.Command = "send"
	cmdStartStop    ctr.Command = "start-stop"
	cmdOpenCheckBox ctr.Command = "open-check"
	cmdCheckBtn     ctr.Command = "check"
	cmdCloseSubtask ctr.Command = "close-subtask"

	cmdNoOp ctr.Command = "no-op"
)

type autodj struct {
	ctr.CallbackAnswerer

	router  *ctr.Router
	auth    Auth
	dj      AutoDJ
	session ctr.Session
	onError bot.ErrorsHandler

	confStorage         storage.Storage[localModels.AutoDJInfo]
	targetUpdateStorage storage.Storage[string]
	msgIdStorage        storage.Storage[int]
}

type Auth interface {
	IsKnown(ctx context.Context, id int64) bool
}

type AutoDJ interface {
	Config(ctx context.Context, id int64) (localModels.AutoDJInfo, error)
	SetConfig(ctx context.Context, id int64, config localModels.AutoDJInfo) error
	StartAutoDJ(ctx context.Context, id int64) error
	StopAutoDJ(ctx context.Context, id int64) error
}

func Register(
	router *ctr.Router,
	auth Auth,
	dj AutoDJ,
	session ctr.Session,
	onError bot.ErrorsHandler,
) {
	a := &autodj{
		router:  router,
		auth:    auth,
		dj:      dj,
		session: session,
		onError: onError,

		confStorage:         storage.New[localModels.AutoDJInfo](),
		targetUpdateStorage: storage.New[string](),
		msgIdStorage:        storage.New[int](),
	}

	router.RegisterCommand(a.init)

	// settings
	router.RegisterCallbackPrefix(cmdUpdate, a.update)
	router.RegisterHandler(cmdGetUpdate, a.getUpdate)
	router.RegisterCallback(cmdReset, a.reset)
	router.RegisterCallbackPrefix(cmdOpenCheckBox, a.openCheckBox)
	router.RegisterCallbackPrefix(cmdCheckBtn, a.getCheckedBtn)
	router.RegisterCallback(cmdCloseSubtask, a.closeSubtask)

	// start, stop
	router.RegisterCallback(cmdStartStop, a.startStop)

	// send
	router.RegisterCallback(cmdSend, a.commitConfig)

	// filler
	router.RegisterCallback(cmdNoOp, a.nullHandler)
}

func (a *autodj) init(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "autodj.init"

	chatId := update.Message.Chat.ID

	if !a.auth.IsKnown(ctx, chatId) {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrUnknown,
		}); err != nil {
			a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	res, err := a.dj.Config(ctx, chatId)
	if err != nil {
		// handle errors
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	a.confStorage.Set(chatId, res)

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        a.configRepr(res),
		ReplyMarkup: a.mainMenuMarkup(res),
		ParseMode:   models.ParseModeHTML,
	})
	if err != nil {
		a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}

	a.msgIdStorage.Set(chatId, msg.ID)
}

func (a *autodj) update(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "autodj.update"

	a.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	a.targetUpdateStorage.Set(chatId, update.CallbackQuery.Data)

	var msg string

	target := a.router.GetState(update.CallbackQuery.Data)

	switch target {
	case "playlist":
		msg = ctr.SchAutoDJAskPlaylist
	}

	a.targetUpdateStorage.Set(chatId, target)
	a.session.Redirect(chatId, a.router.Path(cmdGetUpdate))

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    chatId,
		MessageID: a.msgIdStorage.Get(chatId),
		Text:      msg,
	}); err != nil {
		a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (a *autodj) openCheckBox(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.openCheckBox"

	a.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID
	callback := a.router.GetState(update.CallbackQuery.Data)

	var (
		msg    string
		markup models.InlineKeyboardMarkup
	)

	conf := a.confStorage.Get(chatId)

	switch callback {
	case "genre":
		msg = ctr.SchAutoDJAskGenre
		markup = a.genreChooseMarkup(conf)
	case "mood":
		msg = ctr.SchAutoDJAskMood
		markup = a.moodChooseMarkup(conf)
	case "lang":
		msg = ctr.SchAutoDJAskLanguage
		markup = a.langChooseMarkup(conf)
	}

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   a.msgIdStorage.Get(chatId),
		Text:        msg,
		ReplyMarkup: markup,
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (a *autodj) getCheckedBtn(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.getCheckedBtn"

	a.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID
	callback := a.router.GetState(update.CallbackQuery.Data)

	// tagType in ('genre', 'mood', 'lang')
	// data is tag id (from 1 to GenreNumber, MoodNumber, LangNumber)
	tagType, data, found := strings.Cut(callback, "-")
	if !found {
		a.onError(fmt.Errorf("%s [%d]: invalid callback data \"%s\"", op, chatId, callback))
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}
	id, err := strconv.Atoi(data)
	if err != nil {
		a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	conf := a.confStorage.Get(chatId)

	var (
		msg    string
		markup models.InlineKeyboardMarkup
	)

	switch tagType {
	case "genre":
		conf.Genres[id-1] = !conf.Genres[id-1]
		msg = ctr.LibSearchUpdateAskGenre
		markup = a.genreChooseMarkup(conf)
	case "mood":
		conf.Moods[id-1] = !conf.Moods[id-1]
		msg = ctr.LibSearchUpdateAskMood
		markup = a.moodChooseMarkup(conf)
	case "lang":
		conf.Languages[id-1] = !conf.Languages[id-1]
		msg = ctr.LibSearchUpdateAskLang
		markup = a.langChooseMarkup(conf)
	}

	a.confStorage.Set(chatId, conf)

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   a.msgIdStorage.Get(chatId),
		Text:        msg,
		ReplyMarkup: markup,
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (a *autodj) closeSubtask(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "autodj.closeSubtask"

	a.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	info := a.confStorage.Get(chatId)

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   a.msgIdStorage.Get(chatId),
		Text:        a.configRepr(info),
		ReplyMarkup: a.mainMenuMarkup(info),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (a *autodj) startStop(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "autodj.startStop"

	a.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	conf := a.confStorage.Get(chatId)

	var err error

	switch conf.IsPlaying {
	case true:
		err = a.dj.StopAutoDJ(ctx, chatId)
	case false:
		err = a.dj.StartAutoDJ(ctx, chatId)
	}

	if err != nil {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	conf, err = a.dj.Config(ctx, chatId)
	if err != nil {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   a.msgIdStorage.Get(chatId),
		Text:        a.configRepr(conf),
		ReplyMarkup: a.mainMenuMarkup(conf),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (a *autodj) getUpdate(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "autodj.getUpdate"

	chatId := update.Message.Chat.ID

	msg := update.Message.Text
	conf := a.confStorage.Get(chatId)

	switch a.targetUpdateStorage.Get(chatId) {
	case "playlist":
		conf.Playlists = split.SplitMsg(msg)
	}

	a.targetUpdateStorage.Del(chatId)
	a.confStorage.Set(chatId, conf)

	a.session.Redirect(chatId, ctr.NullStatus)

	if _, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    chatId,
		MessageID: update.Message.ID,
	}); err != nil {
		a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   a.msgIdStorage.Get(chatId),
		Text:        a.configRepr(conf),
		ReplyMarkup: a.mainMenuMarkup(conf),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (a *autodj) reset(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "autodj.reset"

	a.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	currentConf, err := a.dj.Config(ctx, chatId)
	if err != nil {
		// TODO handle errors
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	a.confStorage.Set(chatId, currentConf)

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   a.msgIdStorage.Get(chatId),
		Text:        a.configRepr(currentConf),
		ReplyMarkup: a.mainMenuMarkup(currentConf),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (a *autodj) commitConfig(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "autodj.commitConfig"

	a.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	info := a.confStorage.Get(chatId)

	if err := a.dj.SetConfig(ctx, chatId, info); err != nil {
		// TODO handle errors
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	// handle errors
	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    chatId,
		MessageID: a.msgIdStorage.Get(chatId),
		Text:      ctr.SchAutoDJSuccess,
	}); err != nil {
		a.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (a *autodj) nullHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	a.CallbackAnswer(ctx, b, update.CallbackQuery)
}

func (a *autodj) configRepr(info localModels.AutoDJInfo) string {
	var b strings.Builder

	b.WriteString("<b>Настройки автодиджея:</b>\n")
	b.WriteString(fmt.Sprintf("<b>Жанры:</b> %s\n", slice.Join(slice.Filter(localModels.GenresAvail[:], info.Genres[:]), ", ")))
	b.WriteString(fmt.Sprintf("<b>Плейлисты:</b> %s\n", strings.Join(info.Playlists, ", ")))
	b.WriteString(fmt.Sprintf("<b>Языки:</b> %s\n", slice.Join(slice.Filter(localModels.LangsAvail[:], info.Languages[:]), ", ")))
	b.WriteString(fmt.Sprintf("<b>Настроения:</b> %s\n", slice.Join(slice.Filter(localModels.MoodsAvail[:], info.Moods[:]), ", ")))

	if info.IsPlaying {
		b.WriteString("<b>Сейчас играет</b>")
	} else {
		b.WriteString("<b>Сейчас не играет</b>")
	}

	return b.String()
}
