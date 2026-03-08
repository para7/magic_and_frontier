package grimoire

type Variant struct {
	Cast int `json:"cast"`
	Cost int `json:"cost"`
}

type GrimoireEntry struct {
	ID          string    `json:"id"`
	CastID      int       `json:"castid"`
	Script      string    `json:"script"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Variants    []Variant `json:"variants"`
	UpdatedAt   string    `json:"updatedAt"`
}

type GrimoireState struct {
	Entries []GrimoireEntry `json:"entries"`
}

type SaveInput struct {
	ID          string    `json:"id" validate:"required,uuid_any"`
	CastID      int       `json:"castid" validate:"gte=1"`
	Script      string    `json:"script" validate:"trimmed_required,trimmed_min=1,trimmed_max=20000"`
	Title       string    `json:"title" validate:"trimmed_required,trimmed_min=1,trimmed_max=200"`
	Description string    `json:"description"`
	Variants    []Variant `json:"variants"`
}
