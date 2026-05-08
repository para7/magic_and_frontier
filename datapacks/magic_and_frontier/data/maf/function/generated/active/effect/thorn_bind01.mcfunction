execute as @e[type=#maf:enemymob,distance=..8,nbt={OnGround:1b}] run damage @s 6 minecraft:magic
effect give @e[type=#maf:enemymob,distance=..8,nbt={OnGround:1b}] minecraft:slowness 3 0
execute at @s run particle minecraft:witch ~ ~0.2 ~ 4.0 0.2 4.0 0.01 200 force
execute as @e[type=#maf:enemymob,distance=..8,nbt={OnGround:1b}] at @s run particle minecraft:cloud ~ ~0.5 ~ 0.75 0.25 0.75 0.01 12 force
playsound minecraft:block.vine.place master @a ~ ~ ~ 1.5 0.8
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は ソーンバインド を唱えた！"}]
