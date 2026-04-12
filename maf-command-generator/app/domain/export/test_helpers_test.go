package export

import (
	bowModel "maf_command_editor/app/domain/model/bow"
	enemyModel "maf_command_editor/app/domain/model/enemy"
	enemyskillModel "maf_command_editor/app/domain/model/enemyskill"
	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	itemModel "maf_command_editor/app/domain/model/item"
	passiveModel "maf_command_editor/app/domain/model/passive"
)

type exportMasterStub struct {
	grimoires   []grimoireModel.Grimoire
	passives    []passiveModel.Passive
	bows        []bowModel.BowPassive
	items       []itemModel.Item
	enemySkills []enemyskillModel.EnemySkill
	enemies     []enemyModel.Enemy
}

func (s exportMasterStub) ListGrimoires() []grimoireModel.Grimoire {
	out := make([]grimoireModel.Grimoire, len(s.grimoires))
	copy(out, s.grimoires)
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

func ptrFloat(v float64) *float64 {
	return &v
}

func ptrInt(v int) *int {
	return &v
}
