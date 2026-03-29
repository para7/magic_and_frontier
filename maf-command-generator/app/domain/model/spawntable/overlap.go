package spawntable

// RangesOverlap は2つのスポーンテーブルの座標範囲が重なるかどうかを判定する。
func RangesOverlap(left, right SpawnTable) bool {
	return intervalOverlap(left.MinX, left.MaxX, right.MinX, right.MaxX) &&
		intervalOverlap(left.MinY, left.MaxY, right.MinY, right.MaxY) &&
		intervalOverlap(left.MinZ, left.MaxZ, right.MinZ, right.MaxZ)
}

// AllOverlaps は同じ dimension・sourceMobType の間で重なるペアを返す。
func AllOverlaps(entries []SpawnTable) [][2]string {
	var pairs [][2]string
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			left := entries[i]
			right := entries[j]
			if left.SourceMobType != right.SourceMobType || left.Dimension != right.Dimension {
				continue
			}
			if RangesOverlap(left, right) {
				pairs = append(pairs, [2]string{left.ID, right.ID})
			}
		}
	}
	return pairs
}

func intervalOverlap(leftMin, leftMax, rightMin, rightMax int) bool {
	return leftMin <= rightMax && rightMin <= leftMax
}
