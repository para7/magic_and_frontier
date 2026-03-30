package export_convert

import (
	"encoding/json"
	"math"
	"strconv"
)

func JsonString(value string) string {
	return string(mustJSON(value))
}

func mustJSON(value any) []byte {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return data
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func isFinite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

func ToCountValue(min, max *float64) any {
	minValue := 1.0
	maxValue := 1.0
	if min != nil {
		minValue = *min
	}
	if max != nil {
		maxValue = *max
	}
	if minValue == maxValue {
		return minValue
	}
	return map[string]any{
		"type": "minecraft:uniform",
		"min":  minValue,
		"max":  maxValue,
	}
}

func ToWeight(weight float64) int {
	if !isFinite(weight) || weight <= 0 {
		return 1
	}
	return int(math.Floor(weight))
}
