# effect give @e[type=#p7b:undead,sort=nearest,limit=1] minecraft:instant_health 1 1
# effect give @e[type=#p7b:undead,sort=nearest,limit=1] minecraft:slowness 14 40
# effect give @e[type=#p7b:undead,sort=nearest,limit=1] minecraft:glowing 14 0
# effect give @e[type=#p7b:undead,sort=nearest,limit=1] minecraft:weakness 14 40
execute as @e[type=#p7b:undead,distance=..9,sort=nearest,limit=1] run effect give @s minecraft:instant_health 1 5 
execute at @e[type=#p7b:undead,distance=..9,sort=nearest,limit=1] run particle minecraft:happy_villager ~ ~1 ~ 0.5 0.5 0.5 1 50
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 1 2.0
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は ホーリー を唱えた！"}]
