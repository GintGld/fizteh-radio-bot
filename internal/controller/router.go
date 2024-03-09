package controller

import (
	"regexp"
	"slices"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var (
	pathMatch = regexp.MustCompile(`[a-zA-Z]+`)
	delimiter = "/"
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
func (r *Router) With(s string) *Router {
	if !pathMatch.MatchString(s) {
		panic("invalid path " + s + "\n" + "Must contain only latin letters.")
	}

	if slices.Contains(r.paths, s) {
		panic("doubled path: " + s + ". A the point " + r.prefix)
	}

	return &Router{
		prefix:  r.prefix + delimiter + s,
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
	r.bot.RegisterHandlerMatchFunc(r.matchFunc(cmd), handler)
}

// Register callback to given path.
func (r *Router) RegisterCallback(cmd Command, handler bot.HandlerFunc) {
	r.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, r.callback(cmd), bot.MatchTypeExact, handler)
}

func (r *Router) Path(cmd Command) string {
	return r.prefix + string(cmd)
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
