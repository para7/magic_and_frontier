package master

import (
	"tools2/app/internal/domain/entity"
	"tools2/app/internal/domain/entity/enemies"
	"tools2/app/internal/domain/entity/enemyskills"
	"tools2/app/internal/domain/entity/grimoire"
	"tools2/app/internal/domain/entity/items"
	"tools2/app/internal/domain/entity/loottables"
	"tools2/app/internal/domain/entity/skills"
	"tools2/app/internal/domain/entity/spawntables"
	"tools2/app/internal/domain/entity/treasures"
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
