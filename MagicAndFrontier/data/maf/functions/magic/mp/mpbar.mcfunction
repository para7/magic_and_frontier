# NOTE: @aが多すぎるのでfunction20個に分割した方が処理が軽いか？
bossbar set minecraft:mpbar1 players
execute as @a[scores={mafPlayerID=1}] run bossbar set minecraft:mpbar1 players @s
execute as @a[scores={mafPlayerID=1,p7_castTime=..0}] run bossbar set minecraft:mpbar1 color green
execute as @a[scores={mafPlayerID=1,p7_castTime=..0}] run bossbar set minecraft:mpbar1 style progress
execute as @a[scores={mafPlayerID=1,p7_castTime=..0}] run bossbar set minecraft:mpbar1 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=1,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar1 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=1,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar1 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=1,p7_castTime=1..}] run bossbar set minecraft:mpbar1 color yellow
execute as @a[scores={mafPlayerID=1,p7_castTime=1..}] run bossbar set minecraft:mpbar1 style notched_20
execute as @a[scores={mafPlayerID=1,p7_castTime=1..}] run bossbar set minecraft:mpbar1 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar1.title"}]
execute as @a[scores={mafPlayerID=1,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar1 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=1,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar1 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar2 players
execute as @a[scores={mafPlayerID=2}] run bossbar set minecraft:mpbar2 players @s
execute as @a[scores={mafPlayerID=2,p7_castTime=..0}] run bossbar set minecraft:mpbar2 color green
execute as @a[scores={mafPlayerID=2,p7_castTime=..0}] run bossbar set minecraft:mpbar2 style progress
execute as @a[scores={mafPlayerID=2,p7_castTime=..0}] run bossbar set minecraft:mpbar2 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=2,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar2 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=2,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar2 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=2,p7_castTime=1..}] run bossbar set minecraft:mpbar2 color yellow
execute as @a[scores={mafPlayerID=2,p7_castTime=1..}] run bossbar set minecraft:mpbar2 style notched_20
execute as @a[scores={mafPlayerID=2,p7_castTime=1..}] run bossbar set minecraft:mpbar2 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar2.title"}]
execute as @a[scores={mafPlayerID=2,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar2 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=2,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar2 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar3 players
execute as @a[scores={mafPlayerID=3}] run bossbar set minecraft:mpbar3 players @s
execute as @a[scores={mafPlayerID=3,p7_castTime=..0}] run bossbar set minecraft:mpbar3 color green
execute as @a[scores={mafPlayerID=3,p7_castTime=..0}] run bossbar set minecraft:mpbar3 style progress
execute as @a[scores={mafPlayerID=3,p7_castTime=..0}] run bossbar set minecraft:mpbar3 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=3,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar3 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=3,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar3 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=3,p7_castTime=1..}] run bossbar set minecraft:mpbar3 color yellow
execute as @a[scores={mafPlayerID=3,p7_castTime=1..}] run bossbar set minecraft:mpbar3 style notched_20
execute as @a[scores={mafPlayerID=3,p7_castTime=1..}] run bossbar set minecraft:mpbar3 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar3.title"}]
execute as @a[scores={mafPlayerID=3,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar3 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=3,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar3 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar4 players
execute as @a[scores={mafPlayerID=4}] run bossbar set minecraft:mpbar4 players @s
execute as @a[scores={mafPlayerID=4,p7_castTime=..0}] run bossbar set minecraft:mpbar4 color green
execute as @a[scores={mafPlayerID=4,p7_castTime=..0}] run bossbar set minecraft:mpbar4 style progress
execute as @a[scores={mafPlayerID=4,p7_castTime=..0}] run bossbar set minecraft:mpbar4 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=4,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar4 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=4,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar4 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=4,p7_castTime=1..}] run bossbar set minecraft:mpbar4 color yellow
execute as @a[scores={mafPlayerID=4,p7_castTime=1..}] run bossbar set minecraft:mpbar4 style notched_20
execute as @a[scores={mafPlayerID=4,p7_castTime=1..}] run bossbar set minecraft:mpbar4 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar4.title"}]
execute as @a[scores={mafPlayerID=4,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar4 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=4,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar4 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar5 players
execute as @a[scores={mafPlayerID=5}] run bossbar set minecraft:mpbar5 players @s
execute as @a[scores={mafPlayerID=5,p7_castTime=..0}] run bossbar set minecraft:mpbar5 color green
execute as @a[scores={mafPlayerID=5,p7_castTime=..0}] run bossbar set minecraft:mpbar5 style progress
execute as @a[scores={mafPlayerID=5,p7_castTime=..0}] run bossbar set minecraft:mpbar5 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=5,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar5 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=5,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar5 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=5,p7_castTime=1..}] run bossbar set minecraft:mpbar5 color yellow
execute as @a[scores={mafPlayerID=5,p7_castTime=1..}] run bossbar set minecraft:mpbar5 style notched_20
execute as @a[scores={mafPlayerID=5,p7_castTime=1..}] run bossbar set minecraft:mpbar5 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar5.title"}]
execute as @a[scores={mafPlayerID=5,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar5 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=5,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar5 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar6 players
execute as @a[scores={mafPlayerID=6}] run bossbar set minecraft:mpbar6 players @s
execute as @a[scores={mafPlayerID=6,p7_castTime=..0}] run bossbar set minecraft:mpbar6 color green
execute as @a[scores={mafPlayerID=6,p7_castTime=..0}] run bossbar set minecraft:mpbar6 style progress
execute as @a[scores={mafPlayerID=6,p7_castTime=..0}] run bossbar set minecraft:mpbar6 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=6,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar6 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=6,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar6 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=6,p7_castTime=1..}] run bossbar set minecraft:mpbar6 color yellow
execute as @a[scores={mafPlayerID=6,p7_castTime=1..}] run bossbar set minecraft:mpbar6 style notched_20
execute as @a[scores={mafPlayerID=6,p7_castTime=1..}] run bossbar set minecraft:mpbar6 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar6.title"}]
execute as @a[scores={mafPlayerID=6,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar6 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=6,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar6 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar7 players
execute as @a[scores={mafPlayerID=7}] run bossbar set minecraft:mpbar7 players @s
execute as @a[scores={mafPlayerID=7,p7_castTime=..0}] run bossbar set minecraft:mpbar7 color green
execute as @a[scores={mafPlayerID=7,p7_castTime=..0}] run bossbar set minecraft:mpbar7 style progress
execute as @a[scores={mafPlayerID=7,p7_castTime=..0}] run bossbar set minecraft:mpbar7 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=7,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar7 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=7,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar7 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=7,p7_castTime=1..}] run bossbar set minecraft:mpbar7 color yellow
execute as @a[scores={mafPlayerID=7,p7_castTime=1..}] run bossbar set minecraft:mpbar7 style notched_20
execute as @a[scores={mafPlayerID=7,p7_castTime=1..}] run bossbar set minecraft:mpbar7 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar7.title"}]
execute as @a[scores={mafPlayerID=7,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar7 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=7,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar7 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar8 players
execute as @a[scores={mafPlayerID=8}] run bossbar set minecraft:mpbar8 players @s
execute as @a[scores={mafPlayerID=8,p7_castTime=..0}] run bossbar set minecraft:mpbar8 color green
execute as @a[scores={mafPlayerID=8,p7_castTime=..0}] run bossbar set minecraft:mpbar8 style progress
execute as @a[scores={mafPlayerID=8,p7_castTime=..0}] run bossbar set minecraft:mpbar8 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=8,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar8 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=8,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar8 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=8,p7_castTime=1..}] run bossbar set minecraft:mpbar8 color yellow
execute as @a[scores={mafPlayerID=8,p7_castTime=1..}] run bossbar set minecraft:mpbar8 style notched_20
execute as @a[scores={mafPlayerID=8,p7_castTime=1..}] run bossbar set minecraft:mpbar8 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar8.title"}]
execute as @a[scores={mafPlayerID=8,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar8 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=8,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar8 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar9 players
execute as @a[scores={mafPlayerID=9}] run bossbar set minecraft:mpbar9 players @s
execute as @a[scores={mafPlayerID=9,p7_castTime=..0}] run bossbar set minecraft:mpbar9 color green
execute as @a[scores={mafPlayerID=9,p7_castTime=..0}] run bossbar set minecraft:mpbar9 style progress
execute as @a[scores={mafPlayerID=9,p7_castTime=..0}] run bossbar set minecraft:mpbar9 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=9,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar9 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=9,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar9 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=9,p7_castTime=1..}] run bossbar set minecraft:mpbar9 color yellow
execute as @a[scores={mafPlayerID=9,p7_castTime=1..}] run bossbar set minecraft:mpbar9 style notched_20
execute as @a[scores={mafPlayerID=9,p7_castTime=1..}] run bossbar set minecraft:mpbar9 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar9.title"}]
execute as @a[scores={mafPlayerID=9,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar9 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=9,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar9 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar10 players
execute as @a[scores={mafPlayerID=10}] run bossbar set minecraft:mpbar10 players @s
execute as @a[scores={mafPlayerID=10,p7_castTime=..0}] run bossbar set minecraft:mpbar10 color green
execute as @a[scores={mafPlayerID=10,p7_castTime=..0}] run bossbar set minecraft:mpbar10 style progress
execute as @a[scores={mafPlayerID=10,p7_castTime=..0}] run bossbar set minecraft:mpbar10 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=10,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar10 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=10,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar10 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=10,p7_castTime=1..}] run bossbar set minecraft:mpbar10 color yellow
execute as @a[scores={mafPlayerID=10,p7_castTime=1..}] run bossbar set minecraft:mpbar10 style notched_20
execute as @a[scores={mafPlayerID=10,p7_castTime=1..}] run bossbar set minecraft:mpbar10 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar10.title"}]
execute as @a[scores={mafPlayerID=10,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar10 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=10,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar10 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar11 players
execute as @a[scores={mafPlayerID=11}] run bossbar set minecraft:mpbar11 players @s
execute as @a[scores={mafPlayerID=11,p7_castTime=..0}] run bossbar set minecraft:mpbar11 color green
execute as @a[scores={mafPlayerID=11,p7_castTime=..0}] run bossbar set minecraft:mpbar11 style progress
execute as @a[scores={mafPlayerID=11,p7_castTime=..0}] run bossbar set minecraft:mpbar11 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=11,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar11 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=11,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar11 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=11,p7_castTime=1..}] run bossbar set minecraft:mpbar11 color yellow
execute as @a[scores={mafPlayerID=11,p7_castTime=1..}] run bossbar set minecraft:mpbar11 style notched_20
execute as @a[scores={mafPlayerID=11,p7_castTime=1..}] run bossbar set minecraft:mpbar11 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar11.title"}]
execute as @a[scores={mafPlayerID=11,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar11 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=11,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar11 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar12 players
execute as @a[scores={mafPlayerID=12}] run bossbar set minecraft:mpbar12 players @s
execute as @a[scores={mafPlayerID=12,p7_castTime=..0}] run bossbar set minecraft:mpbar12 color green
execute as @a[scores={mafPlayerID=12,p7_castTime=..0}] run bossbar set minecraft:mpbar12 style progress
execute as @a[scores={mafPlayerID=12,p7_castTime=..0}] run bossbar set minecraft:mpbar12 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=12,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar12 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=12,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar12 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=12,p7_castTime=1..}] run bossbar set minecraft:mpbar12 color yellow
execute as @a[scores={mafPlayerID=12,p7_castTime=1..}] run bossbar set minecraft:mpbar12 style notched_20
execute as @a[scores={mafPlayerID=12,p7_castTime=1..}] run bossbar set minecraft:mpbar12 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar12.title"}]
execute as @a[scores={mafPlayerID=12,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar12 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=12,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar12 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar13 players
execute as @a[scores={mafPlayerID=13}] run bossbar set minecraft:mpbar13 players @s
execute as @a[scores={mafPlayerID=13,p7_castTime=..0}] run bossbar set minecraft:mpbar13 color green
execute as @a[scores={mafPlayerID=13,p7_castTime=..0}] run bossbar set minecraft:mpbar13 style progress
execute as @a[scores={mafPlayerID=13,p7_castTime=..0}] run bossbar set minecraft:mpbar13 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=13,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar13 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=13,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar13 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=13,p7_castTime=1..}] run bossbar set minecraft:mpbar13 color yellow
execute as @a[scores={mafPlayerID=13,p7_castTime=1..}] run bossbar set minecraft:mpbar13 style notched_20
execute as @a[scores={mafPlayerID=13,p7_castTime=1..}] run bossbar set minecraft:mpbar13 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar13.title"}]
execute as @a[scores={mafPlayerID=13,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar13 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=13,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar13 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar14 players
execute as @a[scores={mafPlayerID=14}] run bossbar set minecraft:mpbar14 players @s
execute as @a[scores={mafPlayerID=14,p7_castTime=..0}] run bossbar set minecraft:mpbar14 color green
execute as @a[scores={mafPlayerID=14,p7_castTime=..0}] run bossbar set minecraft:mpbar14 style progress
execute as @a[scores={mafPlayerID=14,p7_castTime=..0}] run bossbar set minecraft:mpbar14 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=14,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar14 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=14,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar14 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=14,p7_castTime=1..}] run bossbar set minecraft:mpbar14 color yellow
execute as @a[scores={mafPlayerID=14,p7_castTime=1..}] run bossbar set minecraft:mpbar14 style notched_20
execute as @a[scores={mafPlayerID=14,p7_castTime=1..}] run bossbar set minecraft:mpbar14 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar14.title"}]
execute as @a[scores={mafPlayerID=14,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar14 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=14,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar14 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar15 players
execute as @a[scores={mafPlayerID=15}] run bossbar set minecraft:mpbar15 players @s
execute as @a[scores={mafPlayerID=15,p7_castTime=..0}] run bossbar set minecraft:mpbar15 color green
execute as @a[scores={mafPlayerID=15,p7_castTime=..0}] run bossbar set minecraft:mpbar15 style progress
execute as @a[scores={mafPlayerID=15,p7_castTime=..0}] run bossbar set minecraft:mpbar15 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=15,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar15 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=15,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar15 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=15,p7_castTime=1..}] run bossbar set minecraft:mpbar15 color yellow
execute as @a[scores={mafPlayerID=15,p7_castTime=1..}] run bossbar set minecraft:mpbar15 style notched_20
execute as @a[scores={mafPlayerID=15,p7_castTime=1..}] run bossbar set minecraft:mpbar15 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar15.title"}]
execute as @a[scores={mafPlayerID=15,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar15 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=15,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar15 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar16 players
execute as @a[scores={mafPlayerID=16}] run bossbar set minecraft:mpbar16 players @s
execute as @a[scores={mafPlayerID=16,p7_castTime=..0}] run bossbar set minecraft:mpbar16 color green
execute as @a[scores={mafPlayerID=16,p7_castTime=..0}] run bossbar set minecraft:mpbar16 style progress
execute as @a[scores={mafPlayerID=16,p7_castTime=..0}] run bossbar set minecraft:mpbar16 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=16,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar16 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=16,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar16 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=16,p7_castTime=1..}] run bossbar set minecraft:mpbar16 color yellow
execute as @a[scores={mafPlayerID=16,p7_castTime=1..}] run bossbar set minecraft:mpbar16 style notched_20
execute as @a[scores={mafPlayerID=16,p7_castTime=1..}] run bossbar set minecraft:mpbar16 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar16.title"}]
execute as @a[scores={mafPlayerID=16,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar16 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=16,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar16 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar17 players
execute as @a[scores={mafPlayerID=17}] run bossbar set minecraft:mpbar17 players @s
execute as @a[scores={mafPlayerID=17,p7_castTime=..0}] run bossbar set minecraft:mpbar17 color green
execute as @a[scores={mafPlayerID=17,p7_castTime=..0}] run bossbar set minecraft:mpbar17 style progress
execute as @a[scores={mafPlayerID=17,p7_castTime=..0}] run bossbar set minecraft:mpbar17 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=17,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar17 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=17,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar17 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=17,p7_castTime=1..}] run bossbar set minecraft:mpbar17 color yellow
execute as @a[scores={mafPlayerID=17,p7_castTime=1..}] run bossbar set minecraft:mpbar17 style notched_20
execute as @a[scores={mafPlayerID=17,p7_castTime=1..}] run bossbar set minecraft:mpbar17 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar17.title"}]
execute as @a[scores={mafPlayerID=17,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar17 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=17,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar17 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar18 players
execute as @a[scores={mafPlayerID=18}] run bossbar set minecraft:mpbar18 players @s
execute as @a[scores={mafPlayerID=18,p7_castTime=..0}] run bossbar set minecraft:mpbar18 color green
execute as @a[scores={mafPlayerID=18,p7_castTime=..0}] run bossbar set minecraft:mpbar18 style progress
execute as @a[scores={mafPlayerID=18,p7_castTime=..0}] run bossbar set minecraft:mpbar18 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=18,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar18 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=18,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar18 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=18,p7_castTime=1..}] run bossbar set minecraft:mpbar18 color yellow
execute as @a[scores={mafPlayerID=18,p7_castTime=1..}] run bossbar set minecraft:mpbar18 style notched_20
execute as @a[scores={mafPlayerID=18,p7_castTime=1..}] run bossbar set minecraft:mpbar18 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar18.title"}]
execute as @a[scores={mafPlayerID=18,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar18 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=18,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar18 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar19 players
execute as @a[scores={mafPlayerID=19}] run bossbar set minecraft:mpbar19 players @s
execute as @a[scores={mafPlayerID=19,p7_castTime=..0}] run bossbar set minecraft:mpbar19 color green
execute as @a[scores={mafPlayerID=19,p7_castTime=..0}] run bossbar set minecraft:mpbar19 style progress
execute as @a[scores={mafPlayerID=19,p7_castTime=..0}] run bossbar set minecraft:mpbar19 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=19,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar19 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=19,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar19 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=19,p7_castTime=1..}] run bossbar set minecraft:mpbar19 color yellow
execute as @a[scores={mafPlayerID=19,p7_castTime=1..}] run bossbar set minecraft:mpbar19 style notched_20
execute as @a[scores={mafPlayerID=19,p7_castTime=1..}] run bossbar set minecraft:mpbar19 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar19.title"}]
execute as @a[scores={mafPlayerID=19,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar19 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=19,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar19 value run scoreboard players get @s p7_castTime
bossbar set minecraft:mpbar20 players
execute as @a[scores={mafPlayerID=20}] run bossbar set minecraft:mpbar20 players @s
execute as @a[scores={mafPlayerID=20,p7_castTime=..0}] run bossbar set minecraft:mpbar20 color green
execute as @a[scores={mafPlayerID=20,p7_castTime=..0}] run bossbar set minecraft:mpbar20 style progress
execute as @a[scores={mafPlayerID=20,p7_castTime=..0}] run bossbar set minecraft:mpbar20 name [{"text":"MP "},{"score":{"name":"@s","objective":"mafMP"}},{"text": " / "}, {"score":{"name":"@s","objective":"mafMaxMP"}}]
execute as @a[scores={mafPlayerID=20,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar20 max run scoreboard players get @s mafMaxMP
execute as @a[scores={mafPlayerID=20,p7_castTime=..0}] run execute store result bossbar minecraft:mpbar20 value run scoreboard players get @s mafMP
execute as @a[scores={mafPlayerID=20,p7_castTime=1..}] run bossbar set minecraft:mpbar20 color yellow
execute as @a[scores={mafPlayerID=20,p7_castTime=1..}] run bossbar set minecraft:mpbar20 style notched_20
execute as @a[scores={mafPlayerID=20,p7_castTime=1..}] run bossbar set minecraft:mpbar20 name [{"text":"Casting... "},{"storage":"p7:mpbar","nbt":"bar20.title"}]
execute as @a[scores={mafPlayerID=20,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar20 max run scoreboard players get @s p7_castTimeMax
execute as @a[scores={mafPlayerID=20,p7_castTime=1..}] run execute store result bossbar minecraft:mpbar20 value run scoreboard players get @s p7_castTime
