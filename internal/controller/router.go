package controller

import (
	"regexp"
	"slices"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var (
	pathMatch       = regexp.MustCompile(`[a-zA-Z]+`)
	delimiter       = "/"
	prefixDelimiter = "_"
)

// Router struct to implement
// http router pattern.
type Router struct {
	bot     *bot.Bot
	session Session
	prefix  string
	paths   []string
}

// NewRouter returns new router instance.
func NewRouter(
	bot *bot.Bot,
	session Session,
) *Router {
	return &Router{
		bot:     bot,
		session: session,
		prefix:  delimiter,
	}
}

// Prefix returns current prefix.
func (r *Router) Prefix() string {
	return r.prefix
}

// With returns router with stacked route path.
func (r *Router) With(cmd Command) *Router {
	if !pathMatch.MatchString(string(cmd)) {
		panic("invalid path " + cmd + "\n" + "Must contain only latin letters.")
	}

	if slices.Contains(r.paths, string(cmd)) {
		panic("doubled path: " + string(cmd) + ". A the point " + r.prefix)
	}

	if strings.Contains(string(cmd), prefixDelimiter) {
		panic("detected forbidden symbol +'" + prefixDelimiter + "'.")
	}

	return &Router{
		prefix:  r.prefix + delimiter + string(cmd),
		session: r.session,
	}
}

// Register command registers command
// by excact matching text message.
func (r *Router) RegisterCommand(handler bot.HandlerFunc) {
	if !pathMatch.MatchString(r.prefix) {
		panic("can't register command to given path: " + r.prefix)
	}

	r.bot.RegisterHandler(bot.HandlerTypeMessageText, r.prefix, bot.MatchTypeExact, handler)
}

// Register hanlder to given cmd.
func (r *Router) RegisterHandler(cmd Command, handler bot.HandlerFunc) {
	if strings.Contains(string(cmd), prefixDelimiter) {
		panic("detected forbidden symbol +'" + prefixDelimiter + "'.")
	}

	r.bot.RegisterHandlerMatchFunc(r.matchFunc(cmd), handler)
}

// RegisterCallback registers callback to given path.
func (r *Router) RegisterCallback(cmd Command, handler bot.HandlerFunc) {
	if strings.Contains(string(cmd), prefixDelimiter) {
		panic("detected forbidden symbol +'" + prefixDelimiter + "'.")
	}

	r.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, r.callback(cmd), bot.MatchTypeExact, handler)
}

// RegisterCallbackPrefix registers callback that
// supports encoding state in callback.
func (r *Router) RegisterCallbackPrefix(cmd Command, handler bot.HandlerFunc) {
	if strings.Contains(string(cmd), prefixDelimiter) {
		panic("detected forbidden symbol +'" + prefixDelimiter + "'.")
	}

	r.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, r.callbackPrefix(cmd), bot.MatchTypePrefix, handler)
}

// Path returns absolute path
// to register redirect
func (r *Router) Path(cmd Command) string {
	return r.prefix + string(cmd)
}

// PathPrefixState returns absolute path
// to register redirect with encoded state in it.
func (r Router) PathPrefixState(cmd Command, state string) string {
	return r.callbackPrefix(cmd) + state
}

func (r *Router) GetState(callback string) string {
	_, res, _ := strings.Cut(callback, prefixDelimiter)
	return res
}

// MatchFunc returns func providing wanted match pattern
func (r *Router) matchFunc(cmd Command) bot.MatchFunc {
	return func(update *models.Update) bool {
		return r.prefix+delimiter+string(cmd) == r.session.Status(update.Message.From.ID)
	}
}

// callback returns configured callback
func (r *Router) callback(cmd Command) string {
	return r.prefix + string(cmd)
}

func (r *Router) callbackPrefix(cmd Command) string {
	return r.callback(cmd) + prefixDelimiter
}
