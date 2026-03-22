package idseq

import "fmt"

type CounterState struct {
	Items       int `json:"items"`
	Grimoire    int `json:"grimoire"`
	Skills      int `json:"skills"`
	EnemySkills int `json:"enemySkills"`
	Enemies     int `json:"enemies"`
	SpawnTables int `json:"spawnTables"`
	Treasures   int `json:"treasures"`
	LootTables  int `json:"loottables"`
	CastIDs     int `json:"castids"`
}

type Kind string

const (
	KindItem       Kind = "items"
	KindGrimoire   Kind = "grimoire"
	KindSkill      Kind = "skill"
	KindEnemySkill Kind = "enemyskill"
	KindEnemy      Kind = "enemy"
	KindSpawnTable Kind = "spawntable"
	KindTreasure   Kind = "treasure"
	KindLootTable  Kind = "loottable"
)

func EmptyCounterState() CounterState {
	return CounterState{}
}

func NextID(state CounterState, kind Kind) (CounterState, string) {
	switch kind {
	case KindItem:
		state.Items++
		return state, fmt.Sprintf("items_%d", state.Items)
	case KindGrimoire:
		state.Grimoire++
		return state, fmt.Sprintf("grimoire_%d", state.Grimoire)
	case KindSkill:
		state.Skills++
		return state, fmt.Sprintf("skill_%d", state.Skills)
	case KindEnemySkill:
		state.EnemySkills++
		return state, fmt.Sprintf("enemyskill_%d", state.EnemySkills)
	case KindEnemy:
		state.Enemies++
		return state, fmt.Sprintf("enemy_%d", state.Enemies)
	case KindSpawnTable:
		state.SpawnTables++
		return state, fmt.Sprintf("spawntable_%d", state.SpawnTables)
	case KindTreasure:
		state.Treasures++
		return state, fmt.Sprintf("treasure_%d", state.Treasures)
	case KindLootTable:
		state.LootTables++
		return state, fmt.Sprintf("loottable_%d", state.LootTables)
	default:
		panic("unsupported counter kind")
	}
}

func NextCastID(state CounterState) (CounterState, int) {
	state.CastIDs++
	return state, state.CastIDs
}
