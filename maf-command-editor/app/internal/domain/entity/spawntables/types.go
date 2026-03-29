package spawntables

type ReplacementEntry struct {
	EnemyID string `json:"enemyId"`
	Weight  int    `json:"weight"`
}

type SpawnTableEntry struct {
	ID            string             `json:"id"`
	SourceMobType string             `json:"sourceMobType"`
	Dimension     string             `json:"dimension"`
	MinX          int                `json:"minX"`
	MaxX          int                `json:"maxX"`
	MinY          int                `json:"minY"`
	MaxY          int                `json:"maxY"`
	MinZ          int                `json:"minZ"`
	MaxZ          int                `json:"maxZ"`
	BaseMobWeight int                `json:"baseMobWeight"`
	Replacements  []ReplacementEntry `json:"replacements"`
	UpdatedAt     string             `json:"updatedAt"`
}

type SaveInput struct {
	ID            string             `json:"id"`
	SourceMobType string             `json:"sourceMobType" validate:"trimmed_required,trimmed_min=3,trimmed_max=120"`
	Dimension     string             `json:"dimension" validate:"trimmed_required,trimmed_oneof=minecraft:overworld minecraft:the_nether minecraft:the_end"`
	MinX          int                `json:"minX" validate:"gte=-30000000,lte=30000000"`
	MaxX          int                `json:"maxX" validate:"gte=-30000000,lte=30000000"`
	MinY          int                `json:"minY" validate:"gte=-30000000,lte=30000000"`
	MaxY          int                `json:"maxY" validate:"gte=-30000000,lte=30000000"`
	MinZ          int                `json:"minZ" validate:"gte=-30000000,lte=30000000"`
	MaxZ          int                `json:"maxZ" validate:"gte=-30000000,lte=30000000"`
	BaseMobWeight int                `json:"baseMobWeight" validate:"gte=0,lte=1000000"`
	Replacements  []ReplacementEntry `json:"replacements" validate:"min=1,dive"`
}
