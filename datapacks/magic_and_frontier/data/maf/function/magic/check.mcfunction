tell @s magic/check

# 杖・魔法の使用判定

execute if entity @s[scores={mafUseWand=1..}] if data entity @s SelectedItem.components."minecraft:custom_data".maf.spell.castid run execute store result score @s mafCastID run data get entity @s SelectedItem.components."minecraft:custom_data".maf.spell.castid

execute if entity @s[scores={mafUseWand=1..}] if data entity @s Inventory[{Slot:-106b}].components."minecraft:custom_data".maf.spell.castid run execute store result score @s mafCastID run data get entity @s Inventory[{Slot:-106b}].components."minecraft:custom_data".maf.spell.castid

# アビリティ持ちの杖処理
