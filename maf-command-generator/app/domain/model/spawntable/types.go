package spawntable

import model "maf_command_editor/app/domain/model"

type SpawnTable struct {
	ID            string                   `json:"id"            validate:"trimmed_required"`
	SourceMobType string                   `json:"sourceMobType" validate:"trimmed_required"`
	Dimension     string                   `json:"dimension"     validate:"trimmed_required,trimmed_oneof=minecraft:overworld minecraft:the_nether minecraft:the_end"`
	MinDistance   int                      `json:"minDistance"   validate:"gte=0,lte=99999999"`
	MaxDistance   int                      `json:"maxDistance"   validate:"gte=0,lte=99999999"`
	MinX          int                      `json:"minX"          validate:"gte=-99999999,lte=99999999"`
	MaxX          int                      `json:"maxX"          validate:"gte=-99999999,lte=99999999"`
	MinY          int                      `json:"minY"          validate:"gte=-99999999,lte=99999999"`
	MaxY          int                      `json:"maxY"          validate:"gte=-99999999,lte=99999999"`
	MinZ          int                      `json:"minZ"          validate:"gte=-99999999,lte=99999999"`
	MaxZ          int                      `json:"maxZ"          validate:"gte=-99999999,lte=99999999"`
	BaseMobWeight int                      `json:"baseMobWeight" validate:"gte=0,lte=1000000"`
	Replacements  []model.ReplacementEntry `json:"replacements"  validate:"min=1"`
}
