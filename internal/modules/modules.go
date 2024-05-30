package modules

import (
	"context"

	"go-admin/internal/modules/auth"
	"go-admin/internal/modules/sys"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

const (
	apiPrefix = "/api/"
)

// Collection of wire providers
var Set = wire.NewSet(
	wire.Struct(new(Modules), "*"),
	auth.Set,
	sys.Set,
)

type Modules struct {
	Auth *auth.Auth
	SYS  *sys.SYS
}

func (a *Modules) Init(ctx context.Context) error {
	if err := a.Auth.Init(ctx); err != nil {
		return err
	}
	if err := a.SYS.Init(ctx); err != nil {
		return err
	}

	return nil
}

func (a *Modules) RouterPrefixes() []string {
	return []string{
		apiPrefix,
	}
}

func (a *Modules) RegisterRouters(ctx context.Context, e *gin.Engine) error {
	gAPI := e.Group(apiPrefix)
	v1 := gAPI.Group("v1")

	if err := a.Auth.RegisterV1Routers(ctx, v1); err != nil {
		return err
	}
	if err := a.SYS.RegisterV1Routers(ctx, v1); err != nil {
		return err
	}

	return nil
}

func (a *Modules) Release(ctx context.Context) error {
	if err := a.Auth.Release(ctx); err != nil {
		return err
	}
	if err := a.SYS.Release(ctx); err != nil {
		return err
	}
	return nil
}
