package app

import (
	"context"
	"testing"
	"time"

	"github.com/lishimeng/app-starter/token"
	shutdown "github.com/lishimeng/go-app-shutdown"
)

func TestTokenValidatorBuild001(t *testing.T) {
	go func() {
		time.Sleep(100 * time.Millisecond)
		shutdown.Exit("test done")
	}()

	_ = New().Start(func(ctx context.Context, builder *ApplicationBuilder) error {
		builder.EnableTokenValidator(func(inject TokenValidatorInjectFunc) {
			prov := token.NewRedisStorage(GetCache())
			inject(prov)
		})
		return nil
	}, func(s string) {})
}
