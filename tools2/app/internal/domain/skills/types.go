package skills

type SkillEntry struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Script    string `json:"script"`
	ItemID    string `json:"itemId"`
	UpdatedAt string `json:"updatedAt"`
}

type SaveInput struct {
	ID     string `json:"id" validate:"required,uuid_any"`
	Name   string `json:"name" validate:"trimmed_required,trimmed_min=1,trimmed_max=80"`
	Script string `json:"script" validate:"trimmed_required,trimmed_min=1,trimmed_max=20000"`
	ItemID string `json:"itemId" validate:"required,uuid_any"`
}
