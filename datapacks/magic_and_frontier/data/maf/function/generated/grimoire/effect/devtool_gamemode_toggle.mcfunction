execute if entity @s[gamemode=creative] run tag @s add maf_devtool_was_creative
execute if entity @s[tag=maf_devtool_was_creative] run gamemode survival @s
execute if entity @s[tag=maf_devtool_was_creative] run tellraw @s [{"text":"[Debug] ゲームモードをサバイバルに切り替えました","color":"yellow"}]
execute unless entity @s[tag=maf_devtool_was_creative] run gamemode creative @s
execute unless entity @s[tag=maf_devtool_was_creative] run tellraw @s [{"text":"[Debug] ゲームモードをクリエイティブに切り替えました","color":"yellow"}]
playsound minecraft:item.flintandsteel.use master @s ~ ~ ~ 1 1
tag @s remove maf_devtool_was_creative
