package sys

import (
	"go-admin/internal/mods/sys/api"
	"go-admin/internal/mods/sys/biz"
	"go-admin/internal/mods/sys/dal"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	wire.Struct(new(SYS), "*"),
	wire.Struct(new(dal.Logger), "*"),
	wire.Struct(new(biz.Logger), "*"),
	wire.Struct(new(api.Logger), "*"),
)
