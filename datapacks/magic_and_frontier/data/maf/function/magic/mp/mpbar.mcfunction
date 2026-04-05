# TODO: この20プレイヤー固定を修正できないだろうか？
bossbar set minecraft:mpbar1 players
bossbar set minecraft:mpbar2 players
bossbar set minecraft:mpbar3 players
bossbar set minecraft:mpbar4 players
bossbar set minecraft:mpbar5 players
bossbar set minecraft:mpbar6 players
bossbar set minecraft:mpbar7 players
bossbar set minecraft:mpbar8 players
bossbar set minecraft:mpbar9 players
bossbar set minecraft:mpbar10 players
bossbar set minecraft:mpbar11 players
bossbar set minecraft:mpbar12 players
bossbar set minecraft:mpbar13 players
bossbar set minecraft:mpbar14 players
bossbar set minecraft:mpbar15 players
bossbar set minecraft:mpbar16 players
bossbar set minecraft:mpbar17 players
bossbar set minecraft:mpbar18 players
bossbar set minecraft:mpbar19 players
bossbar set minecraft:mpbar20 players

execute as @a[scores={mafPlayerID=1..20}] run function maf:magic/mp/mpbar_per_player_dispatch
