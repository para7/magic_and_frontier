マイクラで、魔法書をbookで作り、右クリックで使用検知して発動させたい。特定のタグがついた本をメインハンドに持って、右クリックしたことをトリガーにして、効果を発動させる処理のサンプルや、giveコマンドの例を示して。


できます。
`1.21.2` 以降は、`food` 単体では「使用可能化」されず、**右クリックで使えるようにする役目は `consumable`** が持つようになっています。なので、今の系統では **`book` に `consumable` を付けて、`minecraft:using_item` advancement で検知する** のが基本です。`26.1` は新しい版番号表記で、直前の通常版が `1.21.11` です。([Minecraft Feedback][1])

以下は、**「特定タグ付きの本をメインハンドで右クリックしたら魔法を発動」**する最小サンプルです。
方針はこうです。

1. `book` に `custom_data` で識別タグを付ける
2. `consumable` でその本を「使用可能」にする
3. `minecraft:using_item` advancement で使用開始を検知する
4. 関数側で **メインハンド所持を再確認** してから効果を発動する
5. `use_cooldown` を付けて連打を抑える

`using_item` は「使用中」に反応するため、何も工夫しないと連続で走りやすいです。コミュニティでは `consumable` + `using_item` を使う方法が一般的で、クールダウンや進捗の即時 revoke を併用する形がよく使われています。([Reddit][2])

---

## 1. give コマンド例

まず魔法書を配るコマンドです。

```mcfunction
/give @p minecraft:book[
  minecraft:custom_data={magic_book:"firebolt"},
  minecraft:item_name='"火炎の魔法書"',
  minecraft:lore=['"右クリックで火炎魔法を発動"'],
  minecraft:consumable={
    consume_seconds:99999,
    animation:"bow",
    sound:"entity.blaze.shoot",
    has_consume_particles:false
  },
  minecraft:use_cooldown={
    seconds:1.5
  }
]
```

ポイント:

* `custom_data={magic_book:"firebolt"}` で識別
* `consume_seconds:99999` で「使い切る」前に検知する用途に寄せる
* `has_consume_particles:false` で食事っぽい違和感を減らす
* `use_cooldown` で連打を抑制

`food` は現在「食べたときの効果データ」を持つだけで、使用可能化は `consumable` が担当です。([Minecraft Feedback][1])

---

## 2. advancement で右クリック検知

`data/magic/advancement/use_firebolt.json`

```json
{
  "criteria": {
    "use_book": {
      "trigger": "minecraft:using_item",
      "conditions": {
        "item": {
          "predicates": {
            "minecraft:custom_data": "{magic_book:\"firebolt\"}"
          }
        }
      }
    }
  },
  "rewards": {
    "function": "magic:firebolt_cast"
  }
}
```

これで、`custom_data` が一致する本の使用開始を拾えます。
この方式は `consumable` を付けたアイテムと `using_item` を組み合わせるやり方です。([Reddit][2])

---

## 3. 発動関数

`data/magic/function/firebolt_cast.mcfunction`

```mcfunction
# 進捗を戻して次回も発火できるようにする
advancement revoke @s only magic:use_firebolt

# 本当にメインハンドでその本を使っている時だけ発動
execute if items entity @s weapon.mainhand minecraft:book[minecraft:custom_data~{magic_book:"firebolt"}] run function magic:firebolt_effect
```

ここで **メインハンド限定** にしています。
advancement 側は「そのアイテムの使用」を拾いますが、ユーザー条件どおりに厳密にするなら、関数内で `weapon.mainhand` を再チェックしておくのが安全です。

---

## 4. 実際の魔法効果

`data/magic/function/firebolt_effect.mcfunction`

```mcfunction
# 演出
playsound minecraft:entity.blaze.shoot player @a ~ ~ ~ 1 1
particle minecraft:flame ^ ^1 ^1 0.2 0.2 0.2 0.01 20 force @a

# 前方にいる敵へダメージ例
execute at @s positioned ^ ^1 ^3 as @e[type=!player,distance=..2,limit=1,sort=nearest] run damage @s 6 minecraft:magic

# あるいは自分の前に小爆発風の演出
execute at @s positioned ^ ^1 ^3 run particle minecraft:small_flame ~ ~ ~ 0.3 0.3 0.3 0.02 40 force @a
```

