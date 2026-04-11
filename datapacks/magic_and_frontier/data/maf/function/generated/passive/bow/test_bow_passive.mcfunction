effect give @e[distance=..4] minecraft:glowing 3 10
effect give @e[distance=..4,type=!minecraft:player] minecraft:levitation 1 2
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 1 2 1
tell @a "矢の効果が発動しました"
summon minecraft:zombie ~ ~ ~
