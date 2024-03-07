package controller

import (
	"regexp"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var (
	pathMatch = regexp.MustCompile(`\/[a-zA-Z]+`)
)

// Router struct to implement
// http router pattern.
type Router struct {
	bot     *bot.Bot
	session Session
	prefix  string
}

// NewRouter returns new router instance.
func NewRouter(
	bot *bot.Bot,
	session Session,
) *Router {
	return &Router{
		bot:     bot,
		session: session,
		prefix:  "",
	}
}

// Prefix returns current prefix.
func (r *Router) Prefix() string {
	return r.prefix
}

// With returns router with stacked route path.
func (r *Router) With(s string) *Router {
	if !pathMatch.MatchString(s) {
		panic("invalid path " + s + "\n" + "Correct example: '/help'.")
	}

	return &Router{
		prefix:  r.prefix + s,
		session: r.session,
	}
}

// Register command registers command
// by excact matching text message.
func (r *Router) RegisterCommand(cmd string, handler bot.HandlerFunc) {
	r.bot.RegisterHandler(bot.HandlerTypeMessageText, cmd, bot.MatchTypeExact, handler)
}

// Register hanlder to given path.
func (r *Router) RegisterHandler(path string, handler bot.HandlerFunc) {
	r.bot.RegisterHandlerMatchFunc(r.matchFunc(path), handler)
}

// MatchFunc returns func providing wanted match pattern
func (r *Router) matchFunc(s string) bot.MatchFunc {
	return func(update *models.Update) bool {
		return r.prefix+s == r.session.Status(update.Message.From.ID)
	}
}
