execute at @e[type=#maf:enemymob,distance=..128,sort=nearest,limit=1] run playsound minecraft:entity.generic.explode player @a ~ ~ ~ 10 1
execute at @e[type=#maf:enemymob,sort=nearest,limit=1] run summon creeper ~ ~ ~ {NoAI:1b,powered:1b,ExplosionRadius:5b,Fuse:0,CustomName:'{"text":"ディザスター"}',Attributes:[{Name:generic.attackDamage,Base:4}]}
execute at @e[type=#maf:enemymob,sort=nearest,limit=1] run summon lightning_bolt ~ ~ ~
playsound minecraft:entity.evoker.cast_spell player @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell player @a ~ ~ ~ 2 0.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は ディザスター を唱えた！"}]
