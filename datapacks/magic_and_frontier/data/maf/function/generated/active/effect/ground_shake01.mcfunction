execute as @e[distance=..10,type=#maf:enemymob,nbt={OnGround:1b}] run damage @s 20 minecraft:magic
playsound minecraft:entity.ravager.roar master @a ~ ~ ~ 2 0.8
playsound minecraft:entity.generic.explode master @a ~ ~ ~ 1 0.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は グランドシェイク を唱えた！"}]
