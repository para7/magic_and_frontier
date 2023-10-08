execute as @e[type=#p7b:enemymob,distance=..50,nbt={OnGround:1b}] run damage @s 40 minecraft:fall by @p

tellraw @a[distance=..60] [{"selector":"@s"},{"text":" は グランドシェイクIII を唱えた！"}]

function maf:magic/exec/01_attack/effect
