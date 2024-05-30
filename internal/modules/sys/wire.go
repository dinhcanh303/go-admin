package sys

import (
	"go-admin/internal/modules/sys/api"
	"go-admin/internal/modules/sys/repo"
	"go-admin/internal/modules/sys/service"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	wire.Struct(new(SYS), "*"),
	wire.Struct(new(repo.Logger), "*"),
	wire.Struct(new(service.Logger), "*"),
	wire.Struct(new(api.Logger), "*"),
)
