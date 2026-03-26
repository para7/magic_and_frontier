package grimoire

type Grimoire struct {
	ID          string `json:"id"`
	CastID      int    `json:"castid" validate:"gte=1"`
	CastTime    int    `json:"castTime" validate:"gte=0"`
	MPCost      int    `json:"mpCost" validate:"gte=0"`
	Script      string `json:"script" validate:"trimmed_required,trimmed_min=1"`
	Title       string `json:"title" validate:"trimmed_required,trimmed_min=1"`
	Description string `json:"description"`
	UpdatedAt   string `json:"updatedAt"`
}
