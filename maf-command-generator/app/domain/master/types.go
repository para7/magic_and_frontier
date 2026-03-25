package master

type DBMaster interface {
	HasItem(id string)
	HasGrimoire(id string)
	HasPassive(id string)
	HasEnemySkill(id string)
	HasEnemy(id string)
	HasSpawnTable(id string)
	HasTreasure(id string)
	HasLootTable(id string)
}
