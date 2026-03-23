package httpapi

import (
	"net/http"
	"time"

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
	deps = normalizeDependencies(cfg, deps)
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

func normalizeDependencies(cfg config.Config, deps Dependencies) Dependencies {
	defaults := DefaultDependencies(cfg)
	if deps.ItemRepo == nil {
		deps.ItemRepo = defaults.ItemRepo
	}
	if deps.GrimoireRepo == nil {
		deps.GrimoireRepo = defaults.GrimoireRepo
	}
	if deps.SkillRepo == nil {
		deps.SkillRepo = defaults.SkillRepo
	}
	if deps.EnemySkillRepo == nil {
		deps.EnemySkillRepo = defaults.EnemySkillRepo
	}
	if deps.EnemyRepo == nil {
		deps.EnemyRepo = defaults.EnemyRepo
	}
	if deps.SpawnTableRepo == nil {
		deps.SpawnTableRepo = defaults.SpawnTableRepo
	}
	if deps.TreasureRepo == nil {
		deps.TreasureRepo = defaults.TreasureRepo
	}
	if deps.LootTableRepo == nil {
		deps.LootTableRepo = defaults.LootTableRepo
	}
	if deps.CounterRepo == nil {
		deps.CounterRepo = defaults.CounterRepo
	}
	if deps.ExportSettingsPath == "" {
		deps.ExportSettingsPath = defaults.ExportSettingsPath
	}
	if deps.Now == nil {
		deps.Now = time.Now
	}
	return deps
}

func DefaultDependencies(cfg config.Config) Dependencies {
	return application.DefaultDependencies(cfg)
}
