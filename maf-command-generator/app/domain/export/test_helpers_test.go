package export

import (
	"testing"

	activeModel "maf_command_editor/app/domain/model/active"
	bowModel "maf_command_editor/app/domain/model/bow"
	enemyModel "maf_command_editor/app/domain/model/enemy"
	enemyskillModel "maf_command_editor/app/domain/model/enemyskill"
	itemModel "maf_command_editor/app/domain/model/item"
	passiveModel "maf_command_editor/app/domain/model/passive"
	spawntableModel "maf_command_editor/app/domain/model/spawntable"
)

type exportMasterStub struct {
	actives     []activeModel.Active
	passives    []passiveModel.Passive
	bows        []bowModel.BowPassive
	items       []itemModel.Item
	enemySkills []enemyskillModel.EnemySkill
	enemies     []enemyModel.Enemy
	spawnTables []spawntableModel.SpawnTable
}

func restoreExportUnixNow(t *testing.T, value int64) {
	t.Helper()

	original := exportUnixNow
	exportUnixNow = func() int64 {
		return value
	}
	t.Cleanup(func() {
		exportUnixNow = original
	})
}

func (s exportMasterStub) ListActives() []activeModel.Active {
	out := make([]activeModel.Active, len(s.actives))
	copy(out, s.actives)
	return out
}

func (s exportMasterStub) ListPassives() []passiveModel.Passive {
	out := make([]passiveModel.Passive, len(s.passives))
	copy(out, s.passives)
	return out
}

func (s exportMasterStub) ListBows() []bowModel.BowPassive {
	out := make([]bowModel.BowPassive, len(s.bows))
	copy(out, s.bows)
	return out
}

func (s exportMasterStub) ListItems() []itemModel.Item {
	out := make([]itemModel.Item, len(s.items))
	copy(out, s.items)
	return out
}

func (s exportMasterStub) ListEnemySkills() []enemyskillModel.EnemySkill {
	out := make([]enemyskillModel.EnemySkill, len(s.enemySkills))
	copy(out, s.enemySkills)
	return out
}

func (s exportMasterStub) ListEnemies() []enemyModel.Enemy {
	out := make([]enemyModel.Enemy, len(s.enemies))
	copy(out, s.enemies)
	return out
}

func (s exportMasterStub) ListSpawnTables() []spawntableModel.SpawnTable {
	out := make([]spawntableModel.SpawnTable, len(s.spawnTables))
	copy(out, s.spawnTables)
	return out
}
