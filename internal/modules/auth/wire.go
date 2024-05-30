package auth

import (
	"go-admin/internal/modules/auth/api"
	"go-admin/internal/modules/auth/repo"
	"go-admin/internal/modules/auth/service"

	"github.com/google/wire"
)

// Collection of wire providers
var Set = wire.NewSet(
	wire.Struct(new(Auth), "*"),
	wire.Struct(new(Casbinx), "*"),
	wire.Struct(new(repo.Menu), "*"),
	wire.Struct(new(service.Menu), "*"),
	wire.Struct(new(api.Menu), "*"),
	wire.Struct(new(repo.MenuResource), "*"),
	wire.Struct(new(repo.Role), "*"),
	wire.Struct(new(service.Role), "*"),
	wire.Struct(new(api.Role), "*"),
	wire.Struct(new(repo.RoleMenu), "*"),
	wire.Struct(new(repo.User), "*"),
	wire.Struct(new(service.User), "*"),
	wire.Struct(new(api.User), "*"),
	wire.Struct(new(repo.UserRole), "*"),
	wire.Struct(new(service.Auth), "*"),
	wire.Struct(new(api.Auth), "*"),
)
