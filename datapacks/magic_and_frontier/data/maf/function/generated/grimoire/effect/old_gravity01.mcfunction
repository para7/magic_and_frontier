# effect give @a[distance=..4] minecraft:speed 24 0
execute as @e[type=#maf:enemymob_notboss,nbt={OnGround:0b},distance=..18] run data merge entity @s {Motion:[0.0,-5.5,0.0]}
execute as @e[type=minecraft:bat,nbt={OnGround:0b},distance=..18] run data merge entity @s {Motion:[0.0,-5.5,0.0]}
# execute as @e[type=#maf:enemymob,nbt={OnGround:0b}] run effect give @s minecraft:slowness 5 100 
# execute as @e[type=minecraft:bat,nbt={OnGround:0b}] run data merge entity @s {NoGravity:0b,NoAI:1b}
particle minecraft:cloud ~0 ~0.5 ~ 0.5 0.8 0.5 0.2 50 normal
playsound minecraft:entity.ender_dragon.ambient player @a ~ ~ ~ 0.8 1.5
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5
tellraw @a[distance=..2] [{"selector":"@s"},{"text":" は グラビデ を唱えた！"}]
