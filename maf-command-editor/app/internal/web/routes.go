package web

import (
	"net/http"
	"time"

	"tools2/app/internal/application"
	"tools2/app/internal/config"
)

type Dependencies = application.Dependencies

type App struct {
	cfg  config.Config
	deps Dependencies
}

func RegisterRoutes(mux *http.ServeMux, cfg config.Config, deps Dependencies) {
	defaults := application.DefaultDependencies(cfg)
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
	if deps.Now == nil {
		deps.Now = time.Now
	}
	app := App{cfg: cfg, deps: deps}

	mux.HandleFunc("GET /items", app.itemsPage)
	mux.HandleFunc("GET /items/new", app.itemsNewPage)
	mux.HandleFunc("POST /items/new", app.itemsSubmit)
	mux.HandleFunc("GET /items/edit", app.itemsEditPage)
	mux.HandleFunc("POST /items/edit", app.itemsEditSubmit)
	mux.HandleFunc("POST /items/{id}/delete", app.itemsDelete)

	mux.HandleFunc("GET /grimoire", app.grimoirePage)
	mux.HandleFunc("GET /grimoire/new", app.grimoireNewPage)
	mux.HandleFunc("POST /grimoire/new", app.grimoireSubmit)
	mux.HandleFunc("GET /grimoire/edit", app.grimoireEditPage)
	mux.HandleFunc("POST /grimoire/edit", app.grimoireEditSubmit)
	mux.HandleFunc("POST /grimoire/{id}/delete", app.grimoireDelete)

	mux.HandleFunc("GET /skills", app.skillsPage)
	mux.HandleFunc("GET /skills/new", app.skillsNewPage)
	mux.HandleFunc("POST /skills/new", app.skillsSubmit)
	mux.HandleFunc("GET /skills/edit", app.skillsEditPage)
	mux.HandleFunc("POST /skills/edit", app.skillsEditSubmit)
	mux.HandleFunc("POST /skills/{id}/delete", app.skillsDelete)

	mux.HandleFunc("GET /enemy-skills", app.enemySkillsPage)
	mux.HandleFunc("GET /enemy-skills/new", app.enemySkillsNewPage)
	mux.HandleFunc("POST /enemy-skills/new", app.enemySkillsSubmit)
	mux.HandleFunc("GET /enemy-skills/edit", app.enemySkillsEditPage)
	mux.HandleFunc("POST /enemy-skills/edit", app.enemySkillsEditSubmit)
	mux.HandleFunc("POST /enemy-skills/{id}/delete", app.enemySkillsDelete)

	mux.HandleFunc("GET /treasures", app.treasuresPage)
	mux.HandleFunc("GET /treasures/new", app.treasuresNewPage)
	mux.HandleFunc("POST /treasures/new", app.treasuresSubmit)
	mux.HandleFunc("GET /treasures/edit", app.treasuresEditPage)
	mux.HandleFunc("POST /treasures/edit", app.treasuresEditSubmit)
	mux.HandleFunc("POST /treasures/{id}/delete", app.treasuresDelete)

	mux.HandleFunc("GET /loottables", app.lootTablesPage)
	mux.HandleFunc("GET /loottables/new", app.lootTablesNewPage)
	mux.HandleFunc("POST /loottables/new", app.lootTablesSubmit)
	mux.HandleFunc("GET /loottables/edit", app.lootTablesEditPage)
	mux.HandleFunc("POST /loottables/edit", app.lootTablesEditSubmit)
	mux.HandleFunc("POST /loottables/{id}/delete", app.lootTablesDelete)

	mux.HandleFunc("GET /enemies", app.enemiesPage)
	mux.HandleFunc("GET /enemies/new", app.enemiesNewPage)
	mux.HandleFunc("POST /enemies/new", app.enemiesSubmit)
	mux.HandleFunc("GET /enemies/edit", app.enemiesEditPage)
	mux.HandleFunc("POST /enemies/edit", app.enemiesEditSubmit)
	mux.HandleFunc("POST /enemies/{id}/delete", app.enemiesDelete)

	mux.HandleFunc("GET /spawn-tables", app.spawnTablesPage)
	mux.HandleFunc("GET /spawn-tables/new", app.spawnTablesNewPage)
	mux.HandleFunc("POST /spawn-tables/new", app.spawnTablesSubmit)
	mux.HandleFunc("GET /spawn-tables/edit", app.spawnTablesEditPage)
	mux.HandleFunc("POST /spawn-tables/edit", app.spawnTablesEditSubmit)
	mux.HandleFunc("POST /spawn-tables/{id}/delete", app.spawnTablesDelete)

	mux.HandleFunc("POST /save", app.saveExport)
}
