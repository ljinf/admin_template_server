//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"poem_server_admin/internal/cache"
	"poem_server_admin/internal/handler"
	hSys "poem_server_admin/internal/handler/sys"
	"poem_server_admin/internal/repository"
	repSys "poem_server_admin/internal/repository/sys"
	"poem_server_admin/internal/server"
	"poem_server_admin/internal/service"
	srvSys "poem_server_admin/internal/service/sys"
	"poem_server_admin/pkg/app"
	"poem_server_admin/pkg/jwt"
	"poem_server_admin/pkg/log"
	"poem_server_admin/pkg/server/http"
	"poem_server_admin/pkg/sid"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repSys.NewDictRepository,
	repSys.NewUserRepository,
	repSys.NewMenuRepository,
	repSys.NewRoleRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	srvSys.NewUserService,
	srvSys.NewMenuService,
	srvSys.NewRoleService,
	srvSys.NewDictService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	hSys.NewAccountHandler,
	hSys.NewDictHandler,
	hSys.NewMenuHandler,
	hSys.NewRoleHandler,
	hSys.NewUserHandler,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
)

var cacheSet = wire.NewSet(
	cache.NewAccountCache,
)

// build App
func newApp(
	httpServer *http.Server,
	// job *server.Job,
	// task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer),
		app.WithName("poem-admin-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		cacheSet,
		serviceSet,
		handlerSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}
