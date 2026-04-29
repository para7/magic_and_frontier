# @s視線先の円範囲ターゲット付与（引数必須: forward, radius, particleCount）
tag @e[type=#maf:enemymob,tag=maf_sight_circle_target] remove maf_sight_circle_target
$execute at @s anchored eyes rotated as @s positioned ^ ^ ^$(forward) as @e[type=#maf:enemymob,distance=..$(radius)] run tag @s add maf_sight_circle_target
$execute at @s anchored eyes rotated as @s positioned ^ ^ ^$(forward) run particle minecraft:witch ~ ~0.2 ~ $(radius) 0.1 $(radius) 0.01 $(particleCount) force
