package skills

type SkillEntry struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Script      string `json:"script"`
	UpdatedAt   string `json:"updatedAt"`
}

type SaveInput struct {
	ID          string `json:"id"`
	Name        string `json:"name" validate:"trimmed_max=80"`
	Description string `json:"description" validate:"trimmed_max=400"`
	Script      string `json:"script" validate:"trimmed_required,trimmed_min=1,trimmed_max=20000"`
}
