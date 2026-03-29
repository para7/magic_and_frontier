package grimoire

type Grimoire struct {
	ID          string `json:"id"       validate:"trimmed_required"`
	CastID      int    `json:"castid"   validate:"gte=1"`
	CastTime    int    `json:"castTime" validate:"gte=0,lte=12000"`
	MPCost      int    `json:"mpCost"   validate:"gte=0,lte=1000000"`
	Script      string `json:"script"   validate:"trimmed_required"`
	Title       string `json:"title"    validate:"trimmed_required"`
	Description string `json:"description"`
}
