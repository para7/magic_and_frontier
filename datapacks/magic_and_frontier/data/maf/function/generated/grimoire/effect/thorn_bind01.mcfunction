execute as @e[type=#maf:enemymob,distance=..8,nbt={OnGround:1b}] run damage @s 6 minecraft:magic
effect give @e[type=#maf:enemymob,distance=..8,nbt={OnGround:1b}] minecraft:slowness 3 0
playsound minecraft:block.vine.place master @a ~ ~ ~ 1.5 0.8
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は ソーンバインド を唱えた！"}]
