package item

import (
	"bytes"
	"encoding/json"
)

type Item struct {
	ID        string         `json:"id"                validate:"trimmed_required"`
	ItemID    string         `json:"itemId"            validate:"trimmed_required"`
	Maf       ItemMaf        `json:"maf,omitempty"`
	Minecraft map[string]any `json:"minecraft,omitempty"`
}

func (i *Item) UnmarshalJSON(data []byte) error {
	type itemAlias struct {
		ID        string          `json:"id"`
		ItemID    string          `json:"itemId"`
		Maf       ItemMaf         `json:"maf,omitempty"`
		Minecraft json.RawMessage `json:"minecraft,omitempty"`
	}

	var decoded itemAlias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	minecraft, err := decodeMinecraftField(decoded.Minecraft)
	if err != nil {
		return err
	}

	*i = Item{
		ID:        decoded.ID,
		ItemID:    decoded.ItemID,
		Maf:       decoded.Maf,
		Minecraft: minecraft,
	}
	return nil
}

type ItemMaf struct {
	GrimoireID  string `json:"grimoireId,omitempty"`
	PassiveID   string `json:"passiveId,omitempty"`
	PassiveSlot int    `json:"passiveSlot,omitempty" validate:"omitempty,gte=1,lte=3"`
	BowID       string `json:"bowId,omitempty"`
	MaxMP       *int   `json:"maxmp,omitempty"`
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
