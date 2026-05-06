function maf:common/buff/set {buff_id:"parry",tick:7,tick_function:"parry_tick",destructor_function:"parry_destructor"}
effect give @s minecraft:resistance 1 9 true
execute at @s run particle minecraft:enchanted_hit ~ ~1.0 ~ 0.35 0.45 0.35 0.05 24 force
playsound minecraft:entity.player.attack.sweep master @a[distance=..24] ~ ~ ~ 1.2 1.1
