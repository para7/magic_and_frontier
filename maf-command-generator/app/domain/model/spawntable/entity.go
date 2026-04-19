package spawntable

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cv "maf_command_editor/app/domain/custom_validator"
	model "maf_command_editor/app/domain/model"
	"maf_command_editor/app/files"
)

type SpawnTableEntity struct {
	store files.JsonStore[SpawnTable]
	data  []SpawnTable
}

type spawnTableFile struct {
	Coordinates *spawnTableCoordinates `json:"coordinates"`
	Entries     []SpawnTable           `json:"entries"`
}

type spawnTableCoordinates struct {
	Dimension   string `json:"dimension"`
	MinDistance int    `json:"minDistance"`
	MaxDistance int    `json:"maxDistance"`
	MinX        int    `json:"minX"`
	MaxX        int    `json:"maxX"`
	MinY        int    `json:"minY"`
	MaxY        int    `json:"maxY"`
	MinZ        int    `json:"minZ"`
	MaxZ        int    `json:"maxZ"`
}

func NewSpawnTableEntity(path string) *SpawnTableEntity {
	return &SpawnTableEntity{store: files.NewJsonStore[SpawnTable](path)}
}

func (s *SpawnTableEntity) ValidateJSON(newEntity SpawnTable, mas model.DBMaster) (SpawnTable, []model.ValidationError) {
	var errs []model.ValidationError
	errs = append(errs, s.ValidateStruct(newEntity)...)
	errs = append(errs, s.ValidateRelation(newEntity, mas)...)
	if len(errs) > 0 {
		return SpawnTable{}, errs
	}
	return newEntity, nil
}

func (s *SpawnTableEntity) ValidateStruct(newEntity SpawnTable) []model.ValidationError {
	err := cv.Validate.Struct(newEntity)
	if err == nil {
		return nil
	}
	var errs []model.ValidationError
	for _, fe := range err.(cv.ValidationErrors) {
		errs = append(errs, cv.NewValidationError("spawntable", newEntity.ID, fe))
	}
	return errs
}

func (s *SpawnTableEntity) ValidateRelation(newEntity SpawnTable, mas model.DBMaster) []model.ValidationError {
	var errs []model.ValidationError

	// SourceMobType は namespaced resource ID
	if src := strings.TrimSpace(newEntity.SourceMobType); src != "" && !model.IsNamespacedResourceID(src) {
		errs = append(errs, model.ValidationError{
			Entity: "spawntable", ID: newEntity.ID,
			Field: "sourceMobType",
			Tag:   "format", Param: "invalid minecraft entity id",
		})
	}

	// 座標の Min <= Max チェック
	if newEntity.MinX > newEntity.MaxX {
		errs = append(errs, model.ValidationError{Entity: "spawntable", ID: newEntity.ID, Field: "minX", Tag: "lte", Param: "maxX"})
	}
	if newEntity.MinDistance > newEntity.MaxDistance {
		errs = append(errs, model.ValidationError{Entity: "spawntable", ID: newEntity.ID, Field: "minDistance", Tag: "lte", Param: "maxDistance"})
	}
	if newEntity.MinY > newEntity.MaxY {
		errs = append(errs, model.ValidationError{Entity: "spawntable", ID: newEntity.ID, Field: "minY", Tag: "lte", Param: "maxY"})
	}
	if newEntity.MinZ > newEntity.MaxZ {
		errs = append(errs, model.ValidationError{Entity: "spawntable", ID: newEntity.ID, Field: "minZ", Tag: "lte", Param: "maxZ"})
	}

	// Replacements の参照チェック
	seen := map[string]bool{}
	totalWeight := newEntity.BaseMobWeight
	for i, r := range newEntity.Replacements {
		enemyID := strings.TrimSpace(r.EnemyID)
		if !mas.HasEnemy(enemyID) {
			errs = append(errs, model.ValidationError{
				Entity: "spawntable", ID: newEntity.ID,
				Field: fmt.Sprintf("replacements[%d].enemyId", i),
				Tag:   "relation", Param: "enemy not found",
			})
			continue
		}
		if seen[enemyID] {
			errs = append(errs, model.ValidationError{
				Entity: "spawntable", ID: newEntity.ID,
				Field: fmt.Sprintf("replacements[%d].enemyId", i),
				Tag:   "relation", Param: "duplicate enemy id",
			})
			continue
		}
		if r.Weight <= 0 {
			errs = append(errs, model.ValidationError{
				Entity: "spawntable", ID: newEntity.ID,
				Field: fmt.Sprintf("replacements[%d].weight", i),
				Tag:   "gt", Param: "0",
			})
			continue
		}
		seen[enemyID] = true
		totalWeight += r.Weight
	}
	if totalWeight <= 0 {
		errs = append(errs, model.ValidationError{
			Entity: "spawntable", ID: newEntity.ID,
			Field: "baseMobWeight",
			Tag:   "relation", Param: "total weight must be > 0",
		})
	}

	return errs
}

