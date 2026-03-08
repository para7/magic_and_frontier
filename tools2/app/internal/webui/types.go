package webui

import (
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/skills"
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
		{Path: "/treasures", Label: "Treasures"},
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
}

type ReferenceOption struct {
	ID    string
	Label string
}

type ItemFormData struct {
	ID                  string
	ItemID              string
	Count               string
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
	ID           string
	CastID       string
	Script       string
	Title        string
	Description  string
	VariantsText string
	FieldErrors  map[string]string
	FormError    string
	IsEditing    bool
}

type GrimoirePageData struct {
	Meta    PageMeta
	Notice  *Notice
	Entries []grimoire.GrimoireEntry
	Form    GrimoireFormData
}

type SkillFormData struct {
	ID          string
	Name        string
	Script      string
	ItemID      string
	ItemOptions []ReferenceOption
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
	Name        string
	Script      string
	Cooldown    string
	Trigger     string
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
	Name          string
	LootPoolsText string
	FieldErrors   map[string]string
	FormError     string
	IsEditing     bool
}

type TreasuresPageData struct {
	Meta            PageMeta
	Notice          *Notice
	Entries         []treasures.TreasureEntry
	ItemOptions     []ReferenceOption
	GrimoireOptions []ReferenceOption
	Form            TreasureFormData
}

type EnemyFormData struct {
	ID                string
	Name              string
	HP                string
	Attack            string
	Defense           string
	MoveSpeed         string
	DropTableID       string
	EnemySkillIDs     []string
	EnemySkillOptions []ReferenceOption
	OriginX           string
	OriginY           string
	OriginZ           string
	DistanceMin       string
	DistanceMax       string
	XMin              string
	XMax              string
	YMin              string
	YMax              string
	ZMin              string
	ZMax              string
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
