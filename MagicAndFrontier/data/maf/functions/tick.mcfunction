execute as @a at @s run function maf:magic/tick

# 自動加算はここで消す
# dummyは各tick内で消す
scoreboard players set @a p7_useWand 0
