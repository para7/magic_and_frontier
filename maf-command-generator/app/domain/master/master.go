package master

import (
	"fmt"
	"log"

	model "maf_command_editor/app/domain/model"
	"maf_command_editor/app/domain/model/bow"
	"maf_command_editor/app/domain/model/enemy"
	"maf_command_editor/app/domain/model/enemyskill"
	"maf_command_editor/app/domain/model/grimoire"
	"maf_command_editor/app/domain/model/item"
	"maf_command_editor/app/domain/model/loottable"
	"maf_command_editor/app/domain/model/passive"
	"maf_command_editor/app/domain/model/spawntable"
	"maf_command_editor/app/domain/model/treasure"
	config "maf_command_editor/app/files"
	mc "maf_command_editor/app/minecraft"
)

type DBMaster struct {
	grimoire   model.MafEntity[grimoire.Grimoire]
	item       model.MafEntity[item.Item]
	passive    model.MafEntity[passive.Passive]
	bow        model.MafEntity[bow.BowPassive]
	enemyskill model.MafEntity[enemyskill.EnemySkill]
	enemy      model.MafEntity[enemy.Enemy]
	spawntable model.MafEntity[spawntable.SpawnTable]
	treasure   model.MafEntity[treasure.Treasure]
	loottable  model.MafEntity[loottable.LootTable]

	minecraftLootTableRoot string
}

func NewDBMaster(cfg config.MafConfig) *DBMaster {
	d := &DBMaster{
		minecraftLootTableRoot: cfg.MinecraftLootTableRoot,
	}

	load := func(name string, loader func() error) {
		if err := loader(); err != nil {
			log.Fatalf("failed to load %s: %v", name, err)
		}
	}

	d.grimoire = grimoire.NewGrimoireEntity(cfg.GrimoireStatePath)
	load("grimoire", d.grimoire.Load)

	d.item = item.NewItemEntity(cfg.ItemStatePath)
	load("item", d.item.Load)

	d.passive = passive.NewPassiveEntity(cfg.PassiveStatePath)
	load("passive", d.passive.Load)

	d.bow = bow.NewBowEntity(cfg.BowStatePath)
	load("bow", d.bow.Load)

	d.enemyskill = enemyskill.NewEnemySkillEntity(cfg.EnemySkillStatePath)
	load("enemyskill", d.enemyskill.Load)

	d.enemy = enemy.NewEnemyEntity(cfg.EnemyStatePath)
	load("enemy", d.enemy.Load)

	d.spawntable = spawntable.NewSpawnTableEntity(cfg.SpawnTableStatePath)
	load("spawntable", d.spawntable.Load)

	d.treasure = treasure.NewTreasureEntity(cfg.TreasureStatePath)
	load("treasure", d.treasure.Load)

	d.loottable = loottable.NewLootTableEntity(cfg.LootTablesStatePath)
	load("loottable", d.loottable.Load)

	return d
}

// ------ MafEntity 向けインターフェースの実装 ------

func (d *DBMaster) HasGrimoire(id string) bool {
	_, found := d.grimoire.Find(id)
	return found
}

func (d *DBMaster) HasItem(id string) bool {
	_, found := d.item.Find(id)
	return found
}

func (d *DBMaster) HasPassive(id string) bool {
	_, found := d.passive.Find(id)
	return found
}

func (d *DBMaster) HasBow(id string) bool {
	_, found := d.bow.Find(id)
	return found
}

func (d *DBMaster) HasEnemySkill(id string) bool {
	_, found := d.enemyskill.Find(id)
	return found
}

func (d *DBMaster) HasEnemy(id string) bool {
	_, found := d.enemy.Find(id)
	return found
}

func (d *DBMaster) HasSpawnTable(id string) bool {
	_, found := d.spawntable.Find(id)
	return found
}

func (d *DBMaster) HasTreasure(id string) bool {
	_, found := d.treasure.Find(id)
	return found
}

func (d *DBMaster) HasLootTable(id string) bool {
	_, found := d.loottable.Find(id)
	return found
}

func (d *DBMaster) HasMinecraftLootTable(tablePath string) bool {
	exists, err := mc.Exists(d.minecraftLootTableRoot, tablePath)
	if err != nil {
		return false
	}
	return exists
}

// ------ CLI 向けユースケースの実装 ------

func (d *DBMaster) ValidateAll() [][]model.ValidationError {
	var result [][]model.ValidationError
	result = append(result, d.grimoire.ValidateAll(d)...)
	result = append(result, d.item.ValidateAll(d)...)
	result = append(result, d.passive.ValidateAll(d)...)
	result = append(result, d.bow.ValidateAll(d)...)
	result = append(result, d.enemyskill.ValidateAll(d)...)
	result = append(result, d.enemy.ValidateAll(d)...)
	result = append(result, d.spawntable.ValidateAll(d)...)
	result = append(result, d.treasure.ValidateAll(d)...)
	result = append(result, d.loottable.ValidateAll(d)...)

	// SpawnTable の重複チェック（全体にまたがる検証）
	tables := d.spawntable.GetAll()
	for _, pair := range spawntable.AllOverlaps(tables) {
		result = append(result, []model.ValidationError{{
			Entity: "spawntable",
			ID:     pair[0],
			Field:  "dimension/range",
			Tag:    "overlap",
			Param:  fmt.Sprintf("overlaps with %s", pair[1]),
		}})
	}

	return result
}

// ------ Export 向けインターフェースの実装 ------

func (d *DBMaster) GetGrimoireByID(id string) (grimoire.Grimoire, bool) {
	return d.grimoire.Find(id)
}

func (d *DBMaster) ListGrimoires() []grimoire.Grimoire {
	entries := d.grimoire.GetAll()
	result := make([]grimoire.Grimoire, len(entries))
	copy(result, entries)
	return result
}

func (d *DBMaster) ListItems() []item.Item {
	entries := d.item.GetAll()
	result := make([]item.Item, len(entries))
	copy(result, entries)
	return result
}

func (d *DBMaster) GetItemByID(id string) (item.Item, bool) {
	return d.item.Find(id)
}

func (d *DBMaster) ListPassives() []passive.Passive {
	entries := d.passive.GetAll()
	result := make([]passive.Passive, len(entries))
	copy(result, entries)
	return result
}

func (d *DBMaster) ListBows() []bow.BowPassive {
	entries := d.bow.GetAll()
	result := make([]bow.BowPassive, len(entries))
	copy(result, entries)
	return result
}

func (d *DBMaster) ListEnemySkills() []enemyskill.EnemySkill {
	entries := d.enemyskill.GetAll()
	result := make([]enemyskill.EnemySkill, len(entries))
	copy(result, entries)
	return result
}

func (d *DBMaster) ListEnemies() []enemy.Enemy {
	entries := d.enemy.GetAll()
	result := make([]enemy.Enemy, len(entries))
	copy(result, entries)
	return result
}

func (d *DBMaster) ListSpawnTables() []spawntable.SpawnTable {
	entries := d.spawntable.GetAll()
	result := make([]spawntable.SpawnTable, len(entries))
	copy(result, entries)
	return result
}

func (d *DBMaster) ListTreasures() []treasure.Treasure {
	entries := d.treasure.GetAll()
	result := make([]treasure.Treasure, len(entries))
	copy(result, entries)
	return result
}

func (d *DBMaster) ListLootTables() []loottable.LootTable {
	entries := d.loottable.GetAll()
	result := make([]loottable.LootTable, len(entries))
	copy(result, entries)
	return result
}

func (d *DBMaster) MinecraftLootTableRoot() string {
	return d.minecraftLootTableRoot
}
