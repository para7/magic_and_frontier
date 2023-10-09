import { ZeroPadding } from "@/lib/ZeroPadding";
import { ConvertCategory } from "../ConvertCategory";
import { DBSchemaType } from "../zod/db";
import { RemoveDuplicatedCast } from "../RemoveDuplicatedCast";

const BookModifires = (x: DBSchemaType["data"][number]) => {
  return `[
    {
        "function": "minecraft:set_name",
        "name": {
            "text": "${x.title}",
            "color": "gold",
            "italic": true
        }
    },
    {
        "function": "minecraft:set_lore",
        "lore": [
            {
                "text": "魔法書"
            },
            {
                "text": "cost:${x.cost} / cast:${x.cast}",
                "color": "white"
            },
            {
                "text": "${x.description}",
                "color": "aqua"
            }
        ],
        "replace": true
    }
]`;
};

export const GenerateBookModifires = (_data: DBSchemaType) => {
  const data = RemoveDuplicatedCast(_data);
  const files = data.map((x) => {
    const filename = `grimore_${ZeroPadding(x.effectid, 6)}.json`;
    const command = BookModifires(x);

    return { filename, command };
  });

  return files;
};
