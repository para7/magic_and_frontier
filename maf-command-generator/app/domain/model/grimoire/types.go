package grimoire

type Grimoire struct {
	ID          string   `json:"id"       validate:"trimmed_required,maf_slug_id"`
	CastTime    int      `json:"castTime" validate:"gte=0,lte=12000"`
	CoolTime    int      `json:"coolTime" validate:"gte=0,lte=12000"`
	MPCost      int      `json:"mpCost"   validate:"gte=0,lte=1000000"`
	Script      []string `json:"script"   validate:"min=1"`
	Title       string   `json:"title"    validate:"trimmed_required"`
	Description string   `json:"description"`
}
