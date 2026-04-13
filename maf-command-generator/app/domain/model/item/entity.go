package item

import (
	"errors"
	"fmt"
	"strings"

	cv "maf_command_editor/app/domain/custom_validator"
	model "maf_command_editor/app/domain/model"
	"maf_command_editor/app/files"
)

type ItemEntity struct {
	store files.JsonStore[Item]
	data  []Item
}

func NewItemEntity(path string) *ItemEntity {
	return &ItemEntity{store: files.NewJsonStore[Item](path)}
}

func (s *ItemEntity) ValidateJSON(newEntity Item, mas model.DBMaster) (Item, []model.ValidationError) {
	var errs []model.ValidationError
	errs = append(errs, s.ValidateStruct(newEntity)...)
	errs = append(errs, s.ValidateRelation(newEntity, mas)...)

	// item components 文字列を再生成して component 構造の破綻を検出する
	_, buildErr := BuildItemComponents(newEntity)
	if buildErr != "" {
		errs = append(errs, model.ValidationError{
			Entity: "item", ID: newEntity.ID,
			Field: "minecraft.components",
			Tag:   "format", Param: buildErr,
		})
	}

	if len(errs) > 0 {
		return Item{}, errs
	}
	return newEntity, nil
}

func (s *ItemEntity) ValidateStruct(newEntity Item) []model.ValidationError {
	var errs []model.ValidationError

	err := cv.Validate.Struct(newEntity)
	if err != nil {
		for _, fe := range err.(cv.ValidationErrors) {
			errs = append(errs, cv.NewValidationError("item", newEntity.ID, fe))
		}
	}

	for key, value := range newEntity.Minecraft.Components {
		normalizedKey := strings.TrimSpace(key)
		if normalizedKey == "" {
			errs = append(errs, model.ValidationError{
				Entity: "item", ID: newEntity.ID,
				Field: "minecraft.components",
				Tag:   "format", Param: "component key is empty",
			})
			continue
		}
		if !strings.Contains(normalizedKey, ":") {
			errs = append(errs, model.ValidationError{
				Entity: "item", ID: newEntity.ID,
				Field: "minecraft.components",
				Tag:   "format", Param: fmt.Sprintf("component key must be namespaced: %q", normalizedKey),
			})
		}
		if strings.TrimSpace(value) == "" {
			errs = append(errs, model.ValidationError{
				Entity: "item", ID: newEntity.ID,
				Field: "minecraft.components",
				Tag:   "format", Param: fmt.Sprintf("component value is empty: %q", normalizedKey),
			})
		}
	}
	return errs
}

func (s *ItemEntity) ValidateRelation(newEntity Item, mas model.DBMaster) []model.ValidationError {
	var errs []model.ValidationError
	grimoireID := strings.TrimSpace(newEntity.Maf.GrimoireID)
	if grimoireID != "" && !mas.HasGrimoire(grimoireID) {
		errs = append(errs, model.ValidationError{
			Entity: "item", ID: newEntity.ID,
			Field: "maf.grimoireId",
			Tag:   "relation", Param: "grimoire not found",
		})
	}

	passiveID := strings.TrimSpace(newEntity.Maf.PassiveID)
	if passiveID != "" && !mas.HasPassive(passiveID) {
		errs = append(errs, model.ValidationError{
			Entity: "item", ID: newEntity.ID,
			Field: "maf.passiveId",
			Tag:   "relation", Param: "passive not found",
		})
	}

	bowID := strings.TrimSpace(newEntity.Maf.BowID)
	if bowID == "" {
		return errs
	}

	itemID := strings.TrimSpace(newEntity.Minecraft.ItemID)
	if itemID != "minecraft:bow" && itemID != "minecraft:crossbow" {
		errs = append(errs, model.ValidationError{
			Entity: "item", ID: newEntity.ID,
			Field: "minecraft.itemId",
			Tag:   "relation", Param: "bowId requires minecraft:bow or minecraft:crossbow",
		})
	}
	if passiveID != "" {
		errs = append(errs, model.ValidationError{
			Entity: "item", ID: newEntity.ID,
			Field: "maf.passiveId",
			Tag:   "relation", Param: "bowId cannot be combined with passiveId",
		})
	}
	if grimoireID != "" {
		errs = append(errs, model.ValidationError{
			Entity: "item", ID: newEntity.ID,
			Field: "maf.grimoireId",
			Tag:   "relation", Param: "bowId cannot be combined with grimoireId",
		})
	}
	if !mas.HasBow(bowID) {
		errs = append(errs, model.ValidationError{
			Entity: "item", ID: newEntity.ID,
			Field: "maf.bowId",
			Tag:   "relation", Param: "bow not found",
		})
	}
	return errs
}

func (s *ItemEntity) Create(newEntity Item, mas model.DBMaster) error {
	validated, errs := s.ValidateJSON(newEntity, mas)
	if len(errs) > 0 {
		return fmt.Errorf("%s.%s: %s", errs[0].Entity, errs[0].Field, errs[0].Tag)
	}
	for _, it := range s.data {
		if it.ID == validated.ID {
			return errors.New("item id already exists: " + validated.ID)
		}
	}
	s.data = append(s.data, validated)
	return nil
}

func (s *ItemEntity) Update(newEntity Item, mas model.DBMaster) error {
	validated, errs := s.ValidateJSON(newEntity, mas)
	if len(errs) > 0 {
		return fmt.Errorf("%s.%s: %s", errs[0].Entity, errs[0].Field, errs[0].Tag)
	}
	for i, it := range s.data {
		if it.ID == validated.ID {
			s.data[i] = validated
			return nil
		}
	}
	return errors.New("item not found: " + validated.ID)
}

func (s *ItemEntity) Delete(id string, mas model.DBMaster) error {
	for i, it := range s.data {
		if it.ID == id {
			s.data = append(s.data[:i], s.data[i+1:]...)
			return nil
		}
	}
	return errors.New("item not found: " + id)
}

func (s *ItemEntity) Save() error {
	return s.store.Save(s.data)
}

func (s *ItemEntity) Load() error {
	data, err := s.store.Load()
	if err != nil {
		return err
	}
	s.data = data
	fmt.Printf("[item.Load] Loaded %d records\n", len(data))
	return nil
}

func (s *ItemEntity) ValidateAll(mas model.DBMaster) [][]model.ValidationError {
	var result [][]model.ValidationError
	seenIDs := map[string]bool{}
	for _, it := range s.data {
		if _, errs := s.ValidateJSON(it, mas); len(errs) > 0 {
			result = append(result, errs)
		}
		if seenIDs[it.ID] {
			result = append(result, []model.ValidationError{{
				Entity: "item",
				ID:     it.ID,
				Field:  "id",
				Tag:    "unique",
				Param:  "ID重複を検出",
			}})
			continue
		}
		seenIDs[it.ID] = true
	}
	if len(result) > 0 {
		fmt.Printf("[item.ValidateAll] Found errors in %d record(s)\n", len(result))
	} else {
		fmt.Printf("[item.ValidateAll] No errors found\n")
	}
	return result
}

func (s *ItemEntity) Find(id string) (Item, bool) {
	for _, it := range s.data {
		if it.ID == id {
			return it, true
		}
	}
	return Item{}, false
}

func (s *ItemEntity) GetAll() []Item {
	return s.data
}
