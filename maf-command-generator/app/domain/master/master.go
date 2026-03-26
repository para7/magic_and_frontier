package master

import (
	model "maf_command_editor/app/domain/model"
	"maf_command_editor/app/domain/model/grimoire"
)

type DBMasterImpl struct {
	grimoire model.MafEntity[grimoire.Grimoire]
}

func NewDBMaster(grimoire model.MafEntity[grimoire.Grimoire]) *DBMasterImpl {
	return &DBMasterImpl{grimoire: grimoire}
}

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
