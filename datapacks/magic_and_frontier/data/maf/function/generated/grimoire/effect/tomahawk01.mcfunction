function maf:common/sight/eyes_tagged
execute as @e[type=#maf:enemymob,tag=maf_target,distance=..8,sort=nearest,limit=1] at @s run particle minecraft:cloud ~ ~1 ~ 0.75 0.25 0.75 0.01 12 force
#execute as @e[type=#maf:enemymob,tag=maf_target,distance=..8,sort=nearest,limit=1] run damage @s 20 minecraft:player_attack
damage @e[type=#maf:enemymob,tag=maf_target,distance=..8,sort=nearest,limit=1] 20 minecraft:player_attack by @s
function maf:common/reduce_durability {amount: 5}
playsound minecraft:entity.player.attack.sweep master @a ~ ~ ~ 1.5 0.7
tellraw @s [{"text":"トマホーク！"}]
