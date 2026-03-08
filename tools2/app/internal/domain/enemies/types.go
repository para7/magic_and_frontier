package enemies

type DropRef struct {
	Kind     string   `json:"kind" validate:"trimmed_required,trimmed_oneof=item grimoire"`
	RefID    string   `json:"refId" validate:"required,uuid_any"`
	Weight   float64  `json:"weight" validate:"gte=1,lte=100000"`
	CountMin *float64 `json:"countMin,omitempty" validate:"omitempty,gte=1,lte=64"`
	CountMax *float64 `json:"countMax,omitempty" validate:"omitempty,gte=1,lte=64"`
}

type AxisBounds struct {
	XMin *float64 `json:"xMin,omitempty" validate:"omitempty,gte=-30000000,lte=30000000"`
	XMax *float64 `json:"xMax,omitempty" validate:"omitempty,gte=-30000000,lte=30000000"`
	YMin *float64 `json:"yMin,omitempty" validate:"omitempty,gte=-30000000,lte=30000000"`
	YMax *float64 `json:"yMax,omitempty" validate:"omitempty,gte=-30000000,lte=30000000"`
	ZMin *float64 `json:"zMin,omitempty" validate:"omitempty,gte=-30000000,lte=30000000"`
	ZMax *float64 `json:"zMax,omitempty" validate:"omitempty,gte=-30000000,lte=30000000"`
}

type Vec3 struct {
	X float64 `json:"x" validate:"gte=-30000000,lte=30000000"`
	Y float64 `json:"y" validate:"gte=-30000000,lte=30000000"`
	Z float64 `json:"z" validate:"gte=-30000000,lte=30000000"`
}

type Distance struct {
	Min float64 `json:"min" validate:"gte=0,lte=30000000"`
	Max float64 `json:"max" validate:"gte=0,lte=30000000"`
}

type SpawnRule struct {
	Origin     Vec3        `json:"origin" validate:"required"`
	Distance   Distance    `json:"distance" validate:"required"`
	AxisBounds *AxisBounds `json:"axisBounds,omitempty" validate:"omitempty"`
}

type EnemyEntry struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	HP            float64   `json:"hp"`
	Attack        *float64  `json:"attack,omitempty"`
	Defense       *float64  `json:"defense,omitempty"`
	MoveSpeed     *float64  `json:"moveSpeed,omitempty"`
	DropTableID   string    `json:"dropTableId"`
	EnemySkillIDs []string  `json:"enemySkillIds"`
	SpawnRule     SpawnRule `json:"spawnRule"`
	DropTable     []DropRef `json:"dropTable,omitempty"`
	UpdatedAt     string    `json:"updatedAt"`
}

type SaveInput struct {
	ID            string    `json:"id" validate:"required,uuid_any"`
	Name          string    `json:"name" validate:"trimmed_required,trimmed_min=1,trimmed_max=80"`
	HP            float64   `json:"hp" validate:"gte=1,lte=100000"`
	Attack        *float64  `json:"attack,omitempty" validate:"omitempty,gte=0,lte=100000"`
	Defense       *float64  `json:"defense,omitempty" validate:"omitempty,gte=0,lte=100000"`
	MoveSpeed     *float64  `json:"moveSpeed,omitempty" validate:"omitempty,gte=0,lte=100000"`
	DropTableID   string    `json:"dropTableId" validate:"trimmed_required,trimmed_min=1,trimmed_max=200"`
	EnemySkillIDs []string  `json:"enemySkillIds" validate:"dive,required,uuid_any"`
	SpawnRule     SpawnRule `json:"spawnRule" validate:"required"`
	DropTable     []DropRef `json:"dropTable,omitempty" validate:"dive"`
}
