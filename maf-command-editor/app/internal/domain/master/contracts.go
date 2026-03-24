package master

import (
	"errors"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/loottables"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/spawntables"
	"tools2/app/internal/domain/treasures"
)

var (
	ErrDuplicateID = errors.New("duplicate id")
	ErrNotFound    = errors.New("entry not found")
	ErrRelation    = errors.New("relation validation failed")
)

// 各JSONに対応するエンティティ
type MafEntity[I any, E any] interface {
	Validate(input I, master DBMaster) common.SaveResult[E]
	Create(entry E, master DBMaster) error
	Update(entry E, master DBMaster) error
	Delete(id string, master DBMaster) error
	Save() error
	ListAll() []E
	FindByID(id string) (E, bool)
	HasID(id string) bool
}

// 全体の統括インターフェース
type DBMaster interface {
	HasItem(id string) bool
	HasGrimoire(id string) bool
	HasSkill(id string) bool
	HasEnemySkill(id string) bool
	HasEnemy(id string) bool
	HasTreasure(id string) bool
	HasLootTable(id string) bool
	HasSpawnTable(id string) bool

	Items() MafEntity[items.SaveInput, items.ItemEntry]
	Grimoires() MafEntity[grimoire.SaveInput, grimoire.GrimoireEntry]
	Skills() MafEntity[skills.SaveInput, skills.SkillEntry]
	EnemySkills() MafEntity[enemyskills.SaveInput, enemyskills.EnemySkillEntry]
	Enemies() MafEntity[enemies.SaveInput, enemies.EnemyEntry]
	Treasures() MafEntity[treasures.SaveInput, treasures.TreasureEntry]
	LootTables() MafEntity[loottables.SaveInput, loottables.LootTableEntry]
	SpawnTables() MafEntity[spawntables.SaveInput, spawntables.SpawnTableEntry]

	ValidateSavedAll() ValidationReport
	SaveAll() error
}
