package master

import (
	"log"

	export "maf_command_editor/app/domain/export"
	model "maf_command_editor/app/domain/model"
	"maf_command_editor/app/domain/model/grimoire"
	config "maf_command_editor/app/files"
)

type DBMasterImpl struct {
	grimoire model.MafEntity[grimoire.Grimoire]
}

var _ export.DBMaster = (*DBMasterImpl)(nil)

func NewDBMaster(cfg config.MafConfig) *DBMasterImpl {
	grimoire := grimoire.NewGrimoireEntity(cfg.GrimoireStatePath)
	if err := grimoire.Load(); err != nil {
		log.Fatalf("failed to load grimoire: %v", err)
	}
	return &DBMasterImpl{grimoire: grimoire}
}

// ------ MafEntity 向けインターフェースの実装 ------

func (d *DBMasterImpl) HasGrimoire(id string) bool {
	_, found := d.grimoire.Find(id)
	return found
}

func (d *DBMasterImpl) HasItem(_ string) bool {
	// TODO
	return false
}
func (d *DBMasterImpl) HasPassive(_ string) bool {
	// TODO
	return false
}
func (d *DBMasterImpl) HasEnemySkill(_ string) bool {
	// TODO
	return false
}
func (d *DBMasterImpl) HasEnemy(_ string) bool {
	// TODO
	return false
}
func (d *DBMasterImpl) HasSpawnTable(_ string) bool {
	// TODO
	return false
}
func (d *DBMasterImpl) HasTreasure(_ string) bool {
	// TODO
	return false
}
func (d *DBMasterImpl) HasLootTable(_ string) bool {
	// TODO
	return false
}

// ------ CLI 向けユースケースの実装 ------

func (d *DBMasterImpl) ValidateAll() [][]model.ValidationError {
	var result [][]model.ValidationError
	result = append(result, d.grimoire.ValidateAll(d)...)
	return result
}

// ------ Export 向けインターフェースの実装 ------

func (d *DBMasterImpl) GetGrimoireByID(id string) (grimoire.Grimoire, bool) {
	return d.grimoire.Find(id)
}

func (d *DBMasterImpl) ListGrimoires() []grimoire.Grimoire {
	entries := d.grimoire.GetAll()
	result := make([]grimoire.Grimoire, len(entries))
	copy(result, entries)
	return result
}