これはただの例です。
ここを差し替えれば、回復、テレポート、前方扇形攻撃、召雷など何でもできます。

---

## 5. もっと簡単な 1 ファイル版

とりあえず試すだけなら、効果関数に全部まとめても動きます。

`data/magic/function/firebolt_cast.mcfunction`

```mcfunction
advancement revoke @s only magic:use_firebolt
execute if items entity @s weapon.mainhand minecraft:book[minecraft:custom_data~{magic_book:"firebolt"}] at @s run playsound minecraft:entity.blaze.shoot player @a ~ ~ ~ 1 1
execute if items entity @s weapon.mainhand minecraft:book[minecraft:custom_data~{magic_book:"firebolt"}] at @s run particle minecraft:flame ^ ^1 ^1 0.2 0.2 0.2 0.01 20 force @a
execute if items entity @s weapon.mainhand minecraft:book[minecraft:custom_data~{magic_book:"firebolt"}] at @s positioned ^ ^1 ^3 as @e[type=!player,distance=..2,limit=1,sort=nearest] run damage @s 6 minecraft:magic
```

---

## 6. 複数の魔法書に増やす場合

識別子だけ変えれば増やせます。

```mcfunction
/give @p minecraft:book[
  minecraft:custom_data={magic_book:"heal"},
  minecraft:item_name='"治癒の魔法書"',
  minecraft:consumable={consume_seconds:99999,animation:"bow",has_consume_particles:false},
  minecraft:use_cooldown={seconds:2}
]
```

対応する advancement:

```json
{
  "criteria": {
    "use_book": {
      "trigger": "minecraft:using_item",
      "conditions": {
        "item": {
          "predicates": {
            "minecraft:custom_data": "{magic_book:\"heal\"}"
          }
        }
      }
    }
  },
  "rewards": {
    "function": "magic:heal_cast"
  }
}
```

---

## 7. 注意点

* **本は元々ふつうの右クリック使用アイテムではない**ので、右クリック検知したいなら `consumable` を付けるのが肝です。([Minecraft Feedback][1])
* `minecraft:using_item` は使用中に発火するので、**advancement の revoke** と **`use_cooldown`** を併用した方が扱いやすいです。([Reddit][2])
* `1.21.3` 付近の古いサンプルが、そのまま `1.21.11` で動かないという報告もあります。なので、今作るなら **`food` ベースの古い書き方ではなく `consumable` 前提** で組むのが安全です。([Planet Minecraft][3])

---

## 最小構成まとめ

**give**

```mcfunction
/give @p minecraft:book[
  minecraft:custom_data={magic_book:"firebolt"},
  minecraft:item_name='"火炎の魔法書"',
  minecraft:consumable={consume_seconds:99999,animation:"bow",has_consume_particles:false},
  minecraft:use_cooldown={seconds:1.5}
]
```

**advancement**

```json
{
  "criteria": {
    "use_book": {
      "trigger": "minecraft:using_item",
      "conditions": {
        "item": {
          "predicates": {
            "minecraft:custom_data": "{magic_book:\"firebolt\"}"
          }
        }
      }
    }
  },
  "rewards": {
    "function": "magic:firebolt_cast"
  }
}
```

**function**

```mcfunction
advancement revoke @s only magic:use_firebolt
execute if items entity @s weapon.mainhand minecraft:book[minecraft:custom_data~{magic_book:"firebolt"}] run function magic:firebolt_effect
```

必要なら次に、
**そのまま datapack のフォルダ構成込みでコピペできる完成形** を出します。

[1]: https://feedback.minecraft.net/hc/en-us/articles/31261174284557-Minecraft-Java-Edition-1-21-2-Bundles-of-Bravery?utm_source=chatgpt.com "Minecraft: Java Edition 1.21.2 - Bundles of Bravery"
[2]: https://www.reddit.com/r/MinecraftCommands/comments/1jrg08r/using_advancement_for_right_click_detect/?utm_source=chatgpt.com "Using Advancement for Right Click detect"
[3]: https://www.planetminecraft.com/data-pack/009-right-click-detection-with-edible-item-minecraft-1-21-datapack-tutorials/?utm_source=chatgpt.com "009: Right click detection with edible item [Minecraft 1.21 ..."
