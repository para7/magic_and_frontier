tag @e[type=#maf:enemymob,tag=maf_sight_target] remove maf_sight_target
tag @e[type=#maf:enemymob,tag=maf_sight_candidate] remove maf_sight_candidate
function maf:common/eyes_tagged
execute as @e[type=#maf:enemymob,tag=maf_sight_target,distance=..8,sort=nearest,limit=1] run damage @s 20 minecraft:player_attack
execute as @e[type=#maf:enemymob,tag=maf_sight_target,distance=..8,sort=nearest,limit=1] at @s run particle minecraft:crit ~ ~0.9 ~ 0.3 0.5 0.3 0.01 20 force
tag @e[type=#maf:enemymob,tag=maf_sight_candidate] remove maf_sight_candidate
tag @e[type=#maf:enemymob,tag=maf_sight_target] remove maf_sight_target
playsound minecraft:entity.player.attack.sweep master @a ~ ~ ~ 1.5 0.7
tellraw @s [{"text":"トマホーク！"}]
