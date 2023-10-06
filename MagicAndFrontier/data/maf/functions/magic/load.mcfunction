# 魔導書に設定されているIDをそのまま代入する先
scoreboard objectives add p7_magicID dummy

scoreboard objectives add p7_MP dummy
scoreboard objectives add p7_MaxMP dummy


scoreboard objectives add p7_castCost dummy
scoreboard objectives add p7_castTime dummy
# あえてmagicIDと二重にすることで、コストが違う同効果魔法の実装などを簡単に実現
scoreboard objectives add p7_castID dummy
# 詠唱時間の表示用
scoreboard objectives add p7_castTimeMax dummy
 
# MP自然回復タイマー
scoreboard objectives add p7_MPTick dummy

# scoreboard objectives add p7_

# tellraw @a [{"text":"データベースを設定"}]
function maf:magic/setdb

bossbar add mpbar1 "MP"
bossbar set minecraft:mpbar1 style progress
bossbar set minecraft:mpbar1 color green
bossbar add mpbar2 "MP"
bossbar set minecraft:mpbar2 style progress
bossbar set minecraft:mpbar2 color green
bossbar add mpbar3 "MP"
bossbar set minecraft:mpbar3 style progress
bossbar set minecraft:mpbar3 color green
bossbar add mpbar4 "MP"
bossbar set minecraft:mpbar4 style progress
bossbar set minecraft:mpbar4 color green
bossbar add mpbar5 "MP"
bossbar set minecraft:mpbar5 style progress
bossbar set minecraft:mpbar5 color green
bossbar add mpbar6 "MP"
bossbar set minecraft:mpbar6 style progress
bossbar set minecraft:mpbar6 color green
bossbar add mpbar7 "MP"
bossbar set minecraft:mpbar7 style progress
bossbar set minecraft:mpbar7 color green
bossbar add mpbar8 "MP"
bossbar set minecraft:mpbar8 style progress
bossbar set minecraft:mpbar8 color green
bossbar add mpbar9 "MP"
bossbar set minecraft:mpbar9 style progress
bossbar set minecraft:mpbar9 color green
bossbar add mpbar10 "MP"
bossbar set minecraft:mpbar10 style progress
bossbar set minecraft:mpbar10 color green
bossbar add mpbar11 "MP"
bossbar set minecraft:mpbar11 style progress
bossbar set minecraft:mpbar11 color green
bossbar add mpbar12 "MP"
bossbar set minecraft:mpbar12 style progress
bossbar set minecraft:mpbar12 color green
bossbar add mpbar13 "MP"
bossbar set minecraft:mpbar13 style progress
bossbar set minecraft:mpbar13 color green
bossbar add mpbar14 "MP"
bossbar set minecraft:mpbar14 style progress
bossbar set minecraft:mpbar14 color green
bossbar add mpbar15 "MP"
bossbar set minecraft:mpbar15 style progress
bossbar set minecraft:mpbar15 color green
bossbar add mpbar16 "MP"
bossbar set minecraft:mpbar16 style progress
bossbar set minecraft:mpbar16 color green
bossbar add mpbar17 "MP"
bossbar set minecraft:mpbar17 style progress
bossbar set minecraft:mpbar17 color green
bossbar add mpbar18 "MP"
bossbar set minecraft:mpbar18 style progress
bossbar set minecraft:mpbar18 color green
bossbar add mpbar19 "MP"
bossbar set minecraft:mpbar19 style progress
bossbar set minecraft:mpbar19 color green
bossbar add mpbar20 "MP"
bossbar set minecraft:mpbar20 style progress
bossbar set minecraft:mpbar20 color green
