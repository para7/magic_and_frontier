package web

import (
	"fmt"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/loottables"
	dmaster "tools2/app/internal/domain/master"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/spawntables"
	"tools2/app/internal/domain/treasures"
)

func (a App) masterOrErr() (dmaster.DBMaster, error) {
	if a.deps.Master == nil {
		if a.masterInitErr != nil {
			return nil, a.masterInitErr
		}
		return nil, fmt.Errorf("master is not initialized")
	}
	return a.deps.Master, nil
}

func (a App) loadItemStateFromMaster() (common.EntryState[items.ItemEntry], error) {
	master, err := a.masterOrErr()
	if err != nil {
		return common.EntryState[items.ItemEntry]{}, err
	}
	return common.EntryState[items.ItemEntry]{Entries: master.Items().ListAll()}, nil
}

func (a App) loadGrimoireStateFromMaster() (common.EntryState[grimoire.GrimoireEntry], error) {
	master, err := a.masterOrErr()
	if err != nil {
		return common.EntryState[grimoire.GrimoireEntry]{}, err
	}
	return common.EntryState[grimoire.GrimoireEntry]{Entries: master.Grimoires().ListAll()}, nil
}

func (a App) loadSkillStateFromMaster() (common.EntryState[skills.SkillEntry], error) {
	master, err := a.masterOrErr()
	if err != nil {
		return common.EntryState[skills.SkillEntry]{}, err
	}
	return common.EntryState[skills.SkillEntry]{Entries: master.Skills().ListAll()}, nil
}

func (a App) loadEnemySkillStateFromMaster() (common.EntryState[enemyskills.EnemySkillEntry], error) {
	master, err := a.masterOrErr()
	if err != nil {
		return common.EntryState[enemyskills.EnemySkillEntry]{}, err
	}
	return common.EntryState[enemyskills.EnemySkillEntry]{Entries: master.EnemySkills().ListAll()}, nil
}

func (a App) loadEnemyStateFromMaster() (common.EntryState[enemies.EnemyEntry], error) {
	master, err := a.masterOrErr()
	if err != nil {
		return common.EntryState[enemies.EnemyEntry]{}, err
	}
	return common.EntryState[enemies.EnemyEntry]{Entries: master.Enemies().ListAll()}, nil
}

func (a App) loadSpawnTableStateFromMaster() (common.EntryState[spawntables.SpawnTableEntry], error) {
	master, err := a.masterOrErr()
	if err != nil {
		return common.EntryState[spawntables.SpawnTableEntry]{}, err
	}
	return common.EntryState[spawntables.SpawnTableEntry]{Entries: master.SpawnTables().ListAll()}, nil
}

func (a App) loadTreasureStateFromMaster() (common.EntryState[treasures.TreasureEntry], error) {
	master, err := a.masterOrErr()
	if err != nil {
		return common.EntryState[treasures.TreasureEntry]{}, err
	}
	return common.EntryState[treasures.TreasureEntry]{Entries: master.Treasures().ListAll()}, nil
}

func (a App) loadLootTableStateFromMaster() (common.EntryState[loottables.LootTableEntry], error) {
	master, err := a.masterOrErr()
	if err != nil {
		return common.EntryState[loottables.LootTableEntry]{}, err
	}
	return common.EntryState[loottables.LootTableEntry]{Entries: master.LootTables().ListAll()}, nil
}
