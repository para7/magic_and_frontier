package enemy

import model "maf_command_editor/app/domain/model"

type Enemy struct {
	ID            string          `json:"id"                    validate:"trimmed_required"`
	MobType       string          `json:"mobType"               validate:"trimmed_required"`
	Name          string          `json:"name"`
	HP            float64         `json:"hp"                    validate:"gte=1,lte=100000"`
	Memo          string          `json:"memo"`
	Attack        *float64        `json:"attack,omitempty"      validate:"omitempty,gte=0,lte=100000"`
	Defense       *float64        `json:"defense,omitempty"     validate:"omitempty,gte=0,lte=100000"`
	MoveSpeed     *float64        `json:"moveSpeed,omitempty"   validate:"omitempty,gte=0,lte=100000"`
	Equipment     model.Equipment `json:"equipment"`
	EnemySkillIDs []string        `json:"enemySkillIds"`
	DropMode      string          `json:"dropMode"              validate:"trimmed_required,trimmed_oneof=append replace"`
	Drops         []model.DropRef `json:"drops"                 validate:"dive"`
}
