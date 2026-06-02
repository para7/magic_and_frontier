# ItemToGiveCommand / ItemToUpdateCommand 重複解消

## Context

`convert/item.go` の `ItemToGiveCommand`（L14-46）と `ItemToUpdateCommand`（L48-87）は、コンポーネント解決ロジック（spellMeta解決 → componentValues構築 → consumable付与 → custom_data付与 → ソート → itemID取得）が完全同一。末尾のコマンド文字列組み立てだけが異なる。今後どちらか一方の解決ロジックを変更した際、もう片方への反映漏れが起きやすい。共通部分を private ヘルパーに抽出し、各関数は末尾組み立てのみ担う構造にする。

## 重要な観察

`ItemToGiveCommand` の空コンポーネント分岐（`give @p {id} 1`）と非空分岐（`give @p {id}[...] 1`）は、itemStack 文字列（空なら `minecraft:stone`、非空なら `minecraft:stone[...]`）に統一すれば `give @p {itemStack} 1` 一本に畳める。`ItemToUpdateCommand` の itemStack 構築（L76-79）と同じ形。よって共通ヘルパーは **itemStack 文字列を返す** 形が最適。

## 変更

### `maf-command-generator/app/domain/export/convert/item.go`

private ヘルパーを新設（既存の `itemComponentsForGive` / `sortedItemGiveComponents` / `itemCustomData` を再利用）:

```go
func buildItemStack(
    entry itemModel.Item,
    activesByID map[string]activeModel.Active,
    passivesByID map[string]passiveModel.Passive,
    bowsByID map[string]bowModel.BowPassive,
    ver int64,
) (string, error) {
    spellMeta, err := resolveItemSpellMeta(entry, activesByID, passivesByID, bowsByID)
    if err != nil {
        return "", err
    }
    componentValues, err := itemComponentsForGive(entry, spellMeta)
    if err != nil {
        return "", err
    }
    if spellMeta.hasUseSpell {
        componentValues["minecraft:consumable"] = bookConsumableSNBT
    }
    customData, err := itemCustomData(entry, activesByID, passivesByID, bowsByID, ver)
    if err != nil {
        return "", err
    }
    componentValues["minecraft:custom_data"] = customData

    components := sortedItemGiveComponents(componentValues)
    itemID := strings.TrimSpace(entry.ItemID)
    if len(components) == 0 {
        return itemID, nil
    }
    return fmt.Sprintf("%s[%s]", itemID, strings.Join(components, ",")), nil
}
```

`ItemToGiveCommand` を簡素化:

```go
func ItemToGiveCommand(entry, activesByID, passivesByID, bowsByID, ver ...int64) (string, error) {
    itemStack, err := buildItemStack(entry, activesByID, passivesByID, bowsByID, itemVersion(ver...))
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("give @p %s 1", itemStack), nil
}
```

`ItemToUpdateCommand` を簡素化（末尾4行はそのまま維持）:

```go
func ItemToUpdateCommand(entry, activesByID, passivesByID, bowsByID, ver ...int64) (string, error) {
    itemStack, err := buildItemStack(entry, activesByID, passivesByID, bowsByID, itemVersion(ver...))
    if err != nil {
        return "", err
    }
    return strings.Join([]string{
        "scoreboard players set #maf_update_damage tmp 0",
        `$execute store result score #maf_update_damage tmp run data get entity @s $(equip).components."minecraft:damage"`,
        fmt.Sprintf("$item replace entity @s $(slot) with %s 1", itemStack),
        `$execute store result entity @s $(equip).components."minecraft:damage" int 1 run scoreboard players get #maf_update_damage tmp`,
    }, "\n"), nil
}
```

注意: 両公開関数のシグネチャ（可変長 `ver ...int64`）は変更しない。呼び出し元（`export/item.go:27`, `:46`）は無改修。`buildItemStack` には解決済みの `int64` を渡す（`itemVersion(ver...)` 経由）。

## 不変条件

- `ItemToGiveCommand` の出力は現状とバイト単位で一致（空コンポーネント時 `give @p {id} 1`、非空時 `give @p {id}[...] 1`）。
- `ItemToUpdateCommand` の4行構造・itemStack 形式とも不変。

## 検証

```
cd maf-command-generator && make check
```

特に既存テスト（`convert/item_test.go` の `TestItemToGiveCommand*` / `TestItemToUpdateCommand*`）が無改修で全通過することを確認。出力ゴールデン（`testdata/export/basic/output/function/generated/item/give/`, `.../item/update/`）も差分なしであること。新規テストは不要（挙動不変のため既存テストでカバー済み）。
