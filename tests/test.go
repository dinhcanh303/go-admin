package tests

import (
	"context"
	"net/http"
	"os"
	"testing"

	"go-admin/internal/config"
	"go-admin/internal/wirex"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
)

const (
	baseAPI = "/api/v1"
)

var (
	app *gin.Engine
)

func init() {
	config.MustLoad("")

	_ = os.RemoveAll(config.C.Storage.DB.DSN)
	ctx := context.Background()
	injector, _, err := wirex.BuildInjector(ctx)
	if err != nil {
		panic(err)
	}

	if err := injector.Mods.Init(ctx); err != nil {
		panic(err)
	}

	app = gin.New()
	err = injector.Mods.RegisterRouters(ctx, app)
	if err != nil {
		panic(err)
	}
}

func tester(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(app),
			Jar:       httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}
