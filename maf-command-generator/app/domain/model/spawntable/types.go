package spawntable

import model "maf_command_editor/app/domain/model"

type SpawnTable struct {
	ID            string                   `json:"id"            validate:"trimmed_required"`
	SourceMobType string                   `json:"sourceMobType" validate:"trimmed_required"`
	Dimension     string                   `json:"dimension"     validate:"trimmed_required,trimmed_oneof=minecraft:overworld minecraft:the_nether minecraft:the_end"`
	MinX          int                      `json:"minX"          validate:"gte=-30000000,lte=30000000"`
	MaxX          int                      `json:"maxX"          validate:"gte=-30000000,lte=30000000"`
	MinY          int                      `json:"minY"          validate:"gte=-30000000,lte=30000000"`
	MaxY          int                      `json:"maxY"          validate:"gte=-30000000,lte=30000000"`
	MinZ          int                      `json:"minZ"          validate:"gte=-30000000,lte=30000000"`
	MaxZ          int                      `json:"maxZ"          validate:"gte=-30000000,lte=30000000"`
	BaseMobWeight int                      `json:"baseMobWeight" validate:"gte=0,lte=1000000"`
	Replacements  []model.ReplacementEntry `json:"replacements"  validate:"min=1"`
	UpdatedAt     string                   `json:"updatedAt"`
}
