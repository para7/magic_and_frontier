package passive

type Passive struct {
	ID               string   `json:"id"          validate:"trimmed_required,maf_slug_id"`
	Name             string   `json:"name"        validate:"trimmed_max=80"`
	Lore             string   `json:"lore"        validate:"trimmed_max=200"`
	Condition        string   `json:"condition"   validate:"trimmed_required,trimmed_oneof=always attack none"`
	Slots            []int    `json:"slots"       validate:"min=1,unique,dive,gte=1,lte=3"`
	Description      string   `json:"description"       validate:"trimmed_max=400"`
	Target           string   `json:"target"`
	Range            int      `json:"range"             validate:"gte=0"`
	Script           []string `json:"script"            validate:"min=1"`
	GenerateGrimoire *bool    `json:"generate_grimoire" validate:"required"`
}