func (s *SpawnTableEntity) Load() error {
	paths, err := filepath.Glob(filepath.Join(s.store.Path, "*.json"))
	if err != nil {
		return err
	}
	loaded := make([]SpawnTable, 0)
	for _, path := range paths {
		entries, err := loadSpawnTableFile(path)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", filepath.Base(path), err)
		}
		loaded = append(loaded, entries...)
	}
	s.data = loaded
	fmt.Printf("[spawntable.Load] Loaded %d records\n", len(loaded))
	return nil
}

func (s *SpawnTableEntity) ValidateAll(mas model.DBMaster) [][]model.ValidationError {
	var result [][]model.ValidationError
	seenIDs := map[string]bool{}
	for _, st := range s.data {
		if _, errs := s.ValidateJSON(st, mas); len(errs) > 0 {
			result = append(result, errs)
		}
		if seenIDs[st.ID] {
			result = append(result, []model.ValidationError{{
				Entity: "spawntable",
				ID:     st.ID,
				Field:  "id",
				Tag:    "unique",
				Param:  "ID重複を検出",
			}})
			continue
		}
		seenIDs[st.ID] = true
	}
	if len(result) > 0 {
		fmt.Printf("[spawntable.ValidateAll] Found errors in %d record(s)\n", len(result))
	} else {
		fmt.Printf("[spawntable.ValidateAll] No errors found\n")
	}
	return result
}

func (s *SpawnTableEntity) Find(id string) (SpawnTable, bool) {
	for _, st := range s.data {
		if st.ID == id {
			return st, true
		}
	}
	return SpawnTable{}, false
}

func (s *SpawnTableEntity) GetAll() []SpawnTable {
	return s.data
}

func loadSpawnTableFile(path string) ([]SpawnTable, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var file spawnTableFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, err
	}
	if len(file.Entries) == 0 {
		return []SpawnTable{}, nil
	}

	baseName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	out := make([]SpawnTable, 0, len(file.Entries))
	for i, entry := range file.Entries {
		table := entry
		if file.Coordinates != nil {
			table.Dimension = file.Coordinates.Dimension
			table.MinDistance = file.Coordinates.MinDistance
			table.MaxDistance = file.Coordinates.MaxDistance
			table.MinX = file.Coordinates.MinX
			table.MaxX = file.Coordinates.MaxX
			table.MinY = file.Coordinates.MinY
			table.MaxY = file.Coordinates.MaxY
			table.MinZ = file.Coordinates.MinZ
			table.MaxZ = file.Coordinates.MaxZ
		}
		if strings.TrimSpace(table.ID) == "" {
			table.ID = buildSpawnTableID(baseName, i, table.SourceMobType)
		}
		out = append(out, table)
	}
	return out, nil
}

func buildSpawnTableID(fileBase string, index int, sourceMobType string) string {
	filePart := normalizeSpawnTableIDPart(fileBase, "spawn")
	mobPart := normalizeSpawnTableIDPart(sourceMobType, "mob")
	return fmt.Sprintf("%s_%s_%d", filePart, mobPart, index+1)
}

func normalizeSpawnTableIDPart(raw, fallback string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if strings.Contains(raw, ":") {
		parts := strings.SplitN(raw, ":", 2)
		raw = parts[1]
	}
	var builder strings.Builder
	prevUnderscore := false
	for _, r := range raw {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_':
			builder.WriteRune(r)
			prevUnderscore = false
		default:
			if !prevUnderscore {
				builder.WriteByte('_')
				prevUnderscore = true
			}
		}
	}
	normalized := strings.Trim(builder.String(), "_")
	if normalized == "" {
		return fallback
	}
	return normalized
}
