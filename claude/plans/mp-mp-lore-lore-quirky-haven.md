# Context

`maf.maxmp` が設定された装備アイテムについて、その値（MaxMP ボーナス/ペナルティ）が現状 custom_data には埋め込まれているが lore には表示されていない。プレイヤーがアイテムを見たとき MP 増減量がわかるよう、lore に追加する。挿入位置は「元の lore の直後、アクティブスキルセクションの前」。

---

# 変更ファイル

`maf-command-generator/app/domain/export/convert/item.go`  
`maf-command-generator/app/domain/export/convert/item_test.go`

---

# 実装

## 1. `maxMPLoreComponents` 追加（item.go）

```go
func maxMPLoreComponents(maxMP *int) []any {
    if maxMP == nil {
        return nil
    }
    return []any{
        itemSkillLoreComponent(fmt.Sprintf("MaxMP  : %+d", *maxMP), "white", false),
    }
}
```

- ブランク行なし（元 lore と連続した stat 行として扱う）
- `%+d` で `+20` / `-20` を表示
- スキルセクションのブランク行はアクティブ側が持っているので不要

## 2. `mergedItemSkillLore` のシグネチャ変更（item.go）

```go
// before
func mergedItemSkillLore(existing any, spellMeta itemSpellMeta) ([]any, bool)

// after
func mergedItemSkillLore(existing any, spellMeta itemSpellMeta, maxMP *int) ([]any, bool)
```

ボディ変更：
```go
func mergedItemSkillLore(existing any, spellMeta itemSpellMeta, maxMP *int) ([]any, bool) {
    mpLore := maxMPLoreComponents(maxMP)
    skillLore := itemSkillLoreComponents(spellMeta)
    if len(mpLore) == 0 && len(skillLore) == 0 {
        return nil, false
    }
    lore := anyList(existing)
    lore = append(lore, mpLore...)
    lore = append(lore, skillLore...)
    return lore, true
}
```

## 3. 呼び出し元 2 箇所の更新（item.go）

| 場所 | 変更前 | 変更後 |
|------|--------|--------|
| `itemComponentsForGive` (L95) | `mergedItemSkillLore(..., spellMeta)` | `mergedItemSkillLore(..., spellMeta, entry.Maf.MaxMP)` |
| `itemComponentsForLoot` (L179) | `mergedItemSkillLore(..., spellMeta)` | `mergedItemSkillLore(..., spellMeta, entry.Maf.MaxMP)` |

## 4. テスト追加（item_test.go）

`TestItemComponentsForGiveIncludesMaxMPInLore`：
- MaxMP あり（例: `+20`）のアイテムで `itemComponentsForGive` を呼ぶ
- `minecraft:lore` に `MaxMP  : +20` が含まれることを確認
- スキルなし・MaxMP のみの場合も lore が返ることを確認

`TestMergedItemSkillLoreMaxMPBeforeActiveSkill`（任意）：
- MaxMP あり + activeId ありのアイテムで、`MaxMP  : +20` が `アクティブスキル` より前に来ることを確認

---

# 出力イメージ

元 lore が `["剣士の兜"]` で `maxmp: 15` の場合：

```
剣士の兜
MaxMP  : +15
                      ← active の blank line
アクティブスキル
炎斬り
...
```

MaxMP のみ（スキルなし）：

```
アーマーの説明
MaxMP  : -10
```

---

# 検証

```
cd maf-command-generator && make check
```

- `TestItemComponentsForGiveIncludesMaxMPInLore` が PASS
- 既存のスキル lore テスト（active/passive/bow）が影響を受けないこと
- `make run/export` でアーマー類の give コマンドに `MaxMP  :` 行が出力されること
