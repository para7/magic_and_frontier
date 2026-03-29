package passive

type Passive struct {
	ID          string `json:"id"          validate:"trimmed_required"`
	Name        string `json:"name"        validate:"trimmed_max=80"`
	SkillType   string `json:"skilltype"   validate:"trimmed_oneof=sword bow axe"`
	Description string `json:"description" validate:"trimmed_max=400"`
	Script      string `json:"script"      validate:"trimmed_required"`
	UpdatedAt   string `json:"updatedAt"`
}
