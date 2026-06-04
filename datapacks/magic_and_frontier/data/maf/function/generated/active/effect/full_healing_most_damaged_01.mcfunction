function maf:common/target/most_hp_missing {range:10}
effect give @e[tag=maf_heal_target] minecraft:instant_health 1 255
execute at @e[tag=maf_heal_target] run particle minecraft:heart ~ ~1 ~ 0.5 0.5 0.5 1 12
function maf:common/target/clear
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 1.0
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は フルヒーリング を唱えた！"}]
