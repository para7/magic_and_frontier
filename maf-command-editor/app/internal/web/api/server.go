package api

import (
	"net/http"

	"tools2/app/internal/application"
	"tools2/app/internal/config"
	"tools2/app/internal/web"
)

type Dependencies = application.Dependencies

type apiRouter struct {
	cfg     config.Config
	deps    Dependencies
	service application.Service
}

func NewHandler(cfg config.Config, deps Dependencies) http.Handler {
	service := application.NewService(cfg, deps)
	master, err := service.Master()
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writeInternalError(w, err)
		})
	}
	deps.Master = master
	app := apiRouter{
		cfg:     cfg,
		deps:    deps,
		service: service,
	}

	mux := http.NewServeMux()
	app.registerSystemRoutes(mux)
	web.RegisterRoutes(mux, cfg, deps)
	app.registerItemRoutes(mux)
	app.registerGrimoireRoutes(mux)
	app.registerSkillRoutes(mux)
	app.registerEnemySkillRoutes(mux)
	app.registerEnemyRoutes(mux)
	app.registerSpawnTableRoutes(mux)
	app.registerTreasureRoutes(mux)
	app.registerLootTableRoutes(mux)
	return mux
}
