package model

type ActiveSnapshot struct {
	ID               string
	GenerateGrimoire bool
}

type PassiveSnapshot struct {
	ID               string
	GenerateGrimoire *bool
}

type DBMaster interface {
	HasItem(id string) bool
	HasActive(id string) bool
	GetActive(id string) (ActiveSnapshot, bool)
	HasPassive(id string) bool
	GetPassive(id string) (PassiveSnapshot, bool)
	HasBow(id string) bool
	HasEnemySkill(id string) bool
	HasEnemy(id string) bool
	HasSpawnTable(id string) bool
	HasMinecraftLootTable(tablePath string) bool
}

// Model の共通インターフェース
type MafEntity[T any] interface {
	ValidateJSON(newEntity T, mas DBMaster) (T, []ValidationError)
	Load() error

	// DBMaster との連携用
	ValidateAll(mas DBMaster) [][]ValidationError
	Find(id string) (T, bool)
	GetAll() []T
}
