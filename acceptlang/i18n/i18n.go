package i18n

import (
	"context"

	"github.com/higebu/go-grpc-interceptor/acceptlang"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"google.golang.org/grpc"
)

var defaultLanguage = "en"

var bundle *i18n.Bundle

func SetDefaultLanguage(lang string) {
	defaultLanguage = lang
}

// var _ grpc.UnaryServerInterceptor = UnaryServerInterceptor

type localizerKey struct{}

func UnaryServerInterceptor(origctx context.Context, origreq interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) grpc.UnaryServerInterceptor {
	return acceptlang.UnaryServerInterceptor(origctx, origreq, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		acceptLangs := acceptlang.FromContext(ctx)
		l := HandleI18n(acceptLangs)
		ctx = context.WithValue(ctx, localizerKey{}, l)
		return handler(ctx, req)
	})
}

func HandleI18n(acceptLangs acceptlang.AcceptLanguages) *i18n.Localizer {
	langs := acceptLangs.Languages()
	langs = append(langs, defaultLanguage)
	return i18n.NewLocalizer(bundle, langs...)
}

func MustLocalizer(ctx context.Context) *i18n.Localizer {
	l, ok := ctx.Value(localizerKey{}).(*i18n.Localizer)
	if !ok {
		panic("could not find Localizer from context")
	}
	return l
}

func SetBundle(b *i18n.Bundle) {
	bundle = b
}
