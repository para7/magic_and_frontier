execute as @e[type=#maf:enemymob,distance=..32,limit=1,sort=nearest] at @s run summon minecraft:creeper ~ ~ ~ {NoAI:1b,Fuse:0s,ExplosionRadius:4b,Tags:["vh"]}
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は エクスプロウド を唱えた！"}]
