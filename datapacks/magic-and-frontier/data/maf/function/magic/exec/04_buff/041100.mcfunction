effect give @e[distance=..10,type=#p7b:friendmob] minecraft:strength 60 3
effect give @e[distance=..10,type=#p7b:friendmob] minecraft:night_vision 60 0
effect give @e[distance=..10,type=#p7b:friendmob] minecraft:resistance 60 0

# アサルトドライブ
tellraw @a[distance=..10] [{"selector":"@s"},{"text":" は "},{"nbt":"magictmp.title","storage":"p7:maf"},{"text": "を唱えた！"}]
