import json
import os

TAGS_DIR = "/home/para7/workspaces/minecraft-docker/datapacks/magic_and_frontier/data/maf/tags"

def update_json(filepath, additions):
    with open(filepath, 'r') as f:
        data = json.load(f)
    
    values = set(data.get("values", []))
    for add in additions:
        values.add(add)
    
    # Optional removals? 
    # (e.g. if we want to remove something)

    data["values"] = sorted(list(values))
    
    with open(filepath, 'w') as f:
        json.dump(data, f, indent=2)
        f.write('\n')

# Entity additions
enemymobs = [
    "minecraft:warden",
    "minecraft:bogged",
    "minecraft:breeze"
]

friendmobs = [
    "minecraft:axolotl",
    "minecraft:glow_squid",
    "minecraft:goat",
    "minecraft:tadpole",
    "minecraft:frog",
    "minecraft:allay",
    "minecraft:camel",
    "minecraft:sniffer",
    "minecraft:armadillo",
    "minecraft:strider",
    "minecraft:wandering_trader",
    "minecraft:snow_golem",
    "minecraft:trader_llama",
    "minecraft:bat"
]

update_json(os.path.join(TAGS_DIR, "entity_type", "enemymob.json"), enemymobs)
update_json(os.path.join(TAGS_DIR, "entity_type", "enemymob_notboss.json"), [m for m in enemymobs if m != "minecraft:warden"]) # warden is boss? let's add it to both and maybe remove from notboss later if needed. Actually warden is boss-like, but let's add it. Wait, the user just wants up to date 1.21.11 info.
update_json(os.path.join(TAGS_DIR, "entity_type", "friendmob.json"), friendmobs)
update_json(os.path.join(TAGS_DIR, "entity_type", "mobs.json"), enemymobs + friendmobs)
update_json(os.path.join(TAGS_DIR, "entity_type", "skeletons.json"), ["minecraft:bogged", "minecraft:wither_skeleton"])
update_json(os.path.join(TAGS_DIR, "entity_type", "undead.json"), ["minecraft:bogged"])
update_json(os.path.join(TAGS_DIR, "entity_type", "zombies.json"), []) # No new zombies
update_json(os.path.join(TAGS_DIR, "entity_type", "water_enemy.json"), []) # No new water enemies (maybe bogged? they spawn in swamps, but not exactly "water" like drowneds).

# Items
update_json(os.path.join(TAGS_DIR, "item", "swords.json"), ["minecraft:stone_sword"])
