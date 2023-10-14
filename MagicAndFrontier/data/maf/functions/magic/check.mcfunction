tell @s magic/check

# 杖・魔法の使用判定

execute if entity @s[scores={mafUseWand=1..},nbt={SelectedItem:{tag:{grimoire:1}},Inventory:[{Slot:-106b,tag:{wandID:1}}]}] run execute store result score @s mafMagicID run data get entity @s SelectedItem.tag.magicID

execute if entity @s[scores={mafUseWand=1..},nbt={SelectedItem:{tag:{wandID:1}},Inventory:[{Slot:-106b,tag:{grimoire:1}}]}] run execute store result score @s mafMagicID run data get entity @s Inventory[-1].tag.magicID

# アビリティ持ちの杖処理
