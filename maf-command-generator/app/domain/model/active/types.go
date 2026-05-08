package active

import "encoding/json"

type Active struct {
	ID          string   `json:"id"            validate:"trimmed_required,maf_slug_id"`
	CastTime    int      `json:"castTime"      validate:"gte=0,lte=12000"`
	CoolTime    int      `json:"coolTime"      validate:"gte=0,lte=12000"`
	MPCost      int      `json:"mpCost"        validate:"gte=0,lte=1000000"`
	Script      []string `json:"script"        validate:"min=1"`
	Title       string   `json:"title"         validate:"trimmed_required"`
	Description string   `json:"description"`
	LootEnable  *bool    `json:"loot_enable"`
}

func (g *Active) UnmarshalJSON(data []byte) error {
	type alias Active

	var v alias
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v.LootEnable == nil {
		defaultTrue := true
		v.LootEnable = &defaultTrue
	}

	*g = Active(v)
	return nil
}

func (g Active) IsLootEnabled() bool {
	return g.LootEnable == nil || *g.LootEnable
}
