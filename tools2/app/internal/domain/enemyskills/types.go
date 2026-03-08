package enemyskills

type Trigger string

const (
	TriggerOnSpawn Trigger = "on_spawn"
	TriggerOnHit   Trigger = "on_hit"
	TriggerOnLowHP Trigger = "on_low_hp"
	TriggerOnTimer Trigger = "on_timer"
)

type EnemySkillEntry struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Script    string   `json:"script"`
	Cooldown  *float64 `json:"cooldown,omitempty"`
	Trigger   *Trigger `json:"trigger,omitempty"`
	UpdatedAt string   `json:"updatedAt"`
}

type SaveInput struct {
	ID       string   `json:"id" validate:"required,uuid_any"`
	Name     string   `json:"name" validate:"trimmed_required,trimmed_min=1,trimmed_max=80"`
	Script   string   `json:"script" validate:"trimmed_required,trimmed_min=1,trimmed_max=20000"`
	Cooldown *float64 `json:"cooldown,omitempty" validate:"omitempty,gte=0,lte=12000"`
	Trigger  string   `json:"trigger,omitempty" validate:"omitempty,trimmed_oneof=on_spawn on_hit on_low_hp on_timer"`
}
