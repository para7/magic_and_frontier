package model

type DBMaster interface {
	HasItem(id string) bool
	HasGrimoire(id string) bool
	HasPassive(id string) bool
	HasEnemySkill(id string) bool
	HasEnemy(id string) bool
	HasSpawnTable(id string) bool
	HasTreasure(id string) bool
	HasLootTable(id string) bool
	HasMinecraftLootTable(tablePath string) bool
}

// Model の共通インターフェース
type MafEntity[T any] interface {
	ValidateJSON(newEntity T, mas DBMaster) (T, []ValidationError)

	// 主に web 画面からの操作用、メモリで保持してる配列にデータを追記する。リレーション関係を確認するため validate も行う
	Create(newEntity T, mas DBMaster) error
	Update(newEntity T, mas DBMaster) error
	Delete(id string, mas DBMaster) error
	Save() error
	Load() error

	// DBMaster との連携用
	ValidateAll(mas DBMaster) [][]ValidationError
	Find(id string) (T, bool)
	GetAll() []T
}
