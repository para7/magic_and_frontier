package enemy

import (
	"bytes"
	"encoding/json"
)

import model "maf_command_editor/app/domain/model"

type Enemy struct {
	ID         string         `json:"id"                  validate:"trimmed_required,maf_slug_id"`
	MobType    string         `json:"mobType"             validate:"trimmed_required"`
	Memo       string         `json:"memo"`
	Minecraft  map[string]any `json:"minecraft"`
	Maf        EnemyMaf       `json:"maf"`
	Passengers []Passenger    `json:"passengers,omitempty" validate:"omitempty,dive"`
}

type EnemyMaf struct {
	Equipment     model.Equipment `json:"equipment"`
	EnemySkillIDs []string        `json:"enemySkillIds"`
	DropMode      string          `json:"dropMode" validate:"trimmed_required,trimmed_oneof=append replace"`
	Drops         []any           `json:"drops"`
}

type Passenger struct {
	MobType    string         `json:"mobType"              validate:"trimmed_required"`
	Minecraft  map[string]any `json:"minecraft,omitempty"`
	Maf        PassengerMaf   `json:"maf,omitempty"`
	Passengers []Passenger    `json:"passengers,omitempty" validate:"omitempty,dive"`
}

type PassengerMaf struct {
	Equipment     *model.Equipment `json:"equipment,omitempty"`
	EnemySkillIDs []string         `json:"enemySkillIds,omitempty"`
	Tags          []string         `json:"tags,omitempty"`
}

func (e *Enemy) UnmarshalJSON(data []byte) error {
	type enemyAlias struct {
		ID         string          `json:"id"`
		MobType    string          `json:"mobType"`
		Memo       string          `json:"memo"`
		Minecraft  json.RawMessage `json:"minecraft"`
		Maf        EnemyMaf        `json:"maf"`
		Passengers []Passenger     `json:"passengers,omitempty"`
	}

	var decoded enemyAlias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	minecraft, err := decodeMinecraftField(decoded.Minecraft)
	if err != nil {
		return err
	}

	*e = Enemy{
		ID:         decoded.ID,
		MobType:    decoded.MobType,
		Memo:       decoded.Memo,
		Minecraft:  minecraft,
		Maf:        decoded.Maf,
		Passengers: decoded.Passengers,
	}
	return nil
}

func (p *Passenger) UnmarshalJSON(data []byte) error {
	type passengerAlias struct {
		MobType    string          `json:"mobType"`
		Minecraft  json.RawMessage `json:"minecraft,omitempty"`
		Maf        PassengerMaf    `json:"maf,omitempty"`
		Passengers []Passenger     `json:"passengers,omitempty"`
	}

	var decoded passengerAlias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	minecraft, err := decodeMinecraftField(decoded.Minecraft)
	if err != nil {
		return err
	}

	*p = Passenger{
		MobType:    decoded.MobType,
		Minecraft:  minecraft,
		Maf:        decoded.Maf,
		Passengers: decoded.Passengers,
	}
	return nil
}

func decodeMinecraftField(raw json.RawMessage) (map[string]any, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()

	var minecraft map[string]any
	if err := decoder.Decode(&minecraft); err != nil {
		return nil, err
	}
	if minecraft == nil {
		return nil, nil
	}
	return minecraft, nil
}
