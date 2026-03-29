# effect give @e[type=#p7b:undead,sort=nearest,limit=1] minecraft:instant_health 1 1
# effect give @e[type=#p7b:undead,sort=nearest,limit=1] minecraft:slowness 14 40
# effect give @e[type=#p7b:undead,sort=nearest,limit=1] minecraft:glowing 14 0
# effect give @e[type=#p7b:undead,sort=nearest,limit=1] minecraft:weakness 14 40

effect clear @e[type=#p7b:friendmob,distance=..10] minecraft:blindness
effect clear @e[type=#p7b:friendmob,distance=..10] minecraft:mining_fatigue
effect clear @e[type=#p7b:friendmob,distance=..10] minecraft:nausea
effect clear @e[type=#p7b:friendmob,distance=..10] minecraft:poison
effect clear @e[type=#p7b:friendmob,distance=..10] minecraft:slowness
effect clear @e[type=#p7b:friendmob,distance=..10] minecraft:weakness
effect clear @e[type=#p7b:friendmob,distance=..10] minecraft:wither


execute at @e[distance=..10,type=#p7b:friendmob] run particle minecraft:happy_villager ~ ~1 ~ 0.5 0.5 0.5 1 10

playsound minecraft:entity.player.levelup player @a ~ ~ ~ 1 2.0

playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5

tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は リフレッシュ を唱えた！"}]