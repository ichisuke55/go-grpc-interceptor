package i18n

import (
	"context"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	unaryInfo = &grpc.UnaryServerInfo{
		FullMethod: "TestService.UnaryMethod",
	}
)

func loadTranslation() {
	bundle := i18n.NewBundle(language.English)
	enTranslation := `[{
    "id": "hello",
    "translation": "Hello world"
  }]`
	jaTranslation := `[{
    "id": "hello",
    "translation": "こんにちは"
  }]`
	bundle.MustParseMessageFileBytes([]byte(enTranslation), "en.json")
	bundle.MustParseMessageFileBytes([]byte(jaTranslation), "ja.json")
	SetBundle(bundle)
}

func newMetadataContext(ctx context.Context, val string) context.Context {
	md := metadata.Pairs("accept-language", val)
	return metadata.NewIncomingContext(ctx, md)
}

func TestDefaultLanguage(t *testing.T) {
	loadTranslation()
	req := "request"
	_, err := UnaryServerInterceptor(context.Background(), req, unaryInfo, func(ctx context.Context, _ interface{}) (interface{}, error) {
		localizer := MustLocalizer(ctx)
		if got, want := localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "hello"}), "Hello world"; got != want {
			t.Errorf("expect localizer() = %q, but got %q", want, got)
		}
		return nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRespectAcceptLanguage(t *testing.T) {
	loadTranslation()
	req := "request"
	ctx := newMetadataContext(context.Background(), "ja")
	_, err := UnaryServerInterceptor(ctx, req, unaryInfo, func(ctx context.Context, _ interface{}) (interface{}, error) {
		localizer := MustLocalizer(ctx)
		if got, want := localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "hello"}), "こんにちは"; got != want {
			t.Errorf("expect T() = %q, but got %q", want, got)
		}
		return nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFallbackDefaultLanguage(t *testing.T) {
	loadTranslation()
	req := "request"
	ctx := newMetadataContext(context.Background(), "da")
	_, err := UnaryServerInterceptor(ctx, req, unaryInfo, func(ctx context.Context, _ interface{}) (interface{}, error) {
		localizer := MustLocalizer(ctx)
		if got, want := localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "hello"}), "Hello world"; got != want {
			t.Errorf("expect T() = %q, but got %q", want, got)
		}
		return nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
