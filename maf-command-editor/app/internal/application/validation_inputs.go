package application

import (
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/loottables"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/spawntables"
	"tools2/app/internal/domain/treasures"
)

func itemToInput(entry items.ItemEntry) items.SaveInput {
	return items.SaveInput{
		ID:                  entry.ID,
		ItemID:              entry.ItemID,
		SkillID:             entry.SkillID,
		CustomName:          entry.CustomName,
		Lore:                entry.Lore,
		Enchantments:        entry.Enchantments,
		Unbreakable:         entry.Unbreakable,
		CustomModelData:     entry.CustomModelData,
		RepairCost:          entry.RepairCost,
		HideFlags:           entry.HideFlags,
		PotionID:            entry.PotionID,
		CustomPotionColor:   entry.CustomPotionColor,
		CustomPotionEffects: entry.CustomPotionEffects,
		AttributeModifiers:  entry.AttributeModifiers,
		CustomNBT:           entry.CustomNBT,
	}
}

func grimoireToInput(entry grimoire.GrimoireEntry) grimoire.SaveInput {
	return grimoire.SaveInput{
		ID:          entry.ID,
		CastID:      entry.CastID,
		CastTime:    entry.CastTime,
		MPCost:      entry.MPCost,
		Script:      entry.Script,
		Title:       entry.Title,
		Description: entry.Description,
	}
}

func skillToInput(entry skills.SkillEntry) skills.SaveInput {
	return skills.SaveInput{
		ID:          entry.ID,
		Name:        entry.Name,
		Description: entry.Description,
		Script:      entry.Script,
	}
}

func enemySkillToInput(entry enemyskills.EnemySkillEntry) enemyskills.SaveInput {
	return enemyskills.SaveInput{
		ID:          entry.ID,
		Name:        entry.Name,
		Description: entry.Description,
		Script:      entry.Script,
	}
}

func treasureToInput(entry treasures.TreasureEntry) treasures.SaveInput {
	return treasures.SaveInput{
		ID:        entry.ID,
		TablePath: entry.TablePath,
		LootPools: append([]treasures.DropRef{}, entry.LootPools...),
	}
}

func loottableToInput(entry loottables.LootTableEntry) loottables.SaveInput {
	return loottables.SaveInput{
		ID:        entry.ID,
		LootPools: append([]treasures.DropRef{}, entry.LootPools...),
	}
}

func enemyToInput(entry enemies.EnemyEntry) enemies.SaveInput {
	return enemies.SaveInput{
		ID:            entry.ID,
		MobType:       entry.MobType,
		Name:          entry.Name,
		HP:            entry.HP,
		Attack:        entry.Attack,
		Defense:       entry.Defense,
		MoveSpeed:     entry.MoveSpeed,
		Equipment:     entry.Equipment,
		EnemySkillIDs: append([]string{}, entry.EnemySkillIDs...),
		DropMode:      entry.DropMode,
		Drops:         append([]enemies.DropRef{}, entry.Drops...),
	}
}

func spawnTableToInput(entry spawntables.SpawnTableEntry) spawntables.SaveInput {
	return spawntables.SaveInput{
		ID:            entry.ID,
		SourceMobType: entry.SourceMobType,
		Dimension:     entry.Dimension,
		MinX:          entry.MinX,
		MaxX:          entry.MaxX,
		MinY:          entry.MinY,
		MaxY:          entry.MaxY,
		MinZ:          entry.MinZ,
		MaxZ:          entry.MaxZ,
		BaseMobWeight: entry.BaseMobWeight,
		Replacements:  append([]spawntables.ReplacementEntry{}, entry.Replacements...),
	}
}
