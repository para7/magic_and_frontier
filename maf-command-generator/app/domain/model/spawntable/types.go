package spawntable

import (
	"encoding/json"

	model "maf_command_editor/app/domain/model"
)

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
	BaseMob       *BaseMob                 `json:"baseMob,omitempty"`
	Replacements  []model.ReplacementEntry `json:"replacements"`
}

type BaseMob struct {
	Weight     int                `json:"weight"     validate:"gte=0,lte=1000000"`
	Attributes *BaseMobAttributes `json:"attributes,omitempty"`
}

type BaseMobAttributes struct {
	HP        *float64 `json:"hp,omitempty"        validate:"omitempty,gte=1,lte=100000"`
	Attack    *float64 `json:"attack,omitempty"    validate:"omitempty,gte=0,lte=100000"`
	Defense   *float64 `json:"defense,omitempty"   validate:"omitempty,gte=0,lte=100000"`
	MoveSpeed *float64 `json:"moveSpeed,omitempty" validate:"omitempty,gte=0,lte=100000"`
}

func (s SpawnTable) GetBaseMobWeight() int {
	if s.BaseMob == nil {
		return 0
	}
	return s.BaseMob.Weight
}

func (s SpawnTable) GetBaseMobAttributes() *BaseMobAttributes {
	if s.BaseMob == nil {
		return nil
	}
	return s.BaseMob.Attributes
}

type spawnTableJSON struct {
	ID            string                   `json:"id"`
	SourceMobType string                   `json:"sourceMobType"`
	Dimension     string                   `json:"dimension"`
	MinDistance   int                      `json:"minDistance"`
	MaxDistance   int                      `json:"maxDistance"`
	MinX          int                      `json:"minX"`
	MaxX          int                      `json:"maxX"`
	MinY          int                      `json:"minY"`
	MaxY          int                      `json:"maxY"`
	MinZ          int                      `json:"minZ"`
	MaxZ          int                      `json:"maxZ"`
	BaseMob       *BaseMob                 `json:"baseMob"`
	BaseMobWeight *int                     `json:"baseMobWeight"`
	Replacements  []model.ReplacementEntry `json:"replacements"`
}

func (s *SpawnTable) UnmarshalJSON(data []byte) error {
	var decoded spawnTableJSON
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	baseMob := decoded.BaseMob
	if baseMob == nil && decoded.BaseMobWeight != nil {
		baseMob = &BaseMob{Weight: *decoded.BaseMobWeight}
	}

	*s = SpawnTable{
		ID:            decoded.ID,
		SourceMobType: decoded.SourceMobType,
		Dimension:     decoded.Dimension,
		MinDistance:   decoded.MinDistance,
		MaxDistance:   decoded.MaxDistance,
		MinX:          decoded.MinX,
		MaxX:          decoded.MaxX,
		MinY:          decoded.MinY,
		MaxY:          decoded.MaxY,
		MinZ:          decoded.MinZ,
		MaxZ:          decoded.MaxZ,
		BaseMob:       baseMob,
		Replacements:  decoded.Replacements,
	}
	return nil
}
