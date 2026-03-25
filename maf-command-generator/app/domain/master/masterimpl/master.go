package masterimpl

import (
	"maf_command_editor/app/domain/model/grimoire"
	model "maf_command_editor/app/domain/model/maf_entity"
)

type DBMasterImpl struct {
	grimoire model.MafEntity[grimoire.Grimoire]
}

func NewDBMaster(grimoire model.MafEntity[grimoire.Grimoire]) *DBMasterImpl {
	return &DBMasterImpl{grimoire: grimoire}
}

func (d *DBMasterImpl) HasGrimoire(id string) bool {
	return d.grimoire.Find(id) != (grimoire.Grimoire{})
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
