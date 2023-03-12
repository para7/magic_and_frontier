tellraw @a[nbt={Dimension:"minecraft:overworld"}] [{"selector":"@s"},{"text":" は ベッドで寝ています…"}]

execute as @a[nbt={Dimension:"minecraft:overworld"}] at @s run playsound minecraft:block.bell.use master @s

scoreboard players set @s BedNotification 0
