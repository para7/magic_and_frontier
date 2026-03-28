package master

import (
	"log"

	model "maf_command_editor/app/domain/model"
	"maf_command_editor/app/domain/model/grimoire"
	config "maf_command_editor/app/files"
)

type DBMaster struct {
	grimoire model.MafEntity[grimoire.Grimoire]
}

func NewDBMaster(cfg config.MafConfig) *DBMaster {
	grimoire := grimoire.NewGrimoireEntity(cfg.GrimoireStatePath)
	if err := grimoire.Load(); err != nil {
		log.Fatalf("failed to load grimoire: %v", err)
	}
	return &DBMaster{grimoire: grimoire}
}

// ------ MafEntity 向けインターフェースの実装 ------

func (d *DBMaster) HasGrimoire(id string) bool {
	_, found := d.grimoire.Find(id)
	return found
}

func (d *DBMaster) HasItem(_ string) bool {
	// TODO
	return false
}
func (d *DBMaster) HasPassive(_ string) bool {
	// TODO
	return false
}
func (d *DBMaster) HasEnemySkill(_ string) bool {
	// TODO
	return false
}
func (d *DBMaster) HasEnemy(_ string) bool {
	// TODO
	return false
}
func (d *DBMaster) HasSpawnTable(_ string) bool {
	// TODO
	return false
}
func (d *DBMaster) HasTreasure(_ string) bool {
	// TODO
	return false
}
func (d *DBMaster) HasLootTable(_ string) bool {
	// TODO
	return false
}

// ------ CLI 向けユースケースの実装 ------

func (d *DBMaster) ValidateAll() [][]model.ValidationError {
	var result [][]model.ValidationError
	result = append(result, d.grimoire.ValidateAll(d)...)
	return result
}

// ------ Export 向けインターフェースの実装 ------

func (d *DBMaster) GetGrimoireByID(id string) (grimoire.Grimoire, bool) {
	return d.grimoire.Find(id)
}

func (d *DBMaster) ListGrimoires() []grimoire.Grimoire {
	entries := d.grimoire.GetAll()
	result := make([]grimoire.Grimoire, len(entries))
	copy(result, entries)
	return result
}
