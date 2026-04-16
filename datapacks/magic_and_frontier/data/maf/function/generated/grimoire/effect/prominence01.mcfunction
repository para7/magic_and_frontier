#プロミネンス
fill ~-7 ~-7 ~-7 ~7 ~7 ~7 fire replace air
fill ~-1 ~-1 ~-1 ~1 ~1 ~1 air replace fire
execute as @e[distance=1..8] run effect give @s minecraft:slowness 5 3
execute as @e[distance=1..8,type=#maf:undead] run effect give @s minecraft:instant_health 1 1
execute as @e[distance=1..8,type=!#maf:undead] run effect give @s minecraft:instant_damage 1 1
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5
playsound minecraft:entity.blaze.shoot master @a ~ ~ ~ 1 1
playsound minecraft:entity.generic.explode master @a ~ ~ ~
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は プロミネンス を唱えた！"}]
