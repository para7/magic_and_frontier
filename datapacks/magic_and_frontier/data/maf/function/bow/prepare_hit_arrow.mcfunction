tag @s add hit
# 矢に貫通を付与することで着弾した瞬間に消滅しないようにする
data merge entity @s {PierceLevel:2b}
data merge entity @s {item:{components:{"minecraft:potion_contents":{custom_effects:[{id:"minecraft:dolphins_grace",duration:200,amplifier:80}]}}}}
