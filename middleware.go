package i18nmiddleware

import (
	"clevergo.tech/clevergo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type contextKey int

// LocalizerKey is the context key of localizer.
const LocalizerKey contextKey = 0

// Option is a function for setting middleware.
type Option func(*middleware)

// WithFormField is an option for setting form field.
func WithFormField(field string) Option {
	return func(m *middleware) {
		m.formField = field
	}
}

type middleware struct {
	bundle    *i18n.Bundle
	formField string
}

func (m *middleware) handle(next clevergo.Handle) clevergo.Handle {
	return func(c *clevergo.Context) error {
		localizer := i18n.NewLocalizer(m.bundle, c.FormValue(m.formField), c.GetHeader("Accept-Language"))
		c.WithValue(LocalizerKey, localizer)
		return next(c)
	}
}

// New returns a I18N middleware with the given bundle and optional options.
func New(bundle *i18n.Bundle, opts ...Option) clevergo.MiddlewareFunc {
	m := &middleware{
		bundle:    bundle,
		formField: "lang",
	}
	for _, opt := range opts {
		opt(m)
	}
	return m.handle
}

// Localizer returns the localizer instance from context.
func Localizer(c *clevergo.Context) *i18n.Localizer {
	return c.Value(LocalizerKey).(*i18n.Localizer)
}
