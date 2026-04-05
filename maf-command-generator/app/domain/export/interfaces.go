package export

import (
	enemyModel "maf_command_editor/app/domain/model/enemy"
	enemyskillModel "maf_command_editor/app/domain/model/enemyskill"
	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	itemModel "maf_command_editor/app/domain/model/item"
	passiveModel "maf_command_editor/app/domain/model/passive"
)

type DBMaster interface {
	ListGrimoires() []grimoireModel.Grimoire
	ListPassives() []passiveModel.Passive
	ListItems() []itemModel.Item
	ListEnemySkills() []enemyskillModel.EnemySkill
	ListEnemies() []enemyModel.Enemy
}
