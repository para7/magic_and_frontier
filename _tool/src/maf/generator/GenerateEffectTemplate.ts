import { ZeroPadding } from "@/lib/ZeroPadding";
import { ConvertCategory } from "../ConvertCategory";
import { DBSchemaType } from "../zod/db";
import { RemoveDuplicatedCast } from "../RemoveDuplicatedCast";

const SummonBookCommand = (x: DBSchemaType["data"][number]) => {
  return `give @p book{display:{Name:'{"text":"${x.title}","color":"gold"}',Lore:['{"text":"魔法書"}','{"text":"${x.description}","color":"aqua"}']},grimoire:1,magicID:${x.castid},Enchantments:[{}]} 1`;
};

export const GenerateEffectTemplate = (_data: DBSchemaType) => {
  const data = RemoveDuplicatedCast(_data);
  const files = data.map((x) => {
    const filename = `${ConvertCategory(x.effectid)}/${ZeroPadding(x.effectid, 6)}.mcfunction`;
    const command = SummonBookCommand(x);

    return { filename, command };
  });

  return files;
};
