package master

type DBMaster interface {
	hasItem(id string)
	hasGrimoire(id string)
	hasPassive(id string)
	hasEnemySkill(id string)
	hasEnemy(id string)
	hasSpawnTable(id string)
	hasTreasure(id string)
	hasLootTable(id string)
}
