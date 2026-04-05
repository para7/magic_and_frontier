data remove storage p7:maf grimoire.dispatch
execute store result storage p7:maf grimoire.dispatch.castid int 1 run scoreboard players get @s mafEffectID
function maf:magic/cast/dispatch/read_effect_ref with storage p7:maf grimoire.dispatch
execute if data storage p7:maf grimoire.dispatch.ref run function maf:magic/cast/dispatch/run_effect_ref with storage p7:maf grimoire.dispatch
execute if entity @s[scores={mafEffectID=1001}] run function maf:generated/passive/apply/passive_1_slot1
execute if entity @s[scores={mafEffectID=1011}] run function maf:generated/passive/apply/regeneration_slot1
