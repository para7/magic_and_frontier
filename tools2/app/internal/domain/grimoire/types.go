package grimoire

type GrimoireEntry struct {
	ID          string `json:"id"`
	CastID      int    `json:"castid"`
	CastTime    int    `json:"castTime"`
	MPCost      int    `json:"mpCost"`
	Script      string `json:"script"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	UpdatedAt   string `json:"updatedAt"`
}

type GrimoireState struct {
	Entries []GrimoireEntry `json:"entries"`
}

type SaveInput struct {
	ID          string `json:"id"`
	CastID      int    `json:"castid" validate:"gte=1"`
	CastTime    int    `json:"castTime" validate:"gte=0,lte=12000"`
	MPCost      int    `json:"mpCost" validate:"gte=0,lte=1000000"`
	Script      string `json:"script" validate:"trimmed_required,trimmed_min=1,trimmed_max=20000"`
	Title       string `json:"title" validate:"trimmed_required,trimmed_min=1,trimmed_max=200"`
	Description string `json:"description"`
}
