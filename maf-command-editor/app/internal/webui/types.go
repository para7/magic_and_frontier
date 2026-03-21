package webui

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

type NavItem struct {
	Path  string
	Label string
}

func NavItems() []NavItem {
	return []NavItem{
		{Path: "/items", Label: "Items"},
		{Path: "/grimoire", Label: "Grimoire"},
		{Path: "/skills", Label: "Skills"},
		{Path: "/enemy-skills", Label: "Enemy Skills"},
		{Path: "/spawn-tables", Label: "Spawn Tables"},
		{Path: "/treasures", Label: "Treasures"},
		{Path: "/loottables", Label: "Loottables"},
		{Path: "/enemies", Label: "Enemies"},
	}
}

type Notice struct {
	Kind string
	Text string
}

type PageMeta struct {
	Title       string
	CurrentPath string
	Description string
	CurrentURL  string
}

type ReferenceOption struct {
	ID    string
	Label string
}

type SelectOption struct {
	Value string
	Label string
}

type ItemFormData struct {
	ID                  string
	ReturnTo            string
	ItemID              string
	Count               string
	SkillID             string
	SkillOptions        []ReferenceOption
	CustomName          string
	Lore                string
	Enchantments        string
	Unbreakable         bool
	CustomModelData     string
	RepairCost          string
	HideFlags           string
	PotionID            string
	CustomPotionColor   string
	CustomPotionEffects string
	AttributeModifiers  string
	CustomNBT           string
	FieldErrors         map[string]string
	FormError           string
	IsEditing           bool
}

type ItemsPageData struct {
	Meta    PageMeta
	Notice  *Notice
	Entries []items.ItemEntry
	Form    ItemFormData
}

type GrimoireFormData struct {
	ID          string
	ReturnTo    string
	CastID      string
	CastTime    string
	MPCost      string
	Script      string
	Title       string
	Description string
	FieldErrors map[string]string
	FormError   string
	IsEditing   bool
}

type GrimoirePageData struct {
	Meta    PageMeta
	Notice  *Notice
	Entries []grimoire.GrimoireEntry
	Form    GrimoireFormData
}

type SkillFormData struct {
	ID          string
	ReturnTo    string
	Name        string
	Description string
	Script      string
	FieldErrors map[string]string
	FormError   string
	IsEditing   bool
}

type SkillsPageData struct {
	Meta    PageMeta
	Notice  *Notice
	Entries []skills.SkillEntry
	Form    SkillFormData
}

type EnemySkillFormData struct {
	ID          string
	ReturnTo    string
	Name        string
	Description string
	Script      string
	FieldErrors map[string]string
	FormError   string
	IsEditing   bool
}

type EnemySkillsPageData struct {
	Meta    PageMeta
	Notice  *Notice
	Entries []enemyskills.EnemySkillEntry
	Form    EnemySkillFormData
}

type TreasureFormData struct {
	ID            string
	TablePath     string
	ReturnTo      string
	LootPoolsText string
	FieldErrors   map[string]string
	FormError     string
	IsEditing     bool
	HasSource     bool
	HasOverlay    bool
}

type TreasureListEntry struct {
	ID         string
	TablePath  string
	LootPools  []treasures.DropRef
	UpdatedAt  string
	HasSource  bool
	HasOverlay bool
}

type TreasuresPageData struct {
	Meta            PageMeta
	Notice          *Notice
	Entries         []TreasureListEntry
	ItemOptions     []ReferenceOption
	GrimoireOptions []ReferenceOption
	Form            TreasureFormData
}

type LootTableFormData struct {
	ID            string
	ReturnTo      string
	LootPoolsText string
	FieldErrors   map[string]string
	FormError     string
	IsEditing     bool
}

type LootTablesPageData struct {
	Meta            PageMeta
	Notice          *Notice
	Entries         []loottables.LootTableEntry
	ItemOptions     []ReferenceOption
	GrimoireOptions []ReferenceOption
	Form            LootTableFormData
}

type EnemyFormData struct {
	ID                string
	ReturnTo          string
	MobType           string
	Name              string
	HP                string
	Attack            string
	Defense           string
	MoveSpeed         string
	DropMode          string
	EnemySkillIDs     []string
	EnemySkillOptions []ReferenceOption
	EquipmentText     string
	DropsText         string
	FieldErrors       map[string]string
	FormError         string
	IsEditing         bool
}

type EnemiesPageData struct {
	Meta    PageMeta
	Notice  *Notice
	Entries []enemies.EnemyEntry
	Form    EnemyFormData
}

type SpawnTableFormData struct {
	ID               string
	ReturnTo         string
	SourceMobType    string
	Dimension        string
	DimensionOptions []SelectOption
	MinX             string
	MaxX             string
	MinY             string
	MaxY             string
	MinZ             string
	MaxZ             string
	BaseMobWeight    string
	ReplacementsText string
	FieldErrors      map[string]string
	FormError        string
	IsEditing        bool
}

type SpawnTablesPageData struct {
	Meta    PageMeta
	Notice  *Notice
	Entries []spawntables.SpawnTableEntry
	Form    SpawnTableFormData
}
