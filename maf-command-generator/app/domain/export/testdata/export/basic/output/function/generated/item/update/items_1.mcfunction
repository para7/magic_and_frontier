scoreboard players set #maf_update_damage tmp 0
$execute store result score #maf_update_damage tmp run data get entity @s $(equip).components."minecraft:damage"
$item replace entity @s $(slot) with minecraft:stone[minecraft:custom_data={maf:{item_id:"minecraft:stone",source_id:"items_1",nbt_snapshot:{id:"minecraft:stone",count:1,components:{"minecraft:custom_name":{text:"Starter Stone"}}}}},minecraft:custom_name={text:"Starter Stone"}] 1
$execute store result entity @s $(equip).components."minecraft:damage" int 1 run scoreboard players get #maf_update_damage tmp
