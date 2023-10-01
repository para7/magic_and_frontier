import { ZeroPadding } from "@/lib/ZeroPadding";
import { ConvertCategory } from "../ConvertCategory";
import { DBSchemaType } from "../zod/db";
import { RemoveDuplicatedCast } from "../RemoveDuplicatedCast";

const SummonBookCommand = (x: DBSchemaType["data"][number]) => {
  return `tellraw @s [{"text":"未実装です。開発者に連絡してください。 ID: ${x.effectid}"}]`;
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
