package master

import (
	"maf-command-editor/app/internal/domain/entity"
	"maf-command-editor/app/internal/domain/entity/enemies"
	"maf-command-editor/app/internal/domain/entity/enemyskills"
	"maf-command-editor/app/internal/domain/entity/grimoire"
	"maf-command-editor/app/internal/domain/entity/items"
	"maf-command-editor/app/internal/domain/entity/loottables"
	"maf-command-editor/app/internal/domain/entity/skills"
	"maf-command-editor/app/internal/domain/entity/spawntables"
	"maf-command-editor/app/internal/domain/entity/treasures"
)

var (
	ErrDuplicateID = entity.ErrDuplicateID
	ErrNotFound    = entity.ErrNotFound
	ErrRelation    = entity.ErrRelation
)

// 全体の統括インターフェース
type DBMaster interface {
	entity.MasterRef

	Items() entity.MafEntity[items.SaveInput, items.ItemEntry]
	Grimoires() entity.MafEntity[grimoire.SaveInput, grimoire.GrimoireEntry]
	Skills() entity.MafEntity[skills.SaveInput, skills.SkillEntry]
	EnemySkills() entity.MafEntity[enemyskills.SaveInput, enemyskills.EnemySkillEntry]
	Enemies() entity.MafEntity[enemies.SaveInput, enemies.EnemyEntry]
	Treasures() entity.MafEntity[treasures.SaveInput, treasures.TreasureEntry]
	LootTables() entity.MafEntity[loottables.SaveInput, loottables.LootTableEntry]
	SpawnTables() entity.MafEntity[spawntables.SaveInput, spawntables.SpawnTableEntry]

	ValidateSavedAll() ValidationReport
	SaveAll() error
}
